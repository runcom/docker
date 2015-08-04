package sysinfo

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/opencontainers/runc/libcontainer/cgroups"
)

// New returns a new SysInfo, using the filesystem to detect which features the kernel supports.
func New(quiet bool) *SysInfo {
	sysInfo := &SysInfo{}
	sysInfo.cgroupMemInfo = checkCgroupMem(quiet)
	sysInfo.cgroupCpuInfo = checkCgroupCpu(quiet)
	sysInfo.cgroupBlkioInfo = checkCgroupBlkioInfo(quiet)

	_, err := cgroups.FindCgroupMountpoint("devices")
	sysInfo.CgroupDevicesEnabled = err == nil

	sysInfo.IPv4ForwardingDisabled = !readProcBool("/proc/sys/net/ipv4/ip_forward")
	sysInfo.BridgeNfCallIptablesDisabled = !readProcBool("/proc/sys/net/bridge/bridge-nf-call-iptables")
	sysInfo.BridgeNfCallIp6tablesDisabled = !readProcBool("/proc/sys/net/bridge/bridge-nf-call-ip6tables")

	// Check if AppArmor is supported.
	if _, err := os.Stat("/sys/kernel/security/apparmor"); !os.IsNotExist(err) {
		sysInfo.AppArmor = true
	}

	return sysInfo
}

func checkCgroupMem(quiet bool) *cgroupMemInfo {
	info := &cgroupMemInfo{}
	mountPoint, err := cgroups.FindCgroupMountpoint("memory")
	if err != nil {
		if !quiet {
			logrus.Warnf("Your kernel does not support cgroup memory limit: %v", err)
		}
		return info
	}
	info.MemoryLimit = true

	info.SwapLimit = cgroupEnabled(mountPoint, "memory.memsw.limit_in_bytes")
	if !quiet && !info.SwapLimit {
		logrus.Warn("Your kernel does not support swap memory limit.")
	}
	info.OomKillDisable = cgroupEnabled(mountPoint, "memory.oom_control")
	if !quiet && !info.OomKillDisable {
		logrus.Warnf("Your kernel does not support oom control.")
	}
	info.MemorySwappiness = cgroupEnabled(mountPoint, "memory.swappiness")
	if !quiet && !info.MemorySwappiness {
		logrus.Warnf("Your kernel does not support memory swappiness.")
	}

	return info
}

func checkCgroupCpu(quiet bool) *cgroupCpuInfo {
	info := &cgroupCpuInfo{}
	mountPoint, err := cgroups.FindCgroupMountpoint("cpu")
	if err != nil {
		if !quiet {
			logrus.Warn(err)
		}
		return info
	}

	info.CpuCfsPeriod = cgroupEnabled(mountPoint, "cpu.cfs_period_us")
	if !quiet && !info.CpuCfsPeriod {
		logrus.Warn("Your kernel does not support cgroup cfs period")
	}

	info.CpuCfsQuota = cgroupEnabled(mountPoint, "cpu.cfs_quota_us")
	if !quiet && !info.CpuCfsQuota {
		logrus.Warn("Your kernel does not support cgroup cfs quotas")
	}

	info.CpuShares = cgroupEnabled(mountPoint, "cpu.shares")
	if !quiet && !info.CpuShares {
		logrus.Warn("Your kernel does not support cgroup cpu shares")
	}
	return info
}

func checkCgroupBlkioInfo(quiet bool) *cgroupBlkioInfo {
	info := &cgroupBlkioInfo{}
	mountPoint, err := cgroups.FindCgroupMountpoint("blkio")
	if err != nil {
		if !quiet {
			logrus.Warn(err)
		}
		return info
	}

	info.BlkioWeight = cgroupEnabled(mountPoint, "blkio.weight")
	if !quiet && !info.BlkioWeight {
		logrus.Warn("Your kernel does not support cgroup blkio weight")
	}
	return info
}

func checkCgroupCpusetInfo(quiet bool) *cgroupCpusetInfo {
	info := &cgroupCpusetInfo{}
	_, err := cgroups.FindCgroupMountpoint("cpuset")
	if err != nil {
		if !quiet {
			logrus.Warn(err)
		}
		return info
	}

	info.Cpuset = true
	return info
}

func cgroupEnabled(mountPoint, name string) bool {
	_, err := os.Stat(path.Join(mountPoint, name))
	return err == nil
}

func readProcBool(path string) bool {
	val, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(val)) == "1"
}
