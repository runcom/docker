// +build linux,cgo,!static_build,openssl

package rand

import (
	"io"
	"math/big"

	sslrand "github.com/shanemhansen/gossl/crypto/rand"
)

// Reader is a proxy to openssl crypto Reader
var Reader = sslrand.Reader

// Int is a proxy to openssl crypto Int
func Int(rand io.Reader, max *big.Int) (n *big.Int, err error) {
	return sslrand.Int(rand, max)
}
