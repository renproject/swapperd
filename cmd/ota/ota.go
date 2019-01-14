package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli"
)

var (
	TokenPairFlag = cli.StringFlag{
		Name:  "tokenPair",
		Usage: "The token pair you are trading",
	}

	SecretFlag = cli.StringFlag{
		Name:  "secret",
		Usage: "The secret for identifying a swap",
	}

	SwapObjFlag = cli.StringFlag{
		Name:  "swapObj",
		Usage: "The swap object to do the atomic swap",
	}

	SwapIDFlag = cli.StringFlag{
		Name:  "swapID",
		Usage: "The atomic swap id",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "Swapperd OTA"
	app.Usage = "A command-line demo for one time addresses."
	app.Version = "1.0.0"

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "new",
			Usage: "Start a new swap",
			Flags: []cli.Flag{
				TokenPairFlag, SecretFlag,
			},
			Action: func(c *cli.Context) error {
				return newSwap(c)
			},
		},
		{
			Name:  "bootload",
			Usage: "Bootload your swaps",
			Flags: []cli.Flag{
				SecretFlag,
			},
			Action: func(c *cli.Context) error {
				return bootload(c)
			},
		},
		{
			Name:  "swap",
			Usage: "Submit the swap object",
			Flags: []cli.Flag{
				SecretFlag, SwapObjFlag,
			},
			Action: func(c *cli.Context) error {
				return createSwap(c)
			},
		},
		{
			Name:  "status",
			Usage: "Status of an atomic swap",
			Flags: []cli.Flag{
				SecretFlag, SwapIDFlag,
			},
			Action: func(c *cli.Context) error {
				return status(c)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		// Remove the timestamp for error message
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
		log.Fatal(err)
	}
}

func newSwap(c *cli.Context) error {
	secret := c.String("secret")
	if secret == "" {
		secretBytes32 := [32]byte{}
		rand.Read(secretBytes32[:])
		secret = base64.StdEncoding.EncodeToString(secretBytes32[:])
	}

	tokenPair := c.String("tokenPair")
	if tokenPair == "" {
		return fmt.Errorf("token pair cannot be empty")
	}

	var err error
	var address string

	switch tokenPair {
	case "BTC/ETH", "BTC/WBTC":
		address, err = getAddress("BTC", secret)
		if err != nil {
			return err
		}
	case "ETH/BTC", "WBTC/BTC", "ETH/WBTC", "WBTC/ETH":
		address, err = getAddress("ETH", secret)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported token pair")
	}

	fmt.Printf("\nDeposit Address: %s\n\n", address)
	if c.String("secret") == "" {
		fmt.Printf("\nSecret: %s\n", secret)
	}
	return nil
}

func createSwap(c *cli.Context) error {
	secret := c.String("secret")
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	swap := c.String("swapObj")
	if swap == "" {
		return fmt.Errorf("swap object cannot be empty")
	}

	resp, err := postSwap(swap, secret)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s", string(resp))
	return nil
}

func status(c *cli.Context) error {
	secret := c.String("secret")
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	swapID := c.String("swapID")
	if swapID == "" {
		return fmt.Errorf("swap object cannot be empty")
	}

	resp, err := getStatus(swapID, secret)
	if err != nil {
		return err
	}
	fmt.Printf("Swap Receipt: \n%s", string(resp))
	return nil
}

func bootload(c *cli.Context) error {
	secret := c.String("secret")
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}
	return postBootload(secret)
}

func getAddress(token, secret string) (string, error) {
	req, err := http.NewRequest("GET", "http://localhost:17927/addresses", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(string(respBytes))
	}

	addresses := map[string]string{}
	if err := json.Unmarshal(respBytes, &addresses); err != nil {
		return "", err
	}

	return addresses[token], nil
}

func postSwap(swap, secret string) ([]byte, error) {
	req, err := http.NewRequest("POST", "http://localhost:17927/swaps", bytes.NewBuffer([]byte(swap)))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func postBootload(secret string) error {
	req, err := http.NewRequest("POST", "http://localhost:17927/bootload", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(data))
	}
	return nil
}

func getStatus(swapID, secret string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:17927/swaps?id=%s", url.QueryEscape(swapID)), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
