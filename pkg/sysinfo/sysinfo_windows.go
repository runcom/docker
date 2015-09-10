package sysinfo

// SysInfo holds fields specific to the Windows implementation. See
// CommonSysInfo for standard fields common to all platforms.
type SysInfo struct {
	CommonSysInfo

	// Fields below here are platform specific.
}

// New returns an empty SysInfo for windows for now.
func New(quiet bool) *SysInfo {
	sysInfo := &SysInfo{}
	return sysInfo
}
