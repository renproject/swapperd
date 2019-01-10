package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/republicprotocol/swapperd/driver/keystore"
	"github.com/tyler-smith/go-bip39"
)

const reset = "\033[m"
const cyan = "\033[36m"
const bold = "\033[1m"

func main() {
	mnemonicFlag := flag.String("mnemonic", "", "Mneumonic for restoring an existing account")
	flag.Parse()

	if *mnemonicFlag != "" {
		createKeystore("testnet", *mnemonicFlag)
		createKeystore("mainnet", *mnemonicFlag)
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
	createKeystore("testnet", mnemonic)
	createKeystore("mainnet", mnemonic)
}

func createKeystore(network, mnemonic string) {
	homeDir := getDefaultSwapperHome()
	if _, err := keystore.Wallet(homeDir, network); err == nil {
		fmt.Printf("swapper already exists at the default location (%s)\n", getDefaultSwapperHome())
		return
	}

	if err := createHomeDir(); err != nil {
		panic(err)
	}

	if err := keystore.Generate(homeDir, network, mnemonic); err != nil {
		panic(err)
	}
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
	windows := os.Getenv("programfiles(x86)")
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

	windows := os.Getenv("programfiles(x86)")
	if windows != "" {
		return windows + "\\Swapperd"
	}
	panic("unknown Operating System")
}
