package http

import (
	"github.com/tofa-project/client-daemon/glob"
	service_start "github.com/tofa-project/client-daemon/onion-serv/app-related/service-start"
	app_type "github.com/tofa-project/client-daemon/onion-serv/app-type"
	apps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

func commStartApp(wConn *wconn.WrapppedConn, input glob.J) {
	if iApp, is := apps.Apps.Load(input["appID"].(string)); is {
		service_start.StartService(iApp.(*app_type.App))
	}

	wConn.Send(glob.J{
		"pipeID":  input["pipeID"],
		"type":    "res",
		"success": "App stopped!",
	})
}
