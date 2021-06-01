// Responsible with onion services handling
package oserv

import (
	"log"

	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
	service_start "github.com/tofa-project/client-daemon/onion-serv/app-related/service-start"
	wrap_app_data "github.com/tofa-project/client-daemon/onion-serv/app-related/wrap-data"
	apps "github.com/tofa-project/client-daemon/onion-serv/apps-list"
)

// json map
type j = glob.J

// Initializes onion services.
// Should be called before tserv.Init()
func Init() {
	log.Print("Loading onion services...")

	appDatas := db.GetApps()

	for _, appData := range appDatas {
		tApp := wrap_app_data.WrapAppData(appData)

		apps.Apps.Store(appData["id"].(string), tApp)

		service_start.StartService(tApp)
	}

	log.Print("OK! Apps broadcasted")
}
