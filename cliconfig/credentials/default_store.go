package credentials

import (
	"os/exec"

	"github.com/docker/docker/cliconfig"
)

// DetectDefaultStore sets the default credentials store
// if the host includes the default store helper program.
func DetectDefaultStore(c *cliconfig.ConfigFile) {
	if c.CredentialsStore != nil {
		// user defined
		return
	}

	if defaultCredentialsStore != nil {
		if _, err := exec.LookPath(remoteCredentialsPrefix + defaultCredentialsStore.Name); err == nil {
			c.CredentialsStore = defaultCredentialsStore
		}
	}
}
