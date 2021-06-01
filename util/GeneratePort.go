package util

import (
	"math/rand"
	"strconv"

	"github.com/tofa-project/client-daemon/db"
)

// generates unique service port
func GeneratePort() string {
	port := strconv.Itoa(49152 + rand.Intn(65535-49152))

	for _, a := range db.GetApps() {
		if a["port"].(string) == port {
			return GeneratePort()
		}
	}

	return port
}
