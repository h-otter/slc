package main

import (
	"github.com/h-otter/slc/container"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func Pull(cctx *cli.Context) error {
	state := cctx.String("state-path")
	if cctx.NArg() != 1 {
		return errors.Errorf("want 1 arg, have %d args", cctx.NArg())
	}
	image := cctx.Args().Get(0)

	c, err := container.NewClient(state)
	if err != nil {
		return err
	}

	if err := c.Pull(image); err != nil {
		return err
	}

	return nil
}

func Run(cctx *cli.Context) error {
	state := cctx.String("state-path")

	commands := []string{}
	if cctx.NArg() < 1 {
		return errors.Errorf("want more than 1 arg, have %d args", cctx.NArg())
	} else if cctx.NArg() == 1 {
		commands = []string{"/bin/sh"}
	} else {
		commands = cctx.Args().Tail()
	}
	image := cctx.Args().Get(0)

	c, err := container.NewClient(state)
	if err != nil {
		return err
	}

	if err := c.Run(image, commands); err != nil {
		return err
	}

	return nil
}

func Remove(cctx *cli.Context) error {
	state := cctx.String("state-path")
	if cctx.NArg() != 1 {
		return errors.Errorf("want 1 arg, have %d args", cctx.NArg())
	}
	image := cctx.Args().Get(0)

	c, err := container.NewClient(state)
	if err != nil {
		return err
	}

	if err := c.Remove(image); err != nil {
		return err
	}

	return nil
}

func Clear(cctx *cli.Context) error {
	state := cctx.String("state-path")
	if cctx.NArg() != 0 {
		return errors.Errorf("want 0 arg, have %d args", cctx.NArg())
	}

	c, err := container.NewClient(state)
	if err != nil {
		return err
	}

	if err := c.Clear(); err != nil {
		return err
	}

	return nil
}
