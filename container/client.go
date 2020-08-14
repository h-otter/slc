package container

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
)

type HostMountOption struct {
	Src   string
	Flags uintptr
	Type  string
}

var DefaultHostMounts = []HostMountOption{
	{
		Src:   "/proc",
		Flags: 0,
		Type:  "proc",
	},
	{
		Src:   "/dev",
		Flags: syscall.MS_BIND | syscall.MS_PRIVATE,
		Type:  "dev",
	},
	{
		Src:   "/sys",
		Flags: syscall.MS_BIND | syscall.MS_PRIVATE,
		Type:  "sys",
	},
	{
		Src:   "/etc/resolv.conf",
		Flags: syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_PRIVATE,
		Type:  "",
	},
	{
		Src:   "/etc/passwd",
		Flags: syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_PRIVATE,
		Type:  "",
	},
	{
		Src:   "/etc/group",
		Flags: syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_PRIVATE,
		Type:  "",
	},
	{
		Src:   "/etc/hostname",
		Flags: syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_PRIVATE,
		Type:  "",
	},
	{
		Src:   "/etc/hosts",
		Flags: syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_PRIVATE,
		Type:  "",
	},
	{
		Src:   "/var/run",
		Flags: syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_PRIVATE,
		Type:  "",
	},
}

type SLCClient struct {
	stateDir string

	hostMounts []HostMountOption
}

func NewClient(stateDir string) (*SLCClient, error) {
	c := &SLCClient{
		hostMounts: DefaultHostMounts,
	}

	if filepath.IsAbs(stateDir) {
		c.stateDir = stateDir
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrapf(err, "os.Getwd()")
		}

		c.stateDir = filepath.Join(wd, stateDir)
	}

	if err := os.MkdirAll(c.stateDir, 0755); err != nil {
		return nil, errors.Wrapf(err, "os.MkdirAll(%s)", c.stateDir)
	}

	return c, nil
}

func (c *SLCClient) GetImageDir(image string) string {
	return filepath.Join(c.stateDir, "containers", image)
}

func CheckRoot() bool {
	output, err := exec.Command("id", "-u").Output()
	if err != nil {
		return false
		// return errors.Wrapf(err, "exec.Command(%s, %s).Output()", "id", "-u")
	}

	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		panic(err)
		// return errors.Wrapf(err, "strconv.Atoi(%s)", string(output[:len(output)-1]))
	}

	return i == 0
}
