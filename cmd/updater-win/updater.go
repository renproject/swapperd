package main

import (
	"fmt"

	"github.com/renproject/swapperd/driver/updater"
	"github.com/renproject/swapperd/driver/winexec"
)

func main() {
	updater, err := updater.New(
		func() {
			winexec.StopService("swapperd")
		},
		func() {
			winexec.StartService("swapperd")
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := updater.Update(); err != nil {
		fmt.Println(err)
	}
}
