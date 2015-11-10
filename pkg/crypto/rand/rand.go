// +build !openssl static_build

package rand

import (
	cryptorand "crypto/rand"
	"io"
	"math/big"
)

// Reader is a proxy to golang crypto Reader
var Reader = cryptorand.Reader

// Int is a proxy to golang crypto Int
func Int(rand io.Reader, max *big.Int) (n *big.Int, err error) {
	return cryptorand.Int(rand, max)
}
