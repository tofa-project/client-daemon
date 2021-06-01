package set_data

import (
	"github.com/tofa-project/client-daemon/db"
	service_start "github.com/tofa-project/client-daemon/onion-serv/app-related/service-start"
	service_stop "github.com/tofa-project/client-daemon/onion-serv/app-related/service-stop"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
)

// Set data key. Updates in in database as well.
func SetDataKey(a *osat.App, key string, value interface{}, reloadService bool) {
	if reloadService {
		service_stop.StopService(a)
		defer service_start.StartService(a)
	}

	a.MuxData.Lock()
	defer a.MuxData.Unlock()

	a.Data[key] = value

	db.UpdateApp(a.Data["id"].(string), a.Data)
}
