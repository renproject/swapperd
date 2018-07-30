package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/republicprotocol/atom-go/adapters/configs/network"
)

func main() {
	net, err := network.LoadNetwork("../network.json")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Bitcoin Node IP Address: (<ipaddress>:<port>): ")
	ipAddr, _ := reader.ReadString('\n')
	fmt.Print("Enter Bitcoin RPC UserName: ")
	rpcUser, _ := reader.ReadString('\n')
	fmt.Print("Enter Bitcoin RPC Password: ")
	rpcPass, _ := reader.ReadString('\n')

	net.Bitcoin.Password = strings.Trim(rpcPass, "\n")
	net.Bitcoin.User = strings.Trim(rpcUser, "\n")
	net.Bitcoin.URL = strings.Trim(ipAddr, "\n")

	if err := net.Update(); err != nil {
		panic(err)
	}
}
