package set_data

import (
	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
	service_start "github.com/tofa-project/client-daemon/onion-serv/app-related/service-start"
	service_stop "github.com/tofa-project/client-daemon/onion-serv/app-related/service-stop"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
)

// Set data. Updates it in database as well.
func SetData(a *osat.App, data glob.J) {
	service_stop.StopService(a)

	a.MuxData.Lock()
	defer a.MuxData.Unlock()

	a.Data = data

	db.UpdateApp(a.Data["id"].(string), a.Data)

	service_start.StartService(a)
}
