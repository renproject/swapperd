package main

import (
	"fmt"

	"github.com/renproject/swapperd/driver/updater"
)

func main() {
	updater, err := updater.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := updater.Update(); err != nil {
		fmt.Println(err)
	}
}
