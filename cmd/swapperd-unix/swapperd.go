package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/republicprotocol/swapperd/driver/composer"
)

func main() {
	flag.Parse()
	done := make(chan struct{})

	testnet := composer.New("testnet", "17927")
	go testnet.Run(done)
	mainnet := composer.New("mainnet", "7927")
	go mainnet.Run(done)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
