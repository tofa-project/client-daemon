package create_app

import (
	"encoding/base64"
	"net/url"
	"time"

	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
	"github.com/tofa-project/client-daemon/lib"
	wrap_data "github.com/tofa-project/client-daemon/onion-serv/app-related/wrap-data"
	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
	"github.com/tofa-project/client-daemon/util"
)

// Creates encapsulated app and returns its instance
func CreateApp(name string) *osat.App {
	data := glob.J{
		"name":          name,
		"priv_key":      base64.StdEncoding.EncodeToString(lib.GeneratePrivKey()),
		"port":          util.GeneratePort(),
		"is_registered": "0",
		"created_at":    time.Now().String(),
		"path":          url.QueryEscape(lib.GenerateRandomString(glob.C_APP_PATH_LENGTH)),
		"version":       "0",
	}

	appID := db.MakeApp(data)
	data["id"] = appID

	encapsulatedApp := wrap_data.WrapAppData(data)

	return encapsulatedApp
}
