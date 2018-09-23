package main

import (
	"bufio"
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
	repNet := flag.String("republic", "testnet", "Which republic protocol network to use")
	flag.Parse()
	cmd := exec.Command("mkdir", "-p", *loc)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Print("Enter a passphrase (this is used to encrypt your keystore files): ")
	bytePassphrase, err := terminal.ReadPassword(0)
	if err == nil {
		fmt.Println("\nPassphrase typed: " + string(bytePassphrase))
	}
	passphrase := strings.Trim(string(bytePassphrase), "\r\n")
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
	switch system{
	case "window":
		return os.Getenv("userprofile")
	case "linux", "darwin":
		return os.Getenv("HOME")
	default:
		panic("unknown Operating System")
	}
}

func readAddress() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your RenEx Ethereum address: ")
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
