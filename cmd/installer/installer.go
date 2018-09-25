package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/republicprotocol/renex-swapper-go/driver/config"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/republicprotocol/renex-swapper-go/driver/keystore"
)

func main() {
	home := getHome()
	loc := flag.String("loc", home+"/.swapper", "Location of the swapper's home directory")
	repNet := flag.String("network", "testnet", "Which republic protocol network to use")
	passphrase := flag.String("passphrase", "", "Passphrase to encrypt your key files")

	flag.Parse()

	cmd := exec.Command("mkdir", "-p", *loc)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	cfg ,err  := config.New(*loc, *repNet)
	if err != nil {
		panic(err)
	}
	if err := keystore.GenerateFile(cfg, *passphrase); err != nil {
		panic(err)
	}
	addr := readAddress()
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