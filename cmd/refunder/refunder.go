package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	btcbindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/btc"
	btcclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/btc"
	"github.com/republicprotocol/renex-swapper-go/utils"

	ethbindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/eth"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	configDriver "github.com/republicprotocol/renex-swapper-go/driver/config"
	keystoreDriver "github.com/republicprotocol/renex-swapper-go/driver/keystore"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

}

func refundEther(reader *bufio.Reader, cfg config.EthereumNetwork, key keystore.EthereumKey) {
	fmt.Println("Swap ID : ")
	swapID, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	swapID = strings.Trim(swapID, "\n\r")
	swapIDBytes, err := hex.DecodeString(swapID)
	if err != nil {
		panic(err)
	}
	swapIDBytes32, err := utils.ToBytes32(swapIDBytes)
	if err != nil {
		panic(err)
	}
	conn, err := ethclient.NewConnWithConfig(cfg)
	if err != nil {
		panic(err)
	}
	swapper, err := ethbindings.NewRenExAtomicSwapper(conn.RenExAtomicSwapper, bind.ContractBackend(conn.Client))
	if err != nil {
		panic(err)
	}
	tx, err := swapper.Refund(key.TransactOpts, swapIDBytes32)
	if err != nil {
		panic(err)
	}
	if _, err := conn.PatchedWaitMined(context.Background(), tx); err != nil {
		panic(err)
	}
}

func refundBitcoin(reader *bufio.Reader, cfg config.BitcoinNetwork, key keystore.BitcoinKey) {
	fmt.Println("Script : ")
	script, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	script = strings.Trim(script, "\n\r")
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Script Transaction: ")
	scriptTx, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	scriptTx = strings.Trim(scriptTx, "\n\r")
	scriptTxBytes, err := hex.DecodeString(scriptTx)
	if err != nil {
		panic(err)
	}
	conn, err := btcclient.NewConnWithConfig(cfg)
	if err != nil {
		panic(err)
	}
	if err := btcbindings.Refund(conn, key, scriptBytes, scriptTxBytes); err != nil {
		panic(err)
	}
}

func configgen(reader *bufio.Reader) (config.Config, keystore.Keystore) {
	fmt.Println("Passphrase : ")
	passphrase, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	passphrase = strings.Trim(passphrase, "\n\r")
	conf := configDriver.New(utils.GetHome()+"/.swapper", "testnet")
	ks := keystoreDriver.LoadFromFile("testnet", utils.GetHome()+"/.swapper", passphrase)
	return conf, ks
}

func blockchainSelector(reader *bufio.Reader) {
	fmt.Println("Blockchain : ")
	blockchain, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	blockchain = strings.Trim(blockchain, "\n\r")
	if strings.ToLower(blockchain) 
	
	
	
	
	{

	}
}
