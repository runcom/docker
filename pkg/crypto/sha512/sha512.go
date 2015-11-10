// +build !openssl static_build

package sha512

import (
	"crypto/sha512"
	"hash"
)

// New is a proxy to golang crypto New
func New() hash.Hash {
	return sha512.New()
}
