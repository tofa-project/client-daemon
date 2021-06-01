package service

import (
	service_start "github.com/tofa-project/client-daemon/onion-serv/app-related/service-start"
	service_stop "github.com/tofa-project/client-daemon/onion-serv/app-related/service-stop"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
)

// Reloads an onion service
func ReloadService(a *osat.App) {
	service_stop.StopService(a)
	service_start.StartService(a)
}
