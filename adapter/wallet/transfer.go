package wallet

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/renproject/libbtc-go"
	"github.com/renproject/libeth-go"
	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/tokens"
)

func (wallet *wallet) Transfer(password string, token tokens.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error) {
	switch token.Blockchain {
	case tokens.BITCOIN:
		return wallet.transferBTC(password, to, amount, speed, sendAll)
	case tokens.ETHEREUM:
		return wallet.transferETH(password, to, amount, speed, sendAll)
	case tokens.ERC20:
		return wallet.transferERC20(password, token, to, amount, speed, sendAll)
	default:
		return "", blockchain.Cost{}, tokens.NewErrUnsupportedToken(string(token.Name))
	}
}

func (wallet *wallet) transferBTC(password, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error) {
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
	txHash, txFee, err := account.Transfer(ctx, to, amount.Int64(), libbtc.Fast, sendAll)
	if err != nil {
		return txHash, blockchain.Cost{}, err
	}
	cost[tokens.NameBTC] = big.NewInt(txFee)
	return txHash, cost, nil
}

func (wallet *wallet) transferETH(password, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cost := blockchain.Cost{}
	account, err := wallet.EthereumAccount(password)
	if err != nil {
		return "", cost, err
	}
	tx, err := account.Transfer(ctx, common.HexToAddress(to), amount, libeth.TxExecutionSpeed(speed), 0, sendAll)
	if err != nil {
		return "", cost, err
	}
	cost[tokens.NameETH] = new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	return tx.Hash().String(), cost, nil
}

func (wallet *wallet) transferERC20(password string, token tokens.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error) {
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
	tx, err := erc20.Transfer(ctx, common.HexToAddress(to), amount, libeth.TxExecutionSpeed(speed), sendAll)
	if err != nil {
		return "", cost, err
	}
	cost[tokens.NameETH] = tx.Cost()
	if txFee := token.AdditionalTransactionFee(amount); txFee != nil {
		cost[token.Name] = txFee
	}
	return tx.Hash().String(), cost, nil
}

func (wallet *wallet) Lookup(token tokens.Token, txHash string) (transfer.UpdateReceipt, error) {
	switch token.Blockchain {
	case tokens.BITCOIN:
		return wallet.bitcoinLookup(txHash)
	case tokens.ETHEREUM, tokens.ERC20:
		return wallet.ethereumLookup(txHash)
	default:
		return transfer.UpdateReceipt{}, tokens.NewErrUnsupportedBlockchain(token.Blockchain)
	}
}

func (wallet *wallet) ethereumLookup(txHash string) (transfer.UpdateReceipt, error) {
	client, err := libeth.NewInfuraClient(wallet.config.Ethereum.Network.Name, "172978c53e244bd78388e6d50a4ae2fa")
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
	client, err := libbtc.NewBlockchainInfoClient(wallet.config.Bitcoin.Network.Name)
	if err != nil {
		return transfer.UpdateReceipt{}, err
	}

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
