// +build linux,cgo,!static_build,openssl

package main

import (
	"net/http"

	"github.com/runcom/sslrt/sslrt"
)

func init() {
	http.DefaultTransport = sslrt.NewOpenSSLTransport(nil)
}
