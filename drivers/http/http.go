package main

import (
	"flag"
	"log"
	netHttp "net/http"

	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/http"
	"github.com/republicprotocol/atom-go/adapters/keystore"
)

func main() {

	confPath := flag.String("config", "./config.json", "Location of the config file")
	keystrPath := flag.String("keystore", "./keystore.json", "Location of the keystore file")

	conf, err := config.LoadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	keystr := keystore.NewKeystore(*keystrPath)

	key, err := keystr.LoadKeypair("ethereum")
	if err != nil {
		panic(err)
	}

	httpAdapter, err := http.NewBoxHttpAdapter(conf, key)
	if err != nil {
		panic(err)
	}

	log.Println("Listening on 0.0.0.0:18516")

	log.Fatal(netHttp.ListenAndServe(":18516", http.NewServer(httpAdapter)))
}
