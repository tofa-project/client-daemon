package service_start

import (
	"net/http"
	"time"

	"github.com/tofa-project/client-daemon/glob"
	httpcodes "github.com/tofa-project/client-daemon/http-codes"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
	wsrpc_broadcast "github.com/tofa-project/client-daemon/ws-rpc/broadcast"
	wsrpc_events "github.com/tofa-project/client-daemon/ws-rpc/events-list"
)

// Handles server asks inbound requests
func Ask(rw *http.ResponseWriter, iMap map[string]string, app *osat.App) {
	w := *rw

	if _, is := iMap["description"]; !is {
		w.WriteHeader(400)
		return
	}

	// outputs ask message to rpc ws
	writeSent := false
	writeCode := make(chan int)
	defer close(writeCode)
	wsrpc_broadcast.Broadcast(glob.J{
		"type":           "event",
		"eventName":      wsrpc_events.O_INCOMING_ASK,
		"appName":        app.GetData()["name"].(string),
		"askDescription": iMap["description"],
	}, func(fIMap glob.J) {
		if writeSent {
			return
		}

		if fIMap["allow"].(bool) {
			writeCode <- httpcodes.ACTION_ALLOWED
		} else {
			writeCode <- httpcodes.ACTION_REJECTED
		}
	})

	// handles timeout
	go func() {
		time.Sleep(glob.C_HTTP_RES_TIMEOUT * time.Minute)

		if !writeSent {
			writeCode <- 408
		}
	}()

	wCode := <-writeCode
	writeSent = true
	w.WriteHeader(wCode)
}
