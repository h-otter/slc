package container

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

const imgCommandURL = "https://github.com/genuinetools/img/releases/download/v0.5.7/img-linux-amd64"

func Download(src, dst string, permission os.FileMode) error {
	resp, err := http.Get(src)
	if err != nil {
		return errors.Wrapf(err, "http.Get()")
	}
	defer resp.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "os.Create()")
	}
	defer out.Close()

	if err := os.Chmod(dst, permission); err != nil {
		return errors.Wrapf(err, "os.Chmod()")
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return errors.Wrapf(err, "io.Copy()")
	}

	return nil
}

func (c *SLCClient) Pull(image string) error {
	// sudoでunpackしないとcapabilityが展開されないなどの問題が起こる
	if !CheckRoot() {
		return errors.New("root permission is required")
	}

	imgCommand := filepath.Join(c.stateDir, "img")
	if _, err := os.Stat(imgCommand); err != nil {
		log.Printf("[INFO] Downloading the img command")
		if err := Download(imgCommandURL, imgCommand, 0755); err != nil {
			return errors.Wrapf(err, "Download(%s, %s, %v)", imgCommandURL, imgCommand, 0755)
		}
		log.Printf("[INFO] Download completed")
	}

	pullCommand := []string{imgCommand, "pull", "-s", c.stateDir, image}
	cmd := exec.Command(pullCommand[0], pullCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "exec.Command(%v).Run()", pullCommand)
	}

	unpackCommand := []string{imgCommand, "unpack", "-s", c.stateDir, image}
	cmd = exec.Command(unpackCommand[0], unpackCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// img -o option is not working
	cmd.Dir = c.GetImageDir(image)
	if err := os.MkdirAll(cmd.Dir, 0755); err != nil {
		return errors.Wrapf(err, "os.MkdirAll(%s)", cmd.Dir)
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "exec.Command(%v).Run()", unpackCommand)
	}

	rootfs := filepath.Join(c.GetImageDir(image), "rootfs")
	if err := PrepareMountTargets(rootfs, DefaultHostMounts); err != nil {
		return errors.Wrapf(err, "PrepareMountTargets(%s)", rootfs)
	}

	return nil
}
