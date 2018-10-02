package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)


// Define flags for commands
var (
	networkFlag = cli.StringFlag{
		Name:  "network",
		Usage: "name of the test network",
	}
)

func main() {
	app := cli.NewApp()

	// Define subcommands
	app.Commands = []cli.Command{
		{
			Name:  "http",
			Usage: "start running the swapper ",
			Action: func(c *cli.Context) error {
				panic("Implement the http logic here")
			},
		},
		{
			Name:  "withdraw",
			Usage: "withdraw the funds in the swapper accounts",
			Action: func(c *cli.Context) error {
				panic("Implement the withdraw logic here")
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

