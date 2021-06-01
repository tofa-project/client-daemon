// Package is responsible for initializing user %appdata% related files/folders.
// Default user data entities will be 'data.db' and '/tor'
//
// 1. data.db -> contains all registered 2FA apps and their related info.
// It's encrypted except their IDs which are used to be identified
//
// 2. /tor -> contains tor related files/folders
package appdata

import (
	"log"
	"os"
	"path/filepath"

	g "github.com/tofa-project/client-daemon/glob"
)

// Initializes the appdata default directory tree if not present
func Init() {
	log.Print("Loading appdata... ")

	// must have tor dir
	// this includes the parent g.V_DATA_DIR
	target := filepath.FromSlash(g.V_DATA_DIR + g.C_TOR_DIR)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			panic(err)
		}
	}

	// must have database file
	target = filepath.FromSlash(g.V_DATA_DIR + g.C_DB_DIR)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		f, err := os.Create(target)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	// must have torrc file
	target = filepath.FromSlash(g.V_DATA_DIR + g.C_TOR_DIR + "/torrc")
	if _, err := os.Stat(target); os.IsNotExist(err) {
		f, err := os.Create(target)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	log.Printf("OK\n")
}
