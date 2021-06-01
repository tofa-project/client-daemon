// Responsible for interacting with the SQLite3 'data.db'.
//
// Does not follow ORM style. Contains just optimized methods for fast interaction with database.
package db

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tofa-project/client-daemon/glob"
)

// json map type (map[string]interface{})
type j = glob.J

// One instance to rule them all
var Instance *sql.DB

// Initializes the database system
func Init() {
	log.Print("Loading db... ")

	d, err := sql.Open("sqlite3", filepath.FromSlash(glob.V_DATA_DIR+glob.C_DB_DIR))
	if err != nil {
		panic(err)
	}
	Instance = d

	// Checks for default tables to be present
	r := Query(` SELECT name  FROM sqlite_master  WHERE type='table' AND name='apps'; `)
	if !r.Next() {
		makeAppsTable()
	}
	r.Close()

	r = Query(` SELECT name  FROM sqlite_master  WHERE type='table' AND name='logs'; `)
	if !r.Next() {
		makeLogsTable()
	}
	r.Close()

	log.Print("OK\n")
}

// Default query method. Returns only the rows. Handles errors automatically
func Query(q string, args ...interface{}) *sql.Rows {
	r, err := Instance.Query(q, args...)
	if err != nil {
		panic(err)
	}

	return r
}

// Default exec method. Handles errors automatically
func Exec(q string, args ...interface{}) {
	_, err := Instance.Exec(q, args...)
	if err != nil {
		panic(err)
	}
}
