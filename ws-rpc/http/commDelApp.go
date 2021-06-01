package http

import (
	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
	service_stop "github.com/tofa-project/client-daemon/onion-serv/app-related/service-stop"
	oservapptype "github.com/tofa-project/client-daemon/onion-serv/app-type"
	apps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

// forever deletes app by ID
func commDelApp(wConn *wconn.WrapppedConn, input glob.J) {
	id := input["appID"].(string)

	if iApp, is := apps.Apps.Load(id); is {
		app := iApp.(*oservapptype.App)

		service_stop.StopService(app)

		apps.Apps.Delete(id)

		db.DeleteApp(id)
	}

	wConn.Send(glob.J{
		"pipeID":  input["pipeID"],
		"type":    "res",
		"success": "App deleted!",
	})
}
