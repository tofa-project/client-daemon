// Package responsible for holding global variables used across module.
//
// Some of them have implicit values, some are initialized within
// application thread based on needs.
package glob

// The data dir which contains the database
// and tofa related configuration
var V_DATA_DIR string

// Key used to encrypt/decrypt the database data
var V_DB_KEY string

// RPC port, host to listen to
var V_RPC_PORT int
var V_RPC_HOST string

// websocket rpc host to listen to
var V_WS_RPC_HOST string

// websocket secret header to prevent other malicious processes from connecting
var V_WS_RPC_HEADER_SECRET string
