package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/republicprotocol/renex-swapper-go/driver/config"
	"github.com/republicprotocol/renex-swapper-go/driver/keystore"
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a passphrase (this is used to encrypt your keystore files): ")
	passphraseText, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	passphrase := strings.Trim(passphraseText, "\r\n")
	if err := keystore.GenerateFile(*loc, *repNet, passphrase); err != nil {
		panic(err)
	}
	// addr := readAddress(reader)

	_ = config.New(*loc, *repNet)
	// cfg.AuthorizedAddresses = []string{addr}
	// config.SaveToFile(fmt.Sprintf("%s/config-%s.json", *loc, *repNet), cfg)
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

func readAddress(reader *bufio.Reader) string {
	fmt.Print("Enter your RenEx ethereum address: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	addr := strings.Trim(text, "\r\n")
	if len(addr) == 40 {
		addr = "0x" + addr
	}
	if len(addr) != 42 {
		fmt.Println("Please enter a valid ethereum address")
		return readAddress(reader)
	}
	return addr
}
