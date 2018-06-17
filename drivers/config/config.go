package main

import (
	"flag"

	"github.com/republicprotocol/atom-go/adapters/config"
)

func main() {
	configParam := flag.String("config", "./config.json", "JSON configuration file")
	flag.Parse()

	conf, err := config.Read(*configParam)
	if err != nil {
		panic(err)
	}
	conf.Ethereum.Mainnet.Chain = "Updated"

	if err := config.Write(conf, *configParam); err != nil {
		panic(err)
	}
}
