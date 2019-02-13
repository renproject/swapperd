package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/renproject/swapperd/driver/keystore"
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
	if err := startSwapperd(); err != nil {
		panic(err)
	}
}

func createKeystore(network, mnemonic string) {
	homeDir := getSwapperHome()
	if _, err := keystore.Wallet(homeDir, network); err == nil {
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

func startSwapperd() error {
	home := getSwapperHome()
	switch runtime.GOOS {
	case "linux":
		return startLinuxService(home)
	case "darwin":
		return startDarwinService(home)
	case "windows":
		return startWindowsService(home)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func startLinuxService(swapperdHome string) error {
	serviceLocation := path.Join(os.Getenv("HOME"), ".config", "systemd", "user")
	if err := exec.Command("mkdir", "-p", serviceLocation).Run(); err != nil {
		return err
	}
	serviceContent := fmt.Sprintf("[Unit]\nDescription=Swapper Daemon\nAssertPathExists=%s\n\n[Service]\nWorkingDirectory=%s\nExecStart=%s/bin/swapperd\nRestart=on-failure\nPrivateTmp=true\nNoNewPrivileges=true\n\n# Specifies which signal to use when killing a service. Defaults to SIGTERM.\n# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)\nKillSignal=SIGHUP\n\n[Install]\nWantedBy=default.target", swapperdHome, swapperdHome, swapperdHome)
	servicePath := path.Join(serviceLocation, "swapperd.service")
	if err := ioutil.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return err
	}
	if err := exec.Command("loginctl", "enable-linger", os.Getenv("whoami")).Run(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "--user", "enable", "swapperd.service").Run(); err != nil {
		return err
	}
	return exec.Command("systemctl", "--user", "start", "swapperd.service").Run()
}

func startDarwinService(swapperdHome string) error {
	serviceContent := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\t\n<dict>\t\t\n<key>Label</key>\t\t\n<string>ren.swapperd</string>\t\t\n<key>ProgramArguments</key>\t\t\n<array>\t\t\t\t\n<string>%s/bin/swapperd</string>\t\t\n</array>\t\t\n<key>KeepAlive</key>\t\t\n<true/>\t\t\n<key>StandardOutPath</key>\t\t\n<string>%s/swapperd.log</string>\t\t\n<key>StandardErrorPath</key>\t\t\n<string>%s/swapperd.log</string>\t\n</dict>\n</plist>", swapperdHome, swapperdHome, swapperdHome)
	servicePath := path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "ren.swapperd.plist")
	if err := ioutil.WriteFile(servicePath, []byte(serviceContent), 0755); err != nil {
		return err
	}
	return exec.Command("launchctl", "load", "-w", servicePath).Run()
}

func startWindowsService(swapperdHome string) error {
	if err := exec.Command("cmd", "/C", "sc", "create", "swapperd", "start=", "auto", "binpath=", path.Join(swapperdHome, "bin", "swapperd.exe")).Run(); err != nil {
		return err
	}
	return exec.Command("cmd", "/C", "sc", "start", "swapperd").Run()
}
