package apps

import (
	"sync"

	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
)

// generic type of apps list
type AppList map[string]*osat.App

// Contains all apps from database.
var Apps = sync.Map{}

// Retrieves all registered apps
func GetApps() AppList {
	copyCat := make(AppList)

	Apps.Range(func(key, value interface{}) bool {
		copyCat[key.(string)] = value.(*osat.App)
		return true
	})

	return copyCat
}

// Registers app instance
func RegApp(a *osat.App) {
	Apps.Store(a.GetData()["id"].(string), a)
}
