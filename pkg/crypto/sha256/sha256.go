// +build !openssl static_build

package sha256

import (
	"crypto/sha256"
	"hash"
)

// New is a proxy to golang crypto New
func New() hash.Hash {
	return sha256.New()
}
