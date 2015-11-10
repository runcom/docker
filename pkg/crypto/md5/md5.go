// +build !openssl

package md5

import (
	"crypto/md5"
	"hash"
)

// New is a proxy to golang crypto New implementation
func New() hash.Hash {
	return md5.New()
}
