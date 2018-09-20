package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	btcclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/btc"

	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Script : ")
	script, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Script Transaction: ")
	scriptTx, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	scriptTxBytes, err := hex.DecodeString(scriptTx)
	if err != nil {
		panic(err)
	}
}

func refundEther(cfg config.EthereumNetwork, key keystore.EthereumKey, swapID [32]byte) {
	conn, err := ethclient.NewConnWithConfig()
}

func refundBitcoin(cfg config.BitcoinNetwork, key keystore.BitcoinKey, script, scriptTx []byte) {
	conn, err := btcclient.NewConnWithConfig(cfg)
	if err != nil {
		panic(err)
	}
	if err := bindings.Refund(conn, key, script, scriptTx); err != nil {
		panic(err)
	}
}

func configgen() (config.Config, keystore.Keystore) {
	return nil, nil
}
