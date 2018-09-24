package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/republicprotocol/renex-swapper-go/driver/config"
	"github.com/republicprotocol/renex-swapper-go/driver/keystore"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	home := getHome()
	loc := flag.String("loc", home+"/.swapper", "Location of the swapper's home directory")
	repNet := flag.String("network", "testnet", "Which republic protocol network to use")
	flag.Parse()
	cmd := exec.Command("mkdir", "-p", *loc)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Println("The following passphrase is used to encrypt your keystore files")
	passphrase := readPassphrase()
	if err := keystore.GenerateFile(*loc, *repNet, passphrase); err != nil {
		panic(err)
	}
	addr := readAddress()
	cfg := config.New(*loc, *repNet)
	cfg.AuthorizedAddresses = []string{addr}
	config.SaveToFile(fmt.Sprintf("%s/config-%s.json", *loc, *repNet), cfg)
}

func getHome() string {
	system := runtime.GOOS
	switch system {
	case "window":
		return os.Getenv("userprofile")
	case "linux", "darwin":
		return os.Getenv("HOME")
	default:
		panic("unknown operating system")
	}
}

func readAddress() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter your RenEx Ethereum address: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	addr := strings.Trim(text, "\r\n")
	if len(addr) == 40 {
		addr = "0x" + addr
	}
	if len(addr) != 42 {
		fmt.Println("Please enter a valid Ethereum address")
		return readAddress()
	}
	return addr
}

func readPassphrase() string {
	fmt.Print("passphrase: ")
	bytePassphrase, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Print("reenter passphrase: ")
	bytePassphraseReenter, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	if bytes.Compare(bytePassphrase, bytePassphraseReenter) != 0 {
		fmt.Println("passphrase mismatch, please try again")
		return readPassphrase()
	}
	return strings.Trim(string(bytePassphrase), "\r\n")
}
