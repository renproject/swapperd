package wallet

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
)

func (wallet *wallet) Transfer(password string, token blockchain.Token, to string, amount, fee *big.Int, sendAll bool) (string, blockchain.Cost, error) {
	switch token.Blockchain {
	case blockchain.Bitcoin:
		return wallet.transferBTC(password, to, amount, fee, sendAll)
	case blockchain.Ethereum:
		return wallet.transferETH(password, to, amount, fee, sendAll)
	case blockchain.ERC20:
		return wallet.transferERC20(password, token, to, amount, fee, sendAll)
	default:
		return "", blockchain.Cost{}, blockchain.NewErrUnsupportedToken(token.Name)
	}
}

func (wallet *wallet) transferBTC(password, to string, amount, fee *big.Int, sendAll bool) (string, blockchain.Cost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cost := blockchain.Cost{}
	account, err := wallet.BitcoinAccount(password)
	if err != nil {
		return "", blockchain.Cost{}, err
	}

	if amount == nil {
		amount = big.NewInt(0)
	}
	if fee == nil {
		fee = big.NewInt(10000)
	}
	txHash, err := account.Transfer(ctx, to, amount.Int64(), fee.Int64(), sendAll)
	if err != nil {
		return txHash, blockchain.Cost{}, err
	}
	cost[blockchain.BTC] = fee
	return txHash, cost, nil
}

func (wallet *wallet) transferETH(password, to string, amount, fee *big.Int, sendAll bool) (string, blockchain.Cost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cost := blockchain.Cost{}
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return "", cost, err
	}
	tx, err := account.Transfer(ctx, common.HexToAddress(to), amount, fee, 0, sendAll)
	if err != nil {
		return "", cost, err
	}
	cost[blockchain.ETH] = new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	return tx.Hash().String(), cost, nil
}

func (wallet *wallet) transferERC20(password string, token blockchain.Token, to string, amount, gasPrice *big.Int, sendAll bool) (string, blockchain.Cost, error) {
	cost := blockchain.Cost{}
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return "", cost, err
	}
	erc20, err := account.NewERC20(string(token.Name))
	if err != nil {
		return "", cost, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if sendAll {
		amount, err = erc20.BalanceOf(ctx, account.Address())
	}
	tx, err := erc20.Transfer(ctx, common.HexToAddress(to), amount, gasPrice, sendAll)
	if err != nil {
		return "", cost, err
	}
	cost[blockchain.ETH] = tx.Cost()
	if txFee := token.AdditionalTransactionFee(amount); txFee != nil {
		cost[token.Name] = txFee
	}
	return tx.Hash().String(), cost, nil
}

func (wallet *wallet) Lookup(token blockchain.Token, txHash string) (transfer.UpdateReceipt, error) {
	switch token.Blockchain {
	case blockchain.Bitcoin:
		return wallet.bitcoinLookup(txHash)
	case blockchain.Ethereum, blockchain.ERC20:
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
