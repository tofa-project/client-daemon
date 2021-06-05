package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tofa-project/client-daemon/appdata"
	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
	oserv "github.com/tofa-project/client-daemon/onion-serv"
	tserv "github.com/tofa-project/client-daemon/tor-serv"
	ws_rpc "github.com/tofa-project/client-daemon/ws-rpc/http"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// gracefully shutdown
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt)
		<-sc

		log.Println("daemon quitting...")

		db.Instance.Close()
		tserv.Instance.Close()

		wg.Done()
	}()

	glob.Init()

	reader := bufio.NewReader(os.Stdin)

	// read db key from stdin
	log.Print("Input db key: ")
	dbKey, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	glob.V_DB_KEY = strings.Trim(dbKey, "\n")

	// read ws key from stdin
	log.Print("Input wsrpc key: ")
	wsKey, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	glob.V_WS_RPC_HEADER_SECRET = strings.Trim(wsKey, "\n")

	ws_rpc.Init()

	appdata.Init()

	db.Init()

	// attempt to retrieve data from database
	// this will crash the daemon if password is incorrect
	db.GetApps()

	tserv.Init()

	oserv.Init()

	wg.Wait()
}
