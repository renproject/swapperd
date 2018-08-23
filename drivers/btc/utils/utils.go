package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"

	rpc "github.com/btcsuite/btcd/rpcclient"
)

func main() {
	// cmd := start()
	// defer stop(cmd)

	// chainParams := &chaincfg.RegressionNetParams
	connect, err := normalizeAddress("localhost", "18443")

	if err != nil {
		panic(err)
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         "testuser",
		Pass:         "testpassword",
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	rpcClient, err := rpc.New(connConfig, nil)
	if err != nil {
		panic(fmt.Errorf("rpc connect: %v", err))
	}

	// err = rpcClient.CreateNewAccount("Atom")
	// if err != nil {
	// 	panic(fmt.Errorf("Failed to create an account: %v", err))
	// }

	addrs, err := rpcClient.GetAddressesByAccount("alice")
	if err != nil {
		panic(fmt.Errorf("Failed to get addresses of an account: %v", err))
	}

	for _, addr := range addrs {
		fmt.Println(rpcClient.GetReceivedByAddress(addr))
	}
}

func start() *exec.Cmd {
	cmd := exec.Command("bitcoind", "--regtest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd
}

func stop(cmd *exec.Cmd) {
	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Wait()
}

func normalizeAddress(addr string, defaultPort string) (hostport string, err error) {
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err = net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}
