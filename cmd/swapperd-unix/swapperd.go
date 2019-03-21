package main

import (
	"os"
	"os/signal"

	"github.com/renproject/swapperd/driver/composer"
)

var version = "undefined"

func main() {
	done := make(chan struct{})
	go composer.Run(version, done)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
