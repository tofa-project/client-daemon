package lib

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
// Credits: https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb#file-main-go-L46
func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz\\|]}[{'\";:/?.>,<=+-_)(*&^%$#@!~`"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic("Random string generator failed! What OS are you using?")
			//return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}
