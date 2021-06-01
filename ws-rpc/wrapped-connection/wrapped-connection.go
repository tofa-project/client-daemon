package wconn

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tofa-project/client-daemon/glob"
	"github.com/tofa-project/client-daemon/lib"
)

// A ws connection wrapped with afferent data structures
type WrapppedConn struct {
	ID   string
	Conn *websocket.Conn

	Pipes          sync.Map
	SendMessageMux sync.Mutex
}

// Sends a json message to this connection
func (wc *WrapppedConn) Send(data glob.J) {
	wc.SendMessageMux.Lock()
	defer wc.SendMessageMux.Unlock()

	wc.Conn.WriteJSON(data)
}

// Creates a virtual pipe.
// Sends a message via a pipe and pushes a callback invoked on message receivment on that pipe.
func (wc *WrapppedConn) SendViaPipe(data glob.J, callback func(glob.J)) {

	pipeID := lib.GenerateRandomString(glob.C_WS_PIPE_ID_LENGTH)

	// stores callback
	wc.Pipes.Store(pipeID, func(input glob.J) {
		if input["pipeID"].(string) == pipeID {
			go callback(input)

			wc.Pipes.Delete(pipeID)
		}
	})

	// sends the actual message
	data["pipeID"] = pipeID
	wc.Send(data)

	// handles timeout
	go func() {
		time.Sleep(glob.C_WS_PIPE_RES_TIMEOUT * time.Minute)

		wc.Pipes.Delete(pipeID)
	}()
}
