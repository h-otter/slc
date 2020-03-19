package container

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

type SLCClient struct {
	stateDir string
}

func NewClient(stateDir string) (*SLCClient, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "os.Getwd()")
	}

	c := &SLCClient{
		stateDir: filepath.Join(wd, stateDir),
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
