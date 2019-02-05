package main

import (
	"os"
	"os/signal"

	"github.com/renproject/swapperd/driver/composer"
)

func main() {
	done := make(chan struct{})
	homeDir := os.Getenv("HOME") + "/.swapperd"

	testnet := composer.New(homeDir, "testnet", "17927")
	go testnet.Run(done)
	mainnet := composer.New(homeDir, "mainnet", "7927")
	go mainnet.Run(done)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
