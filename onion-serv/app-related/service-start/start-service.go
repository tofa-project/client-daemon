package service_start

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/tofa-project/client-daemon/glob"

	osat "github.com/tofa-project/client-daemon/onion-serv/app-type"
	tserv "github.com/tofa-project/client-daemon/tor-serv"
	wsrpc_broadcast "github.com/tofa-project/client-daemon/ws-rpc/broadcast"
	wsrpc_events "github.com/tofa-project/client-daemon/ws-rpc/events-list"
)

// Resets req count and rec bytes as long as onion service is running
func resetRecVars(a *osat.App) {
	if !a.Running {
		return
	}

	a.HttpRecBytes = 0
	a.HttpReqCount = 0

	time.Sleep(time.Second * glob.C_ONION_SERV_HTTP_REC_LIMITS_REFRESH_INTERVAL)

	go resetRecVars(a)
}

// Starts onion service for app
func StartService(a *osat.App) error {
	if a.Running {
		return fmt.Errorf("service already running")
	}

	ReloadHttp(a)

	data := a.Data

	log.Print("\tPublishing onion service [", data["name"].(string), "]...")

	port, _ := strconv.Atoi(data["port"].(string))

	// key is base64 encoded; decode it first, then do the stuff
	//
	keyEncoded := data["priv_key"].(string)
	keyDecoded, err := base64.StdEncoding.DecodeString(keyEncoded)
	if err != nil {
		log.Panicf("\tFAIL! %s", err)
		return err
	}
	keyFinal := ed25519.PrivateKey(keyDecoded)

	// Create onion service
	//
	oService, err := tserv.Instance.Listen(context.TODO(), &tor.ListenConf{
		RemotePorts: []int{port},
		Version3:    true,
		Key:         keyFinal,
	})
	if err != nil {
		log.Panicf("\tFAIL! %s", err)
		return err
	}
	a.Onion = oService
	a.Running = true

	// start http server
	//
	sendErr := true
	eCh := make(chan error)
	defer close(eCh)
	go func() {
		e := a.HttpServer.Serve(oService)
		if sendErr {
			eCh <- e
			sendErr = false
		}
	}()

	// handles timeout. must occur if http starts as expected
	go func() {
		time.Sleep(glob.C_HTTP_SERVE_TIMEOUT * time.Millisecond)

		if sendErr {
			eCh <- nil
			sendErr = false
		}
	}()

	tErr := <-eCh
	if tErr == nil {
		log.Printf("\tOK [%s.onion:%s]", oService.ID, data["port"])

		go resetRecVars(a)

		// emit publishing state event to GUI
		wsrpc_broadcast.BroadcastMessage(glob.J{
			"type":      "event",
			"eventName": wsrpc_events.NFO_APP_PUBLISHED,
			"appID":     data["id"].(string),
		})
	} else {
		log.Printf("\tCould not start HTTP server for service! %s\n", tErr)
	}
	return tErr
}
