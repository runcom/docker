// +build linux,cgo,!static_build,openssl

package md5

import (
	"hash"

	"github.com/shanemhansen/gossl/crypto/md5"
)

// New is a proxy to openssl crypto New implementation
func New() hash.Hash {
	return md5.New()
}
