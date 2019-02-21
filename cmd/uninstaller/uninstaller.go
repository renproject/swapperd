package main

import (
	"github.com/renproject/swapperd/driver/service"
)

func main() {
	service.Stop("swapperd-updater")
	service.Stop("swapperd")
	service.Delete("swapperd-updater")
	service.Delete("swapperd")
}
