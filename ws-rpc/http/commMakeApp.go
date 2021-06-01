package http

import (
	"github.com/tofa-project/client-daemon/glob"
	app_create "github.com/tofa-project/client-daemon/onion-serv/app-related/create-app"
	service_start "github.com/tofa-project/client-daemon/onion-serv/app-related/service-start"
	app_type "github.com/tofa-project/client-daemon/onion-serv/app-type"
	apps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

func commMakeApp(wConn *wconn.WrapppedConn, data glob.J) {
	appName := data["appName"].(string)
	appNameTaken := false

	// check if app name is already taken
	//
	apps.Apps.Range(func(key, value interface{}) bool {
		if appName == value.(*app_type.App).GetData()["name"].(string) {
			appNameTaken = true
			return false
		}

		return true
	})

	if appNameTaken {
		wConn.Send(glob.J{
			"pipeID": data["pipeID"],
			"type":   "res",
			"error":  "App name already taken",
		})

		return
	}

	app := app_create.CreateApp(appName)
	apps.RegApp(app)
	service_start.StartService(app)

	// reply
	wConn.Send(glob.J{
		"pipeID":  data["pipeID"],
		"type":    "res",
		"success": "App created and broadcasted!",
	})
}
