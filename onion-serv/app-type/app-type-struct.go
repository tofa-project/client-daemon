package app_type

import (
	"net/http"
	"sync"

	"github.com/cretz/bine/tor"
	"github.com/tofa-project/client-daemon/glob"
)

// An encapsulated app containing afferent DB data and onion service instance
type App struct {
	MuxData  sync.Mutex
	MuxOnion sync.Mutex
	MuxHttp  sync.Mutex

	Data      glob.J
	Onion     *tor.OnionService
	Running   bool
	IsBlocked bool

	HttpServer *http.Server
	ServeMux   *http.ServeMux

	// Received bytes
	HttpRecBytes uint

	// Received requests
	HttpReqCount uint
}

// Retrieves data
func (a *App) GetData() glob.J {
	a.MuxData.Lock()
	defer a.MuxData.Unlock()

	return a.Data
}

// Retrieves onion service instance
func (a *App) GetOnion() *tor.OnionService {
	a.MuxOnion.Lock()
	defer a.MuxOnion.Unlock()

	return a.Onion
}
