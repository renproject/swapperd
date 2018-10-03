package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/republicprotocol/renex-swapper-go/driver/config"
	"github.com/republicprotocol/renex-swapper-go/driver/keystore"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

func main() {
	home := utils.GetHome()
	loc := flag.String("loc", home+"/.swapper", "Location of the swapper's home directory")
	repNet := flag.String("network", "mainnet", "Which republic protocol network to use")
	passphrase := flag.String("keyphrase", "", "Keyphrase to encrypt your key files")
	flag.Parse()
	cmd := exec.Command("mkdir", "-p", *loc)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cfg, err := config.New(*loc, *repNet)
	if err != nil {
		panic(err)
	}
	if err := keystore.GenerateFile(cfg, *passphrase); err != nil {
		if err != keystore.ErrKeyFileExists {
			panic(err)
		}
	}
	addr := readAddress()
	cfg.AuthorizedAddresses = []string{addr}
	config.SaveToFile(fmt.Sprintf("%s/config-%s.json", *loc, *repNet), cfg)
}

func readAddress() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your RenEx Ethereum address: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	addr := strings.TrimSpace(text)
	if len(addr) == 42 && addr[:2] == "0x" {
		addr = addr[2:]
	}
	addrBytes, err := hex.DecodeString(addr)
	if err != nil || len(addrBytes) != 20 {
		fmt.Println("Please enter a valid Ethereum address")
		return readAddress()
	}
	return "0x" + addr
}
