package daemon

import (
	"time"

	"github.com/docker/docker/engine"
)

func (daemon *Daemon) ContainerWait(job *engine.Job) engine.Status {
	if len(job.Args) != 1 {
		return job.Errorf("Usage: %s", job.Name)
	}
	name := job.Args[0]
	autoRemove := job.GetenvBool("autoremove")
	container, err := daemon.Get(name)
	if err != nil {
		return job.Errorf("%s: %v", job.Name, err)
	}
	status, _ := container.WaitStop(-1 * time.Second)
	job.Printf("%d\n", status)
	if autoRemove {
		job.Printf("ciao %s", job.Name)
		//daemon.Rm(container)
	}
	return engine.StatusOK
}
