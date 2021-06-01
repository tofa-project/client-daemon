package service_start

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tofa-project/client-daemon/glob"
	app_type "github.com/tofa-project/client-daemon/onion-serv/app-type"
	oservapps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wsrpc_broadcast "github.com/tofa-project/client-daemon/ws-rpc/broadcast"
	wsrpc_events "github.com/tofa-project/client-daemon/ws-rpc/events-list"

	service_stop "github.com/tofa-project/client-daemon/onion-serv/app-related/service-stop"
)

// binds http server and http mux to app
func bindHttp(a *app_type.App) {
	a.ServeMux = &http.ServeMux{}
	a.HttpServer = &http.Server{
		Handler: a.ServeMux,
	}
	a.HttpServer.MaxHeaderBytes = 1 << 4 // 256K
}

// Checks if request should be allowed
// based on number of bytes received per timeframe interval
func allowRequestBasedOnBytesCount(a *app_type.App, r *http.Request) ([]byte, error) {
	reader := bufio.NewReader(r.Body)

	var bytesCount uint
	var readBytesErr error
	var ret []byte
	var b byte

	// read all bytes until they're exceeded
	b, readBytesErr = reader.ReadByte()
	for readBytesErr == nil {
		bytesCount++

		if bytesCount+a.HttpRecBytes > glob.C_MAX_ONION_SERV_HTTP_REC_BYTES {
			return ret, fmt.Errorf("exceeded bytes num")
		}

		ret = append(ret, b)

		b, readBytesErr = reader.ReadByte()
	}

	// if they're not exceeded append them to current received bytes
	a.HttpRecBytes += bytesCount

	//log.Printf("Bytes count: %d, Total Rec bytes: %d", bytesCount, a.HttpRecBytes)

	// carry on
	return ret, nil
}

// Checks if request should be allowed
// based on number of requests received per timeframe interval
func allowRequestBasedOnReqCount(a *app_type.App) error {
	if a.HttpReqCount+1 > glob.C_MAX_ONION_SERV_HTTP_REC_REQ_COUNT {
		return fmt.Errorf("exceeded max req count")
	}

	a.HttpReqCount++

	//log.Printf("Req count: %d, Total Rec req: %d", 1, a.HttpReqCount)

	return nil
}

// Handles potential DoS
func handlePotentialDoS(a *app_type.App, err error) {
	go service_stop.StopService(a)

	wsrpc_broadcast.BroadcastMessage(glob.J{
		"type":      "event",
		"eventName": wsrpc_events.NFO_APP_DOS,
		"appID":     a.GetData()["id"].(string),
		"message":   err.Error(),
	})
}

// Filters incoming bytes for protection
func filterBytes(a *app_type.App, r *http.Request) ([]byte, error) {
	var res []byte

	// daemon will stop onion service amid potential DoS
	res, DoSBytes := allowRequestBasedOnBytesCount(a, r)
	DoSReqLim := allowRequestBasedOnReqCount(a)

	if DoSBytes != nil {
		return res, DoSBytes
	}

	if DoSReqLim != nil {
		return res, DoSReqLim
	}

	return res, nil
}

// Checks authorization token.
// Returns 200 on success, other code on failure
func checkAuthorizationToken(a *app_type.App, r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	authHeaderSplitted := strings.Split(authHeader, " ")
	if len(authHeaderSplitted) != 2 {
		return http.StatusBadRequest
	}

	authToken := authHeaderSplitted[1]
	app, is := oservapps.Apps.Load(authToken)
	if !is || a != app {
		return http.StatusForbidden
	}

	return 200
}

// Reads whole request bytes into map
func readIntoJson(input []byte) (map[string]string, error) {
	iMap := make(map[string]string)
	err := json.Unmarshal(input, &iMap)
	if err != nil {
		return iMap, err
	}

	return iMap, nil
}

// Reloads http server with mux for *App.
// Must be ran only when service is closed!
func ReloadHttp(a *app_type.App) {
	if a.Running {
		return
	}

	bindHttp(a)

	// proper response to snarly clients
	a.ServeMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := filterBytes(a, r)

		// fake server
		rw.Header().Add("Server", "Apache")
		rw.Header().Add("Connection", "close")
		rw.Header().Add("Content-Type", "text/html; charset=UTF-8")

		// stops onion service amid potential DoS
		if err != nil {
			handlePotentialDoS(a, err)
			rw.WriteHeader(http.StatusServiceUnavailable)
			return
		} else {
			rw.WriteHeader(http.StatusNoContent)
		}
	})

	//	actual webhook handler
	//
	urlPath, _ := url.QueryUnescape(a.GetData()["path"].(string))
	a.ServeMux.HandleFunc("/"+urlPath, func(w http.ResponseWriter, r *http.Request) {

		// stops onion service amid potential DoS
		iBytes, err := filterBytes(a, r)
		if err != nil {
			handlePotentialDoS(a, err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// daemon will never handle two requests per app at the same time!
		if a.IsBlocked {
			w.WriteHeader(http.StatusConflict)
			return
		}

		// Daemon will always return a json
		w.Header().Add("Content-Type", "application/json")

		a.IsBlocked = true
		defer func() { a.IsBlocked = false }()

		switch r.Method {

		// Remote server tests connection before attempting other methods.
		// Sent as preflight to ensure clean Tor circuit between remote server and client
		case "PING":
			w.WriteHeader(http.StatusNoContent)

		// Remote server registers app with the client
		// Required for 2FA to be online
		case "REG":
			// read input into json
			inputJson, err := readIntoJson(iBytes)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// a registration can only occur if application is not registered
			if a.GetData()["is_registered"].(string) == "0" {
				RegisterService(&w, r, a, inputJson)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		// Remote server asks for confirmation from client amid action in relation with an app
		// Client replies with true/false.
		// Useful for 2FA without code type, just quick confirmation
		case "ASK":
			// read input into json
			inputJson, err := readIntoJson(iBytes)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if code := checkAuthorizationToken(a, r); a.GetData()["is_registered"].(string) == "1" && code == 200 {
				Ask(&w, inputJson, a)
			} else {
				w.WriteHeader(code)
			}

		// Remote server outputs information to the client
		// Useful if you want 2FA with code like style
		case "INFO":
			// read input into json
			inputJson, err := readIntoJson(iBytes)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if code := checkAuthorizationToken(a, r); a.GetData()["is_registered"].(string) == "1" && code == 200 {
				Info(&w, inputJson, a)
			} else {
				w.WriteHeader(code)
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
