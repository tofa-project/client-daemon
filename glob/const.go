package glob

// Relative to V_DATA_DIR
const C_TOR_DIR = "/tor"

// Relative to V_DATA_DIR
const C_DB_DIR = "/data.db"

// websocket connection ID length
const C_WS_CONN_ID_LENGTH = 16

// websocket pipe ID length
const C_WS_PIPE_ID_LENGTH = 16

// websocket pipe response timeout after which pipes will be flushed from map
const C_WS_PIPE_RES_TIMEOUT = 5 // minutes

// Timeout if no error is sent to the error channel amid starting a HTTP server
const C_HTTP_SERVE_TIMEOUT = 500 // milliseconds

// Timeout used for onion services, if no ws-rpc client replies
const C_HTTP_RES_TIMEOUT = 1 // minutes

// application path length for onion service
const C_APP_PATH_LENGTH = 64

// max onion service http received bytes per timeframe
const C_MAX_ONION_SERV_HTTP_REC_BYTES = 1024 * 250 // 0.25MB

// onion service http received requests count interval in seconds
const C_MAX_ONION_SERV_HTTP_REC_REQ_COUNT = 30

// onion service http received bytes or requests count refresh timeframe interval in seconds
const C_ONION_SERV_HTTP_REC_LIMITS_REFRESH_INTERVAL = 60
