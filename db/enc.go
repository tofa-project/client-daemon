package db

import (
	"encoding/json"

	aes256gcm "github.com/tgbv/aes-256-gcm"
	"github.com/tofa-project/client-daemon/glob"
)

// Encrypts db related data
func Encrypt(data j) string {
	dataBytes, _ := json.Marshal(data)

	return aes256gcm.Encrypt([]byte(glob.V_DB_KEY), &dataBytes)
}

// Decrypts db related data
func Decrypt(data string) j {
	dataJson := aes256gcm.Decrypt([]byte(glob.V_DB_KEY), data)
	dataMap := make(j, 0)

	json.Unmarshal(dataJson, &dataMap)

	return dataMap
}
