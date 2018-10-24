package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/republicprotocol/swapperd/driver/keystore"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/ssh/terminal"
)

const reset = "\033[m"
const cyan = "\033[36m"
const bold = "\033[1m"

func main() {
	network := flag.String("network", "testnet", "Which republic protocol network to use")
	usernameFlag := flag.String("username", "", "Username for HTTP basic authentication")
	passwordHashFlag := flag.String("passwordHash", "", "Password Hash for HTTP basic authentication")
	flag.Parse()

	if _, err := keystore.LoadAccounts(*network); err == nil {
		fmt.Printf("Swapper already exists at the default location (%s)\n", getDefaultSwapperHome())
		return
	}

	if err := createHomeDir(); err != nil {
		panic(err)
	}

	var username, passwordHash string

	if *usernameFlag != "" && *passwordHashFlag != "" {
		username = *usernameFlag
		passwordHash = *passwordHashFlag
	} else {
		username, passwordHash = getCredentials()
	}

	mnemonic, err := keystore.Generate(*network, username, passwordHash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n%sPlease backup the following mnemonic to restore your swapper wallet:\n", bold)
	fmt.Printf("%s%s%s\n", cyan, mnemonic, reset)
}

func getCredentials() (string, string) {
	var password []byte
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")
	user, _ := reader.ReadString('\n')
	user = strings.Trim(user, "\r\n")

	for {
		fmt.Print("Enter Password: ")
		passwordEnter, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}

		fmt.Print("\nReenter Password: ")
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
	passwordHash := sha3.Sum256(password)
	return user, base64.StdEncoding.EncodeToString(passwordHash[:])
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
		cmd := exec.Command("cmd", "/C", "md", loc)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
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
		return strings.Join(strings.Split(windows, "\\"), "\\\\") + "\\swapperd"
	}
	panic("unknown Operating System")
}
