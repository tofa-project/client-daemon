package lib

import (
	"crypto/ed25519"
	"math/rand"
	"time"
)

// Generates a new private key for onion service
func GeneratePrivKey() ed25519.PrivateKey {
	keyBytes := make([]byte, 32)

	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(keyBytes)
	rand.Seed(time.Now().UTC().UnixNano())

	privKey := ed25519.NewKeyFromSeed(keyBytes)

	return privKey
}
