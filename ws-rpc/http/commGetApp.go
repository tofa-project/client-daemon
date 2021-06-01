package http

import (
	"github.com/tofa-project/client-daemon/glob"
	app_type "github.com/tofa-project/client-daemon/onion-serv/app-type"
	apps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

// retrieves specific information of certain app by its ID
func commGetApp(wConn *wconn.WrapppedConn, input glob.J) {
	id := input["appID"].(string)

	if iApp, is := apps.Apps.Load(id); is {
		app := iApp.(*app_type.App)

		out := glob.J{
			"id":      id,
			"name":    app.GetData()["name"],
			"state":   app.GetData()["is_registered"],
			"running": app.Running,
			"port":    app.GetData()["port"],
			"onion":   app.Onion.ID,
			"path":    app.GetData()["path"],
			"version": app.GetData()["version"],
		}

		wConn.Send(glob.J{
			"pipeID": input["pipeID"],
			"type":   "res",
			"app":    out,
		})
	} else {
		wConn.Send(glob.J{
			"pipeID": input["pipeID"],
			"type":   "res",
			"error":  "App does not exist.",
		})
	}
}
