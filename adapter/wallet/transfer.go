package wallet

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/core/transfer"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error) {
	switch token {
	case blockchain.TokenBTC:
		return wallet.transferBTC(password, to, amount)
	case blockchain.TokenETH:
		return wallet.transferETH(password, to, amount)
	case blockchain.TokenWBTC, blockchain.TokenDGX, blockchain.TokenREN,
		blockchain.TokenTUSD, blockchain.TokenZRX, blockchain.TokenOMG,
		blockchain.TokenGUSD, blockchain.TokenDAI, blockchain.TokenUSDC:
		return wallet.transferERC20(password, token, to, amount)
	default:
		return "", blockchain.NewErrUnsupportedToken(token.Name)
	}
}

func (wallet *wallet) transferBTC(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	account, err := wallet.BitcoinAccount(password)
	if err != nil {
		return "", err
	}
	return account.Transfer(ctx, to, amount.Int64())
}

func (wallet *wallet) transferETH(password, to string, amount *big.Int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return "", err
	}
	return account.Transfer(ctx, common.HexToAddress(to), amount, 0)
}

func (wallet *wallet) transferERC20(password string, token blockchain.Token, to string, amount *big.Int) (string, error) {
	var txHash string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return txHash, err
	}
	tokenAddress, err := account.ReadAddress(string(token.Name))
	if err != nil {
		return txHash, err
	}

	tokenContract, err := erc20.NewCompatibleERC20(tokenAddress, bind.ContractBackend(account.EthClient()))
	if err != nil {
		return txHash, err
	}

	if err := account.Transact(
		ctx,
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := tokenContract.Transfer(tops, common.HexToAddress(to), amount)
			if err != nil {
				return tx, err
			}
			txHash = tx.Hash().String()
			return tx, nil
		},
		nil,
		1,
	); err != nil {
		return txHash, err
	}

	return txHash, nil
}

func (wallet *wallet) Lookup(token blockchain.Token, txHash string) (transfer.UpdateReceipt, error) {
	switch token.Blockchain {
	case blockchain.Bitcoin:
		return wallet.bitcoinLookup(txHash)
	case blockchain.Ethereum:
		return wallet.ethereumLookup(txHash)
	default:
		return transfer.UpdateReceipt{}, blockchain.NewErrUnsupportedBlockchain(token.Blockchain)
	}
}

func (wallet *wallet) ethereumLookup(txHash string) (transfer.UpdateReceipt, error) {
	client, err := beth.Connect(wallet.config.Ethereum.Network.URL)
	if err != nil {
		return transfer.UpdateReceipt{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	txBlockNumber, err := client.TxBlockNumber(ctx, txHash)
	if err != nil {
		return transfer.UpdateReceipt{}, err
	}

	currBlockNumber, err := client.CurrentBlockNumber(ctx)
	if err != nil {
		return transfer.UpdateReceipt{}, err
	}

	confirmations := new(big.Int).Sub(currBlockNumber, txBlockNumber)
	return transfer.NewUpdateReceipt(txHash, func(receipt *transfer.TransferReceipt) {
		receipt.Confirmations = confirmations.Int64()
	}), nil
}

func (wallet *wallet) bitcoinLookup(txHash string) (transfer.UpdateReceipt, error) {
	client := libbtc.NewBlockchainInfoClient(wallet.config.Bitcoin.Network.Name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	confirmations, err := client.Confirmations(ctx, txHash)
	if err != nil {
		return transfer.UpdateReceipt{}, err
	}

	return transfer.NewUpdateReceipt(txHash, func(receipt *transfer.TransferReceipt) {
		receipt.Confirmations = confirmations
	}), nil
}
