// Holds the global tor service instance
package tserv

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/cretz/bine/tor"
	"github.com/tofa-project/client-daemon/glob"
)

var Instance *tor.Tor

// Initializes Tor daemon
func Init() {
	log.Print("Starting Tor daemon...")

	tInst, err := tor.Start(context.Background(),
		&tor.StartConf{
			DebugWriter: os.Stderr,
			ExePath:     filepath.FromSlash(glob.V_TOR_BIN_PATH),
			DataDir:     filepath.FromSlash(glob.V_DATA_DIR + glob.C_TOR_DIR),
			TorrcFile:   filepath.FromSlash(glob.V_DATA_DIR + glob.C_TOR_DIR + "/torrc"),
		},
	)
	if err != nil {
		panic("FAIL " + err.Error())
	}
	Instance = tInst

	log.Print("OK")
}
