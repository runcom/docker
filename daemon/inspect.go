package daemon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/daemon/network"
	"github.com/docker/docker/engine"
	"github.com/docker/docker/runconfig"
)

// needed for backwards compatibility with api < 1.12
type ContainerJSONRaw struct {
	*Container
	HostConfig *runconfig.HostConfig
}

type ContainerJSON struct {
	Id              string
	Created         time.Time
	Path            string
	Args            []string
	Config          *runconfig.Config
	State           *State
	Image           string
	NetworkSettings *network.Settings
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Name            string
	RestartCount    int
	Driver          string
	ExecDriver      string
	MountLabel      string
	ProcessLabel    string
	Volumes         map[string]string
	VolumesRW       map[string]bool
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      *runconfig.HostConfig
}

func (daemon *Daemon) ContainerInspectRaw(name string) (*ContainerJSONRaw, error) {
	container, err := daemon.Get(name)
	if err != nil {
		return nil, err
	}

	container.Lock()
	defer container.Unlock()

	return &ContainerJSONRaw{container, container.hostConfig}, nil
}

func (daemon *Daemon) ContainerInspect(name string) (*ContainerJSON, error) {
	container, err := daemon.Get(name)
	if err != nil {
		return nil, err
	}

	container.Lock()
	defer container.Unlock()

	out := &ContainerJSON{
		Id:              container.ID,
		Created:         container.Created,
		Path:            container.Path,
		Args:            container.Args,
		Config:          container.Config,
		State:           container.State,
		Image:           container.ImageID,
		NetworkSettings: container.NetworkSettings,
		ResolvConfPath:  container.ResolvConfPath,
		HostnamePath:    container.HostnamePath,
		HostsPath:       container.HostsPath,
		LogPath:         container.LogPath,
		Name:            container.Name,
		RestartCount:    container.RestartCount,
		Driver:          container.Driver,
		ExecDriver:      container.ExecDriver,
		MountLabel:      container.MountLabel,
		ProcessLabel:    container.ProcessLabel,
		Volumes:         container.Volumes,
		VolumesRW:       container.VolumesRW,
		AppArmorProfile: container.AppArmorProfile,
		ExecIDs:         container.GetExecIDs(),
	}

	if children, err := daemon.Children(container.Name); err == nil {
		for linkAlias, child := range children {
			container.hostConfig.Links = append(container.hostConfig.Links, fmt.Sprintf("%s:%s", child.Name, linkAlias))
		}
	}
	// we need this trick to preserve empty log driver, so
	// container will use daemon defaults even if daemon change them
	if container.hostConfig.LogConfig.Type == "" {
		container.hostConfig.LogConfig = daemon.defaultLogConfig
		defer func() {
			container.hostConfig.LogConfig = runconfig.LogConfig{}
		}()
	}

	out.HostConfig = container.hostConfig

	container.hostConfig.Links = nil

	return out, nil
}

func (daemon *Daemon) ContainerExecInspect(job *engine.Job) error {
	if len(job.Args) != 1 {
		return fmt.Errorf("usage: %s ID", job.Name)
	}
	id := job.Args[0]
	eConfig, err := daemon.getExecConfig(id)
	if err != nil {
		return err
	}

	b, err := json.Marshal(*eConfig)
	if err != nil {
		return err
	}
	job.Stdout.Write(b)
	return nil
}
