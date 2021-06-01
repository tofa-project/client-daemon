package wrap_data

import (
	"github.com/tofa-project/client-daemon/glob"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
)

// Wraps app data to an App structure
func WrapAppData(data glob.J) *osat.App {
	a := &osat.App{Data: data}

	return a
}
