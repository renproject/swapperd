package main

import (
	"os"
	"os/signal"

	"github.com/renproject/swapperd/driver/composer"
)

func main() {
	done := make(chan struct{})
	go composer.Run(done)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
}
