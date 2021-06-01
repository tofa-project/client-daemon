// Contains websocket rpc events list
package events_list

// Outbound event sent amid an information received from app server.
//
// Requires "appName", "infoDescription"
const O_INCOMING_INFO = "incoming_info"

// Outbound event sent amid an "action which requires confirmation"
// is received from app server.
//
// Requires "appName", "askDescription"
const O_INCOMING_ASK = "incoming_ask"

// Outbound event sent amid a registration request received from app server.
//
// Requires "appName", "appDescription"
const O_INCOMING_REGISTRATION = "incoming_registration"

// Sent once app is published to HSDIR
//
// Requires "appID"
const NFO_APP_PUBLISHED = "app_published"

// Sent once app is disconnected from Tor
//
// Requires "appID"
const NFO_APP_UNPUBLISHED = "app_unpublished"

// Sent amid potential DoS attack
//
// Requires "appID"
const NFO_APP_DOS = "app_under_dos"
