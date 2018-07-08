package main

import (
	"flag"
	"fmt"
	"log"
	netHttp "net/http"

	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/http"
	"github.com/republicprotocol/atom-go/adapters/keystore"
)

func main() {

	port := flag.String("port", "18516", "HTTP Atom port")
	confPath := flag.String("config", "./config.json", "Location of the config file")
	keystrPath := flag.String("keystore", "./keystore.json", "Location of the keystore file")

	flag.Parse()

	conf, err := config.LoadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	keystr := keystore.NewKeystore(*keystrPath)

	httpAdapter, err := http.NewBoxHttpAdapter(conf, keystr)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("0.0.0.0:%s", *port))
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%s", *port), http.NewServer(httpAdapter)))
}
