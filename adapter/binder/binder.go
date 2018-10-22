package binder

import (
	"fmt"

	"github.com/republicprotocol/swapperd/adapter/account"
	"github.com/republicprotocol/swapperd/adapter/binder/btc"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/adapter/binder/eth"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type builder struct {
	account.Accounts
	swapper.Logger
}

func NewBuilder(accounts account.Accounts, logger swapper.Logger) swapper.BinderBuilder {
	return &builder{
		accounts,
		logger,
	}
}

func (builder *builder) BuildBinders(query swapper.Query) (swapper.SwapContractBinder, swapper.SwapContractBinder, error) {
	native, foreign, err := builder.buildComplementarySwaps(query.Swap)
	if err != nil {
		return nil, nil, err
	}
	nativeSCB, err := builder.buildBinder(native, query.Password)
	if err != nil {
		return nil, nil, err
	}
	foreignSCB, err := builder.buildBinder(foreign, query.Password)
	if err != nil {
		return nil, nil, err
	}
	return nativeSCB, foreignSCB, nil
}

func (builder *builder) buildBinder(swap foundation.SwapTry, password string) (swapper.SwapContractBinder, error) {
	switch swap.Token {
	case foundation.TokenBTC:
		return btc.NewBTCSwapContractBinder(builder.GetBitcoinAccount(), swap, builder.Logger)
	case foundation.TokenETH:
		return eth.NewETHSwapContractBinder(builder.GetEthereumAccount(), swap, builder.Logger)
	case foundation.TokenWBTC:
		return erc20.NewERC20SwapContractBinder(builder.GetEthereumAccount(), swap, builder.Logger)
	default:
		return nil, foundation.NewErrUnsupportedToken(swap.Token.Name)
	}
}

func (builder *builder) buildComplementarySwaps(swap foundation.Swap) (foundation.SwapTry, foundation.SwapTry, error) {
	spendingAddr, fundingAddr, err := builder.calculateAddresses(swap)
	if err != nil {
		return foundation.SwapTry{}, foundation.SwapTry{}, err
	}
	nativeExpiry, foreignExpiry := builder.calculateTimeLocks(swap)
	return builder.buildNativeSwap(swap, nativeExpiry, spendingAddr), builder.buildForeignSwap(swap, foreignExpiry, fundingAddr), nil
}

func (builder *builder) buildNativeSwap(swap foundation.Swap, timelock int64, fundingAddress string) foundation.SwapTry {
	return foundation.SwapTry{
		ID:              swap.ID,
		Token:           swap.SendToken,
		Value:           swap.SendValue,
		SecretHash:      swap.SecretHash,
		TimeLock:        swap.TimeLock,
		SpendingAddress: swap.SendToAddress,
		FundingAddress:  fundingAddress,
	}
}

func (builder *builder) buildForeignSwap(swap foundation.Swap, timelock int64, spendingAddress string) foundation.SwapTry {
	return foundation.SwapTry{
		ID:              swap.ID,
		Token:           swap.ReceiveToken,
		Value:           swap.ReceiveValue,
		SecretHash:      swap.SecretHash,
		TimeLock:        swap.TimeLock,
		SpendingAddress: spendingAddress,
		FundingAddress:  swap.ReceiveFromAddress,
	}
}

func (builder *builder) calculateTimeLocks(swap foundation.Swap) (native, foreign int64) {
	if swap.IsFirst {
		native = swap.TimeLock
		foreign = swap.TimeLock - 24*60*60
		return
	}
	native = swap.TimeLock - 24*60*60
	foreign = swap.TimeLock
	return
}

func (builder *builder) calculateAddresses(swap foundation.Swap) (string, string, error) {
	ethAccount := builder.GetEthereumAccount()
	btcAccount := builder.GetBitcoinAccount()

	ethAddress := ethAccount.Address()
	btcAddress, err := btcAccount.Address()
	if err != nil {
		return "", "", err
	}

	if swap.SendToken.Blockchain == "ethereum" && swap.ReceiveToken.Blockchain == "bitcoin" {
		return ethAddress.String(), btcAddress.EncodeAddress(), nil
	}

	if swap.SendToken.Blockchain == "bitcoin" && swap.ReceiveToken.Blockchain == "ethereum" {
		return btcAddress.EncodeAddress(), ethAddress.String(), nil
	}

	if swap.SendToken.Blockchain == "ethereum" && swap.ReceiveToken.Blockchain == "ethereum" {
		return ethAddress.String(), ethAddress.String(), nil
	}

	return "", "", fmt.Errorf("unsupported blockchain pair : %s <--> %s", swap.SendToken.Blockchain, swap.ReceiveToken.Blockchain)
}
