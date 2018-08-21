package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	config "github.com/republicprotocol/renex-swapper-go/adapters/configs/general"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/network"
)

func main() {
	home := getHome()
	ethNet := flag.String("ethereum", "kovan", "Which ethereum network to use")
	btcNet := flag.String("bitcoin", "testnet", "Which bitcoin network to use")

	keystore.NewKeystore([]uint32{0, 1}, []string{*btcNet, *ethNet}, home+"/.swapper/keystore.json")

	cfg, err := config.LoadConfig(home + "/.swapper/config.json")
	if err != nil {
		panic(err)
	}

	addresses := []string{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your ethereum address(es): (\033[32mClick Enter to Finish\033[m) \nAddress>")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if text == "\n" {
			break
		}
		addresses = append(addresses, strings.Trim(text, "\r\n"))
		fmt.Print("Address>")
	}
	cfg.AuthorizedAddresses = addresses
	cfg.StoreLoc = home + "/.swapper/db"
	cfg.RenGuardAddr = "renex-watchdog-testnet.herokuapp.com"

	if err := cfg.Update(); err != nil {
		panic(err)
	}

	net, err := network.LoadNetwork(home + "/.swapper/network.json")
	if err != nil {
		panic(err)
	}

	fmt.Print("Enter Bitcoin Node IP Address: (<ipaddress>:<port>): ")
	ipAddr, _ := reader.ReadString('\n')
	fmt.Print("Enter Bitcoin RPC UserName: ")
	rpcUser, _ := reader.ReadString('\n')
	fmt.Print("Enter Bitcoin RPC Password: ")
	rpcPass, _ := reader.ReadString('\n')

	net.Bitcoin.Password = strings.Trim(rpcPass, "\r\n")
	net.Bitcoin.User = strings.Trim(rpcUser, "\r\n")
	net.Bitcoin.URL = strings.Trim(ipAddr, "\r\n")

	if err := net.Update(); err != nil {
		panic(err)
	}
}

func getHome() string {
	winHome := os.Getenv("userprofile")
	unixHome := os.Getenv("HOME")

	if winHome != "" {
		return winHome
	}

	if unixHome != "" {
		return unixHome
	}

	panic("unknown Operating System")
}
