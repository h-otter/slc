package container

import (
	"context"
	"os"
	"path/filepath"

	"github.com/h-otter/slc/container/unpack"
	"github.com/pkg/errors"
)

func createMountTarget(src, dst string, ignoreNoSourceError bool) error {
	if _, err := os.Stat(dst); err != nil {
		srcStat, err := os.Stat(src)
		if err != nil {
			if ignoreNoSourceError {
				return nil
			} else {
				return errors.Wrapf(err, "os.Stat(%s)", src)
			}
		}

		if srcStat.IsDir() {
			if err := os.MkdirAll(dst, 0755); err != nil {
				return errors.Wrapf(err, "os.MkdirAll(%s)", dst)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
				return errors.Wrapf(err, "os.MkdirAll(%s)", filepath.Base(dst))
			}

			if _, err := os.Create(dst); err != nil {
				return errors.Wrapf(err, "os.Create(%s)", dst)
			}
		}
	}

	return nil
}

// TODO: root で実行する必要があるため、pullの段階で準備する
func PrepareMountTargets(rootfs string, options []HostMountOption) error {
	for _, v := range options {
		dst := filepath.Join(rootfs, v.Src)
		if err := createMountTarget(v.Src, dst, v.IgnoreNoSourceError); err != nil {
			return errors.Wrapf(err, "createMountTarget(%s, %s, %v)", v.Src, dst, v.IgnoreNoSourceError)
		}
	}

	dst := filepath.Join(rootfs, OLD_ROOTFS)
	if err := createMountTarget("/", dst, false); err != nil {
		return errors.Wrapf(err, "createMountTarget(%s, %s, %v)", "/", dst, false)
	}

	return nil
}

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
