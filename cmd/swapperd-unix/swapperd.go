package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/republicprotocol/swapperd/driver/composer"
)

func main() {
	network := flag.String("network", "testnet", "Which network to use")
	port := flag.String("port", "7927", "Which network to use")
	flag.Parse()

	done := make(chan struct{})
	composer := composer.New(*network, *port)
	go composer.Run(done)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
