package http

import (
	"github.com/tofa-project/client-daemon/glob"
	service_stop "github.com/tofa-project/client-daemon/onion-serv/app-related/service-stop"
	app_type "github.com/tofa-project/client-daemon/onion-serv/app-type"
	apps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

func commStopApp(wConn *wconn.WrapppedConn, input glob.J) {
	if iApp, is := apps.Apps.Load(input["appID"].(string)); is {
		service_stop.StopService(iApp.(*app_type.App))
	}

	wConn.Send(glob.J{
		"pipeID":  input["pipeID"],
		"type":    "res",
		"success": "App stopped!",
	})
}
