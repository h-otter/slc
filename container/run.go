package container

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const OLD_ROOTFS = "oldrootfs"

func (c *SLCClient) Run(image string, argv []string) error {
	rootfs := filepath.Join(c.GetImageDir(image), "rootfs")

	if f, err := os.Stat(rootfs); os.IsNotExist(err) || !f.IsDir() {
		return errors.New("image is not found, please pull the image")
	}

	if err := syscall.Unshare(syscall.CLONE_FS | syscall.CLONE_NEWNS); err != nil {
		return errors.Wrap(err, "syscall.Unshare()")
	}

	// https://stackoverflow.com/questions/41561136/unshare-mount-namespace-not-working-as-expected
	if err := syscall.Mount("none", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return errors.Wrapf(err, "syscall.Mount(%s, %s, %s, %v, %s)", "none", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")
	}

	for _, v := range c.MountOptions {
		dst := filepath.Join(rootfs, v.Src)
		if err := syscall.Mount(v.Src, dst, v.Type, v.Flags, ""); err != nil {
			return errors.Wrapf(err, "syscall.Mount(%s, %s, %s, %v, %s)", v.Src, dst, v.Type, v.Flags, "")
		}

		if v.Flags&syscall.MS_RDONLY != 0 {
			flag := v.Flags | syscall.MS_REMOUNT
			if err := syscall.Mount(v.Src, dst, v.Type, flag, ""); err != nil {
				return errors.Wrapf(err, "syscall.Mount(%s, %s, %s, %v, %s)", v.Src, dst, v.Type, flag, "")
			}
		}
	}

	if err := syscall.Mount(rootfs, rootfs, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return errors.Wrapf(err, "syscall.Mount(%s, %s, %v, %s)", rootfs, rootfs, syscall.MS_BIND|syscall.MS_REC|syscall.MS_RDONLY, "")
	}
	// アプリケーションによってはroでマウントするのは良くないかもしれない
	// 本当は overlayfs などを使いたい
	// if err := syscall.Mount("none", rootfs, "", syscall.MS_BIND|syscall.MS_REC|syscall.MS_RDONLY|syscall.MS_REMOUNT, ""); err != nil {
	// 	return errors.Wrapf(err, "syscall.Mount(%s, %s, %v, %s)", rootfs, rootfs, syscall.MS_BIND|syscall.MS_REC|syscall.MS_RDONLY, "")
	// }
	if err := syscall.PivotRoot(rootfs, filepath.Join(rootfs, OLD_ROOTFS)); err != nil {
		return errors.Wrap(err, "syscall.PivotRoot()")
	}
	// TODO: mount /tmp

	// if err := syscall.Exec(argv[0], argv, os.Environ()); err != nil {
	// 	return errors.Wrapf(err, "syscall.Exec(%v, %v)", argv, os.Environ())
	// }

	cmd := exec.Command("/bin/sh", "-c", strings.Join(argv, " "))
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				os.Exit(s.ExitStatus())
			} else {
				return errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus.")
			}
		} else {
			return errors.Wrap(err, "cmd.Run()")
		}
	}

	return nil
}
