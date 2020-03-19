package container

import (
	"os"

	"github.com/pkg/errors"
)

func (c *SLCClient) Clear() error {
	if !CheckRoot() {
		return errors.New("root permission is required")
	}

	if err := os.RemoveAll(c.stateDir); err != nil {
		return errors.Wrapf(err, "os.RemoveAll(%s)", c.stateDir)
	}

	return nil
}
