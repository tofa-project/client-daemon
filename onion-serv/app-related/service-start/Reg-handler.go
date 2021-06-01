package service_start

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
	httpcodes "github.com/tofa-project/client-daemon/http-codes"

	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
	oservapps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wsrpc_broadcast "github.com/tofa-project/client-daemon/ws-rpc/broadcast"
	wsrpc_events "github.com/tofa-project/client-daemon/ws-rpc/events-list"
)

// Handles service registration inbound request
func RegisterService(rw *http.ResponseWriter, r *http.Request, rApp *osat.App, iMap map[string]string) {
	w := *rw

	name, isName := iMap["name"]
	description, isDes := iMap["description"]
	if !isName || !isDes {
		w.WriteHeader(400)
		return
	}

	// waits for response from ws pipes
	//
	resChan := make(chan int)
	defer close(resChan)
	var app *osat.App

	// state whether WriteHeader has been invoked or not
	sent := false

	// broadcast request to all connected ws clients
	// ws clients will respond with a JSON containing the verdict
	wsrpc_broadcast.Broadcast(glob.J{
		"type":                     "event",
		"eventName":                wsrpc_events.O_INCOMING_REGISTRATION,
		"localAppName":             rApp.GetData()["name"].(string),
		"remoteServiceName":        name,
		"remoteServiceDescription": description,
	}, func(fIMap glob.J) {
		if sent {
			return
		}

		if fIMap["allow"].(bool) {
			resChanSent := false

			oservapps.Apps.Range(func(key, value interface{}) bool {
				rApp := value.(*osat.App)
				host := fmt.Sprintf("%s.onion:%s", rApp.Onion.ID, rApp.GetData()["port"])

				if host == r.Host {
					app = rApp

					// reply code
					resChanSent = true

					// update data
					rApp.MuxData.Lock()
					rApp.Data["is_registered"] = "1"
					db.UpdateApp(rApp.Data["id"].(string), rApp.Data)
					rApp.MuxData.Unlock()

					resChan <- httpcodes.ACTION_ALLOWED

					return false
				}

				return true
			})

			if !resChanSent {
				resChan <- httpcodes.CL_DA_CONFLICT
			}
		} else {
			resChan <- httpcodes.ACTION_REJECTED
		}
	})

	// handles timeout
	go func() {
		time.Sleep(glob.C_HTTP_RES_TIMEOUT * time.Minute)

		if !sent {
			resChan <- http.StatusRequestTimeout
		}
	}()

	resCode := <-resChan

	w.WriteHeader(resCode)
	if resCode == httpcodes.ACTION_ALLOWED {
		w.Write([]byte(fmt.Sprintf(`{"auth_token":"%s"}`, app.GetData()["id"])))
	}

	sent = true
}
