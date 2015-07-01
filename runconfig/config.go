package runconfig

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/docker/docker/pkg/nat"
)

// Entrypoint encapsulates the container entrypoint.
// It might be represented as a string or an array of strings.
// We need to override the json decoder to accept both options.
// The JSON decoder will fail if the api sends an string and
//  we try to decode it into an array of string.
type Entrypoint struct {
	List
}

func NewEntrypoint(parts ...string) *Entrypoint {
	return &Entrypoint{List{parts}}
}

type Command struct {
	List
}

func (e *Command) ToString() string {
	return strings.Join(e.Slice(), " ")
}

func NewCommand(parts ...string) *Command {
	return &Command{List{parts}}
}

// Note: the Config structure should hold only portable information about the container.
// Here, "portable" means "independent from the host we are running on".
// Non-portable information *should* appear in HostConfig.
type Config struct {
	Hostname        string
	Domainname      string
	User            string
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	ExposedPorts    map[nat.Port]struct{}
	PublishService  string
	Tty             bool // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool // Open stdin
	StdinOnce       bool // If true, close stdin after the 1 attached client disconnects.
	Env             []string
	Cmd             *Command
	Image           string // Name of the image as it was passed by the operator (eg. could be symbolic)
	Volumes         map[string]struct{}
	VolumeDriver    string
	WorkingDir      string
	Entrypoint      *Entrypoint
	NetworkDisabled bool
	MacAddress      string
	OnBuild         []string
	Labels          map[string]string
}

type ContainerConfigWrapper struct {
	*Config
	InnerHostConfig *HostConfig `json:"HostConfig,omitempty"`
	Cpuset          string      `json:",omitempty"` // Deprecated. Exported for backwards compatibility.
	*HostConfig                 // Deprecated. Exported to read attrubutes from json that are not in the inner host config structure.

}

func (w *ContainerConfigWrapper) GetHostConfig() *HostConfig {
	hc := w.HostConfig

	if hc == nil && w.InnerHostConfig != nil {
		hc = w.InnerHostConfig
	} else if w.InnerHostConfig != nil {
		if hc.Memory != 0 && w.InnerHostConfig.Memory == 0 {
			w.InnerHostConfig.Memory = hc.Memory
		}
		if hc.MemorySwap != 0 && w.InnerHostConfig.MemorySwap == 0 {
			w.InnerHostConfig.MemorySwap = hc.MemorySwap
		}
		if hc.CpuShares != 0 && w.InnerHostConfig.CpuShares == 0 {
			w.InnerHostConfig.CpuShares = hc.CpuShares
		}

		hc = w.InnerHostConfig
	}

	if hc != nil && w.Cpuset != "" && hc.CpusetCpus == "" {
		hc.CpusetCpus = w.Cpuset
	}

	return hc
}

// DecodeContainerConfig decodes a json encoded config into a ContainerConfigWrapper
// struct and returns both a Config and an HostConfig struct
// Be aware this function is not checking whether the resulted structs are nil,
// it's your business to do so
func DecodeContainerConfig(src io.Reader) (*Config, *HostConfig, error) {
	decoder := json.NewDecoder(src)

	var w ContainerConfigWrapper
	if err := decoder.Decode(&w); err != nil {
		return nil, nil, err
	}

	return w.Config, w.GetHostConfig(), nil
}
