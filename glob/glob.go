package glob

import (
	"flag"
	"log"
	"os"
)

var flagsParserNFO = map[string]string{
	"rpc_host":   "Websocket RPC host to listen to",
	"rpc_secret": `Websocket RPC required secret header to prevent other malicious processes from connecting. Empty string disables it.`,
}

// Initializes variables with default values based on user input
func Init() {
	log.Printf("Loading glob... ")

	// setup flags
	//

	uhdp, _ := os.UserHomeDir()
	fV_DATA_DIR := flag.String("data-dir", uhdp+"/tofa", "Tofa's default data directory. By default: "+uhdp+"/tofa")
	fV_WS_RPC_HOST := flag.String("ws-rpc-host", "127.0.0.1:50750", flagsParserNFO["rpc_host"])
	fV_WS_RPC_HEADER_SECRET := flag.String("ws-secret", "", flagsParserNFO["rpc_secret"])

	flag.Parse()

	// assignments && additional checks
	//

	V_DATA_DIR = *fV_DATA_DIR
	V_WS_RPC_HOST = *fV_WS_RPC_HOST
	V_WS_RPC_HEADER_SECRET = *fV_WS_RPC_HEADER_SECRET

	log.Printf("OK\n")
}
