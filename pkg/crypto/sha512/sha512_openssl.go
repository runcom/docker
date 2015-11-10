// +build linux,cgo,!static_build,openssl

package sha512

import (
	"hash"

	"github.com/shanemhansen/gossl/crypto/sha512"
)

// New is a proxy to openssl crypto New
func New() hash.Hash {
	return sha512.New()
}
