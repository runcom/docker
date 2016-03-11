package credentials

import (
	"github.com/docker/docker/cliconfig"
)

var defaultCredentialsStore = cliconfig.NewCredentialsStore("wincred", nil)
