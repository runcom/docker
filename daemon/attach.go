package daemon

import (
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/version"
)

type ContainerAttachWithLogsConfig struct {
	UseStdin, UseStdout, UseStderr bool
	InStream                       io.ReadCloser
	OutStream                      io.Writer
	Version                        version.Version
	Logs, Stream                   bool
}

func (daemon *Daemon) ContainerAttachWithLogs(name string, config *ContainerAttachWithLogsConfig) error {
	container, err := daemon.Get(name)
	if err != nil {
		return err
	}

	var errStream io.Writer

	if !container.Config.Tty && config.Version.GreaterThanOrEqualTo("1.6") {
		errStream = stdcopy.NewStdWriter(config.OutStream, stdcopy.Stderr)
		config.OutStream = stdcopy.NewStdWriter(config.OutStream, stdcopy.Stdout)
	} else {
		errStream = config.OutStream
	}

	var stdin io.ReadCloser
	var stdout, stderr io.Writer

	if config.UseStdin {
		stdin = config.InStream
	}
	if config.UseStdout {
		stdout = config.OutStream
	}
	if config.UseStderr {
		stderr = errStream
	}

	return container.AttachWithLogs(stdin, stdout, stderr, config.Logs, config.Stream)
}
