package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var version = "undefined"

func main() {
	app := cli.NewApp()
	app.Name = "slc"
	app.Version = version
	app.Usage = "Super Lightweight Container"

	app.Commands = []cli.Command{
		{
			Name:      "pull",
			Usage:     "pull container image",
			ArgsUsage: "[image]",
			Action:    Pull,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "state-path",
					Value: "./state",
				},
			},
		},
		{
			Name:      "run",
			Usage:     "run command on container image",
			ArgsUsage: "[image] [command]...",
			Action:    Run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "state-path",
					Value: "./state",
				},
			},
		}, {
			Name:      "rm",
			Usage:     "remove container image",
			ArgsUsage: "[image]",
			Action:    Remove,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "state-path",
					Value: "./state",
				},
			},
		}, {
			Name:   "clear",
			Usage:  "clear all state of the slc",
			Action: Clear,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "state-path",
					Value: "./state",
				},
			},
		},
	}

	log.SetFlags(log.Llongfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
