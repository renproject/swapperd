package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/republicprotocol/renex-swapper-go/driver/config"
	"github.com/republicprotocol/renex-swapper-go/driver/keystore"
	"github.com/republicprotocol/renex-swapper-go/utils"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	home := utils.GetHome()
	loc := flag.String("loc", home+"/.swapper", "Location of the swapper's home directory")
	repNet := flag.String("network", "mainnet", "Which republic protocol network to use")
	flag.Parse()
	cmd := exec.Command("mkdir", "-p", *loc)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	cfg, err := config.New(*loc, *repNet)
	if err != nil {
		panic(err)
	}
	fmt.Println("The following passphrase is used to encrypt your keystore files")
	passphrase := readPassphrase()
	if err := keystore.GenerateFile(cfg, passphrase); err != nil {
		panic(err)
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
	addr := strings.Trim(text, "\r\n")
	if len(addr) == 42 && addr[:2] == "0x" {
		addr = addr[2:]
	}
	addrBytes, err := hex.DecodeString(addr)
	if err != nil || len(addrBytes) == 40 {
		fmt.Println("Please enter a valid Ethereum address")
		return readAddress()
	}
	return addr
}

func readPassphrase() string {
	fmt.Print("Passphrase: ")
	bytePassphrase, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Print("Reenter passphrase: ")
	bytePassphraseReenter, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	if bytes.Compare(bytePassphrase, bytePassphraseReenter) != 0 {
		fmt.Println("\nPassphrase mismatch, please try again")
		return readPassphrase()
	}
	return strings.Trim(string(bytePassphrase), "\r\n")
}
