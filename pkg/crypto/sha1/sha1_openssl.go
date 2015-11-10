// +build linux,cgo,!static_build,openssl

package sha1

import (
	"hash"

	"github.com/shanemhansen/gossl/crypto/sha1"
)

// New is a proxy to openssl crypto New
func New() hash.Hash {
	return sha1.New()
}
