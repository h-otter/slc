package container

import (
	"context"
	"path/filepath"

	"github.com/h-otter/slc/container/unpack"
	"github.com/pkg/errors"
)

func (c *SLCClient) Pull(image string) error {
	// sudoでunpackしないとcapabilityが展開されないなどの問題が起こる
	if !CheckRoot() {
		return errors.New("root permission is required")
	}

	target := filepath.Join(c.GetImageDir(image), "rootfs")

	if err := unpack.UnpackImage(context.Background(), image, target); err != nil {
		return errors.Wrap(err, "UnpackImage()")
	}

	if err := PrepareMountTargets(target, DefaultHostMounts); err != nil {
		return errors.Wrapf(err, "PrepareMountTargets(%s)", target)
	}

	return nil
}
