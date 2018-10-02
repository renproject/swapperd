package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/driver/config"
	keystoreDriver "github.com/republicprotocol/renex-swapper-go/driver/keystore"
	"github.com/republicprotocol/renex-swapper-go/driver/swapper"
	"github.com/urfave/cli"
)

const (
	BTCDecimals = 8
	ETHDecimals = 18
)


// Define flags for commands
var (
	networkFlag = cli.StringFlag{
		Name:  "network",
		Value: "mainnet",
		Usage: "name of the test network",
	}

	keyPhraseFlag = cli.StringFlag{
		Name:  "keyphrase",
		Value: "",
		Usage: "keyphrase to unlock the keystore file",
	}

	toFlag = cli.StringFlag{
		Name:  "to",
		Value: "",
		Usage: "receiver address for withdraw",
	}

	valueFlag = cli.Float64Flag{
		Name:  "value",
		Value: 0,
		Usage: "amount of token you want to withdraw,", // todo: specify the unit here
	}
)

func main() {
	app := cli.NewApp()

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "http",
			Usage: "start running the swapper ",
			Flags: []cli.Flag{networkFlag},
			Action: func(c *cli.Context) error {
				// swapper, err  := initializeSwapper(c)
				// if err != nil {
				// 	return err
				// }
				panic("Implement the http logic here")
			},
		},
		{
			Name:  "withdraw",
			Usage: "withdraw the funds in the swapper accounts",
			Action: func(c *cli.Context) error {
				swapper, err  := initializeSwapper(c)
				if err != nil {
					return err
				}

				return withdraw(c,swapper)
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

	cfg, err := config.New(path.Join(os.Getenv("HOME"), fmt.Sprintf(".swapper/%v-config.json", network)), network)
	if err != nil {
		return swapper.Swapper{}, err
	}

	ks, err := keystoreDriver.LoadFromFile(cfg, keyPhrase)
	if err != nil {
		return swapper.Swapper{}, err
	}

	return swapper.NewSwapper(cfg,ks), nil
}

func withdraw(ctx *cli.Context, swapper swapper.Swapper) error {
	// Parse and validate the receiver address
	receiver := ctx.String("to")
	if receiver == ""{
		return errors.New("receiver address cannot be empty")
	}
	if !strings.HasPrefix(receiver, "0x"){
		receiver = "0x" + receiver
	}
	if len(receiver) != 42 {
		return errors.New("invalid receiver address")
	}

	// Parse and validate the value to withdraw
	value := ctx.Float64("value")
	if value == 0 {
		return errors.New("please enter a valid withdraw amount ")
	}
	var valueBig *big.Int

	// Parse and validate the token
	tokenStr := strings.ToLower(strings.TrimSpace(ctx.String("token")))
	var tk token.Token
	switch tokenStr{
	case "btc", "bitcoin":
		tk = token.BTC
		valueBig, _ = big.NewFloat(value * math.Pow10(BTCDecimals)).Int(nil)
	case "eth", "ethereum", "ether":
		tk = token.ETH
		valueBig, _ = big.NewFloat(value * math.Pow10(ETHDecimals)).Int(nil)
	default:
		return errors.New("unknown token")
	}

	swapper.Withdraw(receiver, tk, valueBig)

	return nil
}