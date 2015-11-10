// +build linux,cgo,!static_build,openssl

package sha256

import (
	"hash"

	"github.com/shanemhansen/gossl/crypto/sha256"
)

// New is a proxy to openssl crypto New
func New() hash.Hash {
	return sha256.New()
}
