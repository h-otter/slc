package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

var version = "undefined"

func main() {
	app := cli.NewApp()
	app.Name = "slc"
	app.Version = version
	app.Usage = "Super Lightweight Container"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "state-path",
			Usage: "SLC store container images in the state directory.",
			Value: "/tmp/slc-state",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:      "pull",
			Usage:     "pull container image",
			ArgsUsage: "[image]",
			Action:    Pull,
		},
		{
			Name:      "run",
			Usage:     "run command on container image",
			ArgsUsage: "[image] [command]...",
			Action:    Run,
		},
		{
			Name:      "rm",
			Usage:     "remove container image",
			ArgsUsage: "[image]",
			Action:    Remove,
		},
		{
			Name:   "clear",
			Usage:  "clear all state of the slc",
			Action: Clear,
		},
	}

	log.SetFlags(log.Llongfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
