// +build !openssl static_build

package sha1

import (
	"crypto/sha1"
	"hash"
)

// New is a proxy to golang crypto New
func New() hash.Hash {
	return sha1.New()
}
