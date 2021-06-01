package http

import (
	"fmt"
	"log"
	netHttp "net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tofa-project/client-daemon/glob"
	"github.com/tofa-project/client-daemon/lib"
	clients "github.com/tofa-project/client-daemon/ws-rpc/clients"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

// Default websocket upgrader
var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *netHttp.Request) bool { return true },
}

// Default http server handler for incoming connections
//
var httpServerMux = &netHttp.ServeMux{}
var httpServer = &netHttp.Server{
	Handler: httpServerMux,
}

// Closes wrapped websocket connection and removes dude from list (since it's closed now)
func closeWConn(wConn *wconn.WrapppedConn) {
	wConn.Conn.Close()
	clients.Clients.Delete(wConn.ID)
}

// Handles commands arrived from GUI client
func commHandler(wConn *wconn.WrapppedConn, data glob.J) {
	switch data["comm"].(string) {

	case "get-apps":
		commGetApps(wConn, data)
	case "get-app":
		commGetApp(wConn, data)
	case "del-app":
		commDelApp(wConn, data)
	case "stop-app":
		commStopApp(wConn, data)
	case "start-app":
		commStartApp(wConn, data)
	case "make-app":
		commMakeApp(wConn, data)

	default:
		wConn.Send(glob.J{
			"type":   "res",
			"pipeID": data["pipeID"],
		})

	}
}

// Handles websocket messages
func wsMsgHandler(wConn *wconn.WrapppedConn) {
	defer closeWConn(wConn)

	for {
		data := make(glob.J)

		err := wConn.Conn.ReadJSON(&data)
		if err != nil {
			log.Printf("ERROR! %s", err)
			return
		}

		// checks if data is command arrived from client
		if data["type"].(string) == "comm" {
			commHandler(wConn, data)
			continue
		}

		// otherwise iterate all pipes from connection with given data as param
		wConn.Pipes.Range(func(_, f interface{}) bool {
			go f.(func(glob.J))(data)

			return true
		})
	}
}

// checks request authorization header
func checkRequestAuthorizationHeader(value string) error {
	splitted := strings.Split(value, " ")

	if len(splitted) > 1 {
		if splitted[1] == glob.V_WS_RPC_HEADER_SECRET {
			return nil
		}
	}

	return fmt.Errorf("invalid authorization header")
}

// Default handler for incoming connections
func requestHandler(w netHttp.ResponseWriter, r *netHttp.Request) {

	// If header secret is present in globals,
	// authorization header must be present in request!
	if glob.V_WS_RPC_HEADER_SECRET != "" {
		if err := checkRequestAuthorizationHeader(r.Header.Get("Authorization")); err != nil {
			w.WriteHeader(403)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	connID := lib.GenerateRandomString(glob.C_WS_CONN_ID_LENGTH)

	wConn := wconn.WrapppedConn{
		ID:             connID,
		Conn:           conn,
		SendMessageMux: sync.Mutex{},
	}

	clients.Clients.Store(connID, &wConn)

	go wsMsgHandler(&wConn)
}

// Initializes the websocket server
func Init() {
	log.Print("Starting websocket RPC server...")

	httpServer.Addr = glob.V_WS_RPC_HOST
	httpServerMux.HandleFunc("/", requestHandler)

	errCh := make(chan error)

	// start server
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			errCh <- err
		}
	}()

	// invoked in case server started successfully
	go func() {
		time.Sleep(glob.C_HTTP_SERVE_TIMEOUT * time.Millisecond)

		errCh <- nil
	}()

	err := <-errCh
	if err == nil {
		log.Print("OK")
	} else {
		log.Printf("FAIL! %s", err)
	}
}
