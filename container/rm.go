package container

import (
	"os"

	"github.com/pkg/errors"
)

func (c *SLCClient) Remove(image string) error {
	if !CheckRoot() {
		return errors.New("root permission is required")
	}

	if err := os.RemoveAll(c.GetImageDir(image)); err != nil {
		return errors.Wrapf(err, "os.RemoveAll(%s)", c.GetImageDir(image))
	}

	return nil
}
