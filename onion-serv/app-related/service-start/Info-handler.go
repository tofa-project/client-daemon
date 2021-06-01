package service_start

import (
	"log"
	"net/http"

	"github.com/tofa-project/client-daemon/glob"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
	wsrpc_broadcast "github.com/tofa-project/client-daemon/ws-rpc/broadcast"
	wsrpc_events "github.com/tofa-project/client-daemon/ws-rpc/events-list"
)

// Handles info messages from server
func Info(rw *http.ResponseWriter, iMap map[string]string, rApp *osat.App) {
	w := *rw

	if _, is := iMap["description"]; !is {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("bad description not existent")

		return
	}

	wsrpc_broadcast.Broadcast(glob.J{
		"type":            "event",
		"eventName":       wsrpc_events.O_INCOMING_INFO,
		"appName":         rApp.GetData()["name"].(string),
		"infoDescription": iMap["description"],
	}, func(fIMap glob.J) {})

	// done
	w.WriteHeader(200)
}
