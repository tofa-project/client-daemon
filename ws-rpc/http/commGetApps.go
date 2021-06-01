package http

import (
	"github.com/tofa-project/client-daemon/glob"
	apps_list "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

// Handles incoming "get-apps" command
func commGetApps(wConn *wconn.WrapppedConn, data glob.J) {
	apps := make([]glob.J, 0)

	for id, app := range apps_list.GetApps() {
		apps = append(apps, glob.J{
			"id":      id,
			"name":    app.GetData()["name"],
			"state":   app.GetData()["is_registered"],
			"running": app.Running,
		})
	}

	wConn.Send(glob.J{
		"type":   "res",
		"pipeID": data["pipeID"],
		"apps":   apps,
	})
}
