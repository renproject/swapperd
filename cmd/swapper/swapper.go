package main

import (
	"errors"
	"log"
	"os"

	"github.com/republicprotocol/swapperd/driver/config"
	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/republicprotocol/swapperd/driver/swapper"
	"github.com/republicprotocol/swapperd/utils"
	"github.com/urfave/cli"
)

// Define flags for commands
var (
	locationFlag = cli.StringFlag{
		Name:  "location",
		Value: utils.GetDefaultSwapperHome(),
		Usage: "Home directory for RenEx Swapper",
	}
	networkFlag = cli.StringFlag{
		Name:  "network",
		Value: "mainnet",
		Usage: "name of the network",
	}
	keyPhraseFlag = cli.StringFlag{
		Name:  "keyphrase",
		Usage: "keyphrase to unlock the keystore file",
	}
	toFlag = cli.StringFlag{
		Name:  "to",
		Usage: "receiver address you want to withdraw the tokens to",
	}
	tokenFlag = cli.StringFlag{
		Name:  "token",
		Usage: "type of token you want to withdraw (options ETH, BTC)",
	}
	valueFlag = cli.Float64Flag{
		Name:  "value",
		Value: 0,
		Usage: "amount of token you want to withdraw (to withdraw 1 Eth " +
			"use --value 1 --token eth) If this flag is not set the entire " +
			"balance will be withdrawn",
	}
	portFlag = cli.Int64Flag{
		Name:  "port",
		Value: 18516,
		Usage: "port on which the http server runs,",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "RenEx Swapper CLI"
	app.Usage = ""

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "http",
			Usage: "start the RenEx swapper's http server",
			Flags: []cli.Flag{networkFlag, locationFlag, keyPhraseFlag, portFlag},
			Action: func(c *cli.Context) {
				swapper, err := initializeSwapper(c)
				if err != nil {
					panic(err)
				}
				swapper.Http(c.Int64("port"))
			},
		},
		{
			Name:  "withdraw",
			Usage: "withdraw the funds in the swapper accounts",
			Flags: []cli.Flag{networkFlag, keyPhraseFlag, locationFlag, toFlag, tokenFlag, valueFlag},
			Action: func(c *cli.Context) error {
				swapper, err := initializeSwapper(c)
				if err != nil {
					return err
				}
				return withdraw(c, swapper)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeSwapper(ctx *cli.Context) (swapper.Swapper, error) {
	network := ctx.String("network")
	keyPhrase := ctx.String("keyphrase")
	location := ctx.String("location")

	cfg, err := config.New(location, network)
	if err != nil {
		return nil, err
	}

	ks, err := keystore.LoadFromFile(cfg, keyPhrase)
	if err != nil {
		return nil, err
	}

	return swapper.NewSwapper(cfg, ks), nil
}

func withdraw(ctx *cli.Context, swapper swapper.Swapper) error {
	receiver := ctx.String("to")
	if receiver == "" {
		return errors.New("receiver address cannot be empty")
	}
	token := ctx.String("token")
	if token == "" {
		return errors.New("please enter a valid withdraw token")
	}
	return swapper.Withdraw(token, receiver, ctx.Float64("value"))
}
