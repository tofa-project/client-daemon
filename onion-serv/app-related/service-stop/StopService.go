package service_stop

import (
	"log"

	"github.com/tofa-project/client-daemon/glob"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
	wsrpc_broadcast "github.com/tofa-project/client-daemon/ws-rpc/broadcast"
	wsrpc_events "github.com/tofa-project/client-daemon/ws-rpc/events-list"
)

// Stops an onion service
func StopService(a *osat.App) {
	log.Printf("Stopping onion service [%s]...", a.Data["name"])

	a.MuxOnion.Lock()
	a.MuxHttp.Lock()
	defer a.MuxOnion.Unlock()
	defer a.MuxHttp.Unlock()

	errO := a.Onion.Close()
	a.HttpServer.Close()

	if errO == nil {
		a.Running = false
		log.Print("OK")

		// emit publishing state event to GUI
		wsrpc_broadcast.BroadcastMessage(glob.J{
			"type":      "event",
			"eventName": wsrpc_events.NFO_APP_UNPUBLISHED,
			"appID":     a.GetData()["id"].(string),
		})
	} else {
		log.Panicf("FAIL! %s", errO)
	}
}
