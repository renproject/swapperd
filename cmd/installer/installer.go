package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/renproject/swapperd/driver/keystore"
	"github.com/renproject/swapperd/driver/service"
	"github.com/renproject/swapperd/driver/updater"
	"github.com/tyler-smith/go-bip39"
)

const reset = "\033[m"
const cyan = "\033[36m"
const bold = "\033[1m"

func main() {
	mnemonicFlag := flag.String("mnemonic", "", "Mneumonic for restoring an existing account")
	flag.Parse()
	mnemonic := *mnemonicFlag
	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			panic(err)
		}

		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			panic(err)
		}
	}
	createKeystore("testnet", mnemonic)
	createKeystore("mainnet", mnemonic)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if err := updateSwapperd(dir); err != nil {
		panic(err)
	}
	if err := startServices(); err != nil {
		panic(err)
	}
}

func updateSwapperd(binDir string) error {
	service.Stop("swapperd-updater")
	service.Stop("swapperd")
	service.Delete("swapperd-updater")
	service.Delete("swapperd")
	updater, err := updater.New()
	updater.UseNightlyChannel()
	if err != nil {
		return err
	}
	return updater.Update()
}

func createKeystore(network, mnemonic string) {
	homeDir := getSwapperHome()
	if _, err := keystore.Wallet(homeDir, network, nil); err == nil {
		return
	}
	if err := keystore.Generate(homeDir, network, mnemonic); err != nil {
		panic(err)
	}
}

func createHomeDir(loc string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		return exec.Command("mkdir", "-p", loc).Run()
	case "windows":
		return nil
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func getSwapperHome() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(filepath.Dir(ex))
}

func startServices() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	homeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	if err := service.Create("swapperd", path.Join(homeDir, fmt.Sprintf("swapperd%s", path.Ext(ex)))); err != nil {
		return err
	}
	if err := service.Start("swapperd"); err != nil {
		return err
	}
	if err := service.Create("swapperd-updater", path.Join(homeDir, fmt.Sprintf("swapperd-updater%s", path.Ext(ex)))); err != nil {
		return err
	}
	if err := service.Start("swapperd-updater"); err != nil {
		return err
	}
	return nil
}
