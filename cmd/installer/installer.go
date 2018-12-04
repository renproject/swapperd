package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

const reset = "\033[m"
const cyan = "\033[36m"
const bold = "\033[1m"

func main() {
	usernameFlag := flag.String("username", "", "Username for HTTP basic authentication")
	passwordFlag := flag.String("password", "", "Password for HTTP basic authentication")
	mnemonicFlag := flag.String("mnemonic", "", "Mneumonic for restoring an existing account")
	flag.Parse()

	var username, password string
	if *usernameFlag != "" && *passwordFlag != "" {
		username = *usernameFlag
		password = *passwordFlag
	} else {
		username, password = credentials()
	}

	if *mnemonicFlag != "" {
		createKeystore("testnet", username, password, *mnemonicFlag)
		createKeystore("mainnet", username, password, *mnemonicFlag)
		return
	}

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		panic(err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}
	createKeystore("testnet", username, password, mnemonic)
	createKeystore("mainnet", username, password, mnemonic)
}

func createKeystore(network, username, password, mnemonic string) {
	if _, err := keystore.FundManager(network); err == nil {
		fmt.Printf("swapper already exists at the default location (%s)\n", getDefaultSwapperHome())
		return
	}

	if err := createHomeDir(); err != nil {
		panic(err)
	}

	if err := keystore.Generate(network, username, password, mnemonic); err != nil {
		panic(err)
	}
}

func credentials() (string, string) {
	var password []byte
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose a Username: ")
	user, _ := reader.ReadString('\n')
	user = strings.Trim(user, "\r\n")

	for {
		fmt.Print("Choose a Password: ")
		passwordEnter, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}

		fmt.Print("\nReenter the Password: ")
		passwordReenter, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}

		if bytes.Compare(passwordEnter, passwordReenter) == 0 {
			password = passwordEnter
			break
		}

		fmt.Println("password mismatch, please try again")
	}
	return user, string(password)
}

func createHomeDir() error {
	loc := getDefaultSwapperHome()
	unix := os.Getenv("HOME")
	if unix != "" {
		cmd := exec.Command("mkdir", "-p", loc)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return nil
	}
	return errors.New("unknown Operating System")
}

func getDefaultSwapperHome() string {
	unix := os.Getenv("HOME")
	if unix != "" {
		return unix + "/.swapperd"
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return "C:\\windows\\system32\\config\\systemprofile\\swapperd"
	}
	panic("unknown Operating System")
}
