package clients

import (
	"sync"
)

// All active connections
var Clients = sync.Map{}
