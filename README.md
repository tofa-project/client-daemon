## Client Daemon
Manages the apps, Tor process, and how external services communicate with Client amid 2FA related activities. Written in Go. 

- Uses [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) for local database
- Database is encrypted using [tgbv/aes-256-gcm](https://github.com/tgbv/aes-256-gcm)
- Uses [gorilla/websocket](https://github.com/gorilla/websocket) for async RPC control/events.
- Uses [cretz/bine](https://github.com/cretz/bine) to control Tor process.

Each package is initiated procedurally in `main.go` file.

## Building
- You need `go 1.15`
- You need a GCC 32bit compiler to compile modules of  mattn/go-sqlite3
- `go build main.go`

Successfully compiled it on Ubuntu 20 x64 and Windows 10 x64. 

## Appendix
Daemon is useless without a websocket RPC client which answers incoming events. For that, [tofa-project/client-gui](https://github.com/tofa-project/client-gui) was born. 
