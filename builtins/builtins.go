package builtins

import (
	"runtime"

	"github.com/docker/docker/api"
	"github.com/docker/docker/autogen/dockerversion"
	"github.com/docker/docker/pkg/parsers/kernel"
)

type Version struct {
	Version       string
	ApiVersion    string
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string `json:",omitempty"`
}

// builtins jobs independent of any subsystem
func DockerVersion() *Version {
	v := &Version{
		Version:    dockerversion.VERSION,
		ApiVersion: api.APIVERSION,
		GitCommit:  dockerversion.GITCOMMIT,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}
	if kernelVersion, err := kernel.GetKernelVersion(); err == nil {
		v.KernelVersion = kernelVersion.String()
	}

	return v
}
