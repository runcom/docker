package systemd

import (
	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/util"
)

// Send a message to the init daemon. It is common to ignore the error.
func SdNotify(state string) error {
	return daemon.SdNotify(state)
}

// SdBooted checks whether the host was booted with systemd
func SdBooted() bool {
	return util.IsRunningSystemd()
}
