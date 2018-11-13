package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/republicprotocol/swapperd/driver/swapperd"
)

func main() {
	network := flag.String("network", "testnet", "Which network to use")
	port := flag.String("port", "7777", "Which network to use")
	flag.Parse()

	done := make(chan struct{})
	swapperd.Run(done, *network, *port)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
