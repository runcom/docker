// +build linux

package syslog

import (
	"fmt"
	"io"
	"log/syslog"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
)

const name = "syslog"

type Syslog struct {
	writer *syslog.Writer
}

func init() {
	if err := logger.RegisterLogDriver(name, New); err != nil {
		logrus.Fatal(err)
	}
}

func New(ctx logger.Context) (logger.Logger, error) {
	log, err := syslog.New(syslog.LOG_DAEMON, fmt.Sprintf("%s/%s", path.Base(os.Args[0]), getTag(ctx)))
	if err != nil {
		return nil, err
	}

	return &Syslog{
		writer: log,
	}, nil
}

func (s *Syslog) Log(msg *logger.Message) error {
	if msg.Source == "stderr" {
		return s.writer.Err(string(msg.Line))
	}
	return s.writer.Info(string(msg.Line))
}

func (s *Syslog) Close() error {
	return s.writer.Close()
}

func (s *Syslog) Name() string {
	return name
}

func (s *Syslog) GetReader() (io.Reader, error) {
	return nil, logger.ReadLogsNotSupported
}

func getTag(ctx logger.Context) string {
	switch ctx.Config["syslog-tag"] {
	case "name":
		return ctx.ContainerID[:12]
	case "":
	case "id":
		return ctx.ContainerName
	}
	return ctx.ContainerID[:12]
}
