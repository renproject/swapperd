package main

import (
	"os"
	"os/signal"
	"path/filepath"

	"github.com/renproject/swapperd/driver/composer"
)

func main() {
	done := make(chan struct{})
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	homeDir := filepath.Dir(filepath.Dir(ex))

	testnet := composer.New(homeDir, "testnet", "17927")
	go testnet.Run(done)
	mainnet := composer.New(homeDir, "mainnet", "7927")
	go mainnet.Run(done)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
