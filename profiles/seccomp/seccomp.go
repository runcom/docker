// +build linux

package seccomp

import (
	"encoding/json"
	"fmt"

	"github.com/docker/engine-api/types"
	"github.com/opencontainers/specs/specs-go"
	libseccomp "github.com/seccomp/libseccomp-golang"
)

//go:generate go run -tags 'seccomp' generate.go

// GetDefaultProfile returns the default seccomp profile.
func GetDefaultProfile(rs *specs.Spec) (*specs.Seccomp, error) {
	return setupSeccomp(DefaultProfile(rs), rs)
}

// LoadProfile takes a file path and decodes the seccomp profile.
func LoadProfile(body string, rs *specs.Spec) (*specs.Seccomp, error) {
	var config types.Seccomp
	if err := json.Unmarshal([]byte(body), &config); err != nil {
		return nil, fmt.Errorf("Decoding seccomp profile failed: %v", err)
	}

	return setupSeccomp(&config, rs)
}

func setupSeccomp(config *types.Seccomp, rs *specs.Spec) (newConfig *specs.Seccomp, err error) {
	if config == nil {
		return nil, nil
	}

	// No default action specified, no syscalls listed, assume seccomp disabled
	if config.DefaultAction == "" && len(config.Syscalls) == 0 {
		return nil, nil
	}

	newConfig = &specs.Seccomp{}

	// if config.Architectures == 0 then libseccomp will figure out the architecture to use
	if len(config.Architectures) > 0 {
		for _, arch := range config.Architectures {
			newConfig.Architectures = append(newConfig.Architectures, specs.Arch(arch))
		}
	}

	newConfig.DefaultAction = specs.Action(config.DefaultAction)

	var arch string
	native, err := libseccomp.GetNativeArch()
	if err == nil {
		arch = native.String()
	}

	// Loop through all syscall blocks and convert them to libcontainer format
	for _, call := range config.Syscalls {
		var skip bool
		var cap string
		for _, a := range call.Includes.Arches {
			if a == arch {
				break
			}
			skip = true
		}
		if skip {
			continue
		}
		skip = false
		for _, c := range call.Includes.Caps {
			for _, cap = range rs.Process.Capabilities {
				if cap == c {
					break
				}
			}
			skip = true
		}
		if skip {
			continue
		}
		newCall := specs.Syscall{
			Name:   call.Name,
			Action: specs.Action(call.Action),
		}

		// Loop through all the arguments of the syscall and convert them
		for _, arg := range call.Args {
			newArg := specs.Arg{
				Index:    arg.Index,
				Value:    arg.Value,
				ValueTwo: arg.ValueTwo,
				Op:       specs.Operator(arg.Op),
			}

			newCall.Args = append(newCall.Args, newArg)
		}

		newConfig.Syscalls = append(newConfig.Syscalls, newCall)
	}

	return newConfig, nil
}
