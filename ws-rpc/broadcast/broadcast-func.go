// Utility packages dependant of wsrpc
package broadcast

import (
	"github.com/tofa-project/client-daemon/glob"
	clients "github.com/tofa-project/client-daemon/ws-rpc/clients"
	wconn "github.com/tofa-project/client-daemon/ws-rpc/wrapped-connection"
)

// Sends a JSON message for each client using SendViaPipe(...)
func Broadcast(data glob.J, f func(incomingData glob.J)) {
	clients.Clients.Range(func(key, wc interface{}) bool {
		wc.(*wconn.WrapppedConn).SendViaPipe(data, f)

		return true
	})
}

// Sends a JSON message for each client without waiting for response
func BroadcastMessage(data glob.J) {
	clients.Clients.Range(func(key, wc interface{}) bool {
		wc.(*wconn.WrapppedConn).Send(data)

		return true
	})
}
