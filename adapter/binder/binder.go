package binder

import (
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/renproject/swapperd/adapter/binder/btc"
	"github.com/renproject/swapperd/adapter/binder/erc20"
	"github.com/renproject/swapperd/adapter/binder/eth"
	"github.com/renproject/swapperd/adapter/wallet"
	"github.com/renproject/swapperd/core/wallet/swapper/immediate"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/sirupsen/logrus"
)

type builder struct {
	wallet.Wallet
	logrus.FieldLogger
}

func NewBuilder(wallet wallet.Wallet, logger logrus.FieldLogger) immediate.ContractBuilder {
	return &builder{
		wallet,
		logger,
	}
}

func (builder *builder) BuildSwapContracts(req immediate.SwapRequest) (immediate.Contract, immediate.Contract, error) {
	native, foreign, err := builder.buildComplementarySwaps(req.Blob)
	if err != nil {
		return nil, nil, err
	}
	nativeBinder, err := builder.buildBinder(native, req.SendCost, req.Blob.Password)
	if err != nil {
		return nil, nil, err
	}
	foreignBinder, err := builder.buildBinder(foreign, req.ReceiveCost, req.Blob.Password)
	if err != nil {
		return nil, nil, err
	}
	return nativeBinder, foreignBinder, nil
}

func (builder *builder) buildBinder(swap swap.Swap, cost blockchain.Cost, password string) (immediate.Contract, error) {
	switch swap.Token.Blockchain {
	case tokens.BITCOIN:
		btcAccount, err := builder.BitcoinAccount(password)
		if err != nil {
			return nil, err
		}
		return btc.NewBTCSwapContractBinder(btcAccount, swap, cost, builder.FieldLogger)
	case tokens.ETHEREUM:
		ethAccount, err := builder.EthereumAccount(password)
		if err != nil {
			return nil, err
		}
		return eth.NewETHSwapContractBinder(ethAccount, swap, cost, builder.FieldLogger)
	case tokens.ERC20:
		ethAccount, err := builder.EthereumAccount(password)
		if err != nil {
			return nil, err
		}
		return erc20.NewERC20SwapContractBinder(ethAccount, swap, cost, builder.FieldLogger)
	default:
		return nil, tokens.NewErrUnsupportedToken(string(swap.Token.Name))
	}
}

func (builder *builder) buildComplementarySwaps(blob swap.SwapBlob) (swap.Swap, swap.Swap, error) {
	fundingAddr, spendingAddr, err := builder.calculateAddresses(blob)
	if err != nil {
		return swap.Swap{}, swap.Swap{}, err
	}
	nativeExpiry, foreignExpiry := builder.calculateTimeLocks(blob)

	nativeSwap, err := builder.buildNativeSwap(blob, nativeExpiry, fundingAddr)
	if err != nil {
		return swap.Swap{}, swap.Swap{}, err
	}
	foreignSwap, err := builder.buildForeignSwap(blob, foreignExpiry, spendingAddr)
	if err != nil {
		return swap.Swap{}, swap.Swap{}, err
	}
	return nativeSwap, foreignSwap, nil
}

func (builder *builder) buildNativeSwap(blob swap.SwapBlob, timelock int64, fundingAddress string) (swap.Swap, error) {
	token, err := tokens.PatchToken(blob.SendToken)
	if err != nil {
		return swap.Swap{}, err
	}
	value, ok := new(big.Int).SetString(blob.SendAmount, 10)
	if !ok {
		return swap.Swap{}, fmt.Errorf("corrupted send value: %v", blob.SendAmount)
	}

	brokerFee := big.NewInt(0)
	if blob.BrokerFee != 0 {
		if err := builder.Wallet.VerifyAddress(token.Blockchain, blob.BrokerSendTokenAddr); err != nil {
			return swap.Swap{}, fmt.Errorf("corrupted send broker address: %v", blob.BrokerSendTokenAddr)
		}
		brokerFee = new(big.Int).Div(new(big.Int).Mul(value, big.NewInt(blob.BrokerFee)), big.NewInt(10000))
	}

	secretHash, err := unmarshalSecretHash(blob.SecretHash)
	if err != nil {
		return swap.Swap{}, err
	}

	return swap.Swap{
		ID:              blob.ID,
		Token:           token,
		Value:           value,
		Speed:           blob.Speed,
		SecretHash:      secretHash,
		TimeLock:        blob.TimeLock,
		SpendingAddress: blob.SendTo,
		FundingAddress:  fundingAddress,
		BrokerAddress:   blob.BrokerSendTokenAddr,
		BrokerFee:       brokerFee,
	}, nil
}

func (builder *builder) buildForeignSwap(blob swap.SwapBlob, timelock int64, spendingAddress string) (swap.Swap, error) {
	token, err := tokens.PatchToken(string(blob.ReceiveToken))
	if err != nil {
		return swap.Swap{}, err
	}

	value, ok := new(big.Int).SetString(blob.ReceiveAmount, 10)
	if !ok {
		return swap.Swap{}, fmt.Errorf("corrupted receive value: %v", blob.ReceiveAmount)
	}

	brokerFee := big.NewInt(0)
	if blob.BrokerFee != 0 {
		if err := builder.Wallet.VerifyAddress(token.Blockchain, blob.BrokerReceiveTokenAddr); err != nil {
			return swap.Swap{}, fmt.Errorf("corrupted receive broker address: %v", blob.BrokerReceiveTokenAddr)
		}
		brokerFee = new(big.Int).Div(new(big.Int).Mul(value, big.NewInt(blob.BrokerFee)), big.NewInt(10000))
	}

	secretHash, err := unmarshalSecretHash(blob.SecretHash)
	if err != nil {
		return swap.Swap{}, err
	}

	withdrawAddress := spendingAddress
	if blob.WithdrawAddress != "" {
		withdrawAddress = blob.WithdrawAddress
	}

	return swap.Swap{
		ID:              blob.ID,
		Token:           token,
		Value:           value,
		Speed:           blob.Speed,
		SecretHash:      secretHash,
		TimeLock:        blob.TimeLock,
		SpendingAddress: spendingAddress,
		FundingAddress:  blob.ReceiveFrom,
		WithdrawAddress: withdrawAddress,
		BrokerAddress:   blob.BrokerReceiveTokenAddr,
		BrokerFee:       brokerFee,
	}, nil
}

func (builder *builder) calculateTimeLocks(swap swap.SwapBlob) (native, foreign int64) {
	if swap.ShouldInitiateFirst {
		native = swap.TimeLock
		foreign = swap.TimeLock - 24*60*60
		return
	}
	native = swap.TimeLock - 24*60*60
	foreign = swap.TimeLock
	return
}

func (builder *builder) calculateAddresses(swap swap.SwapBlob) (string, string, error) {
	sendToken, err := tokens.PatchToken(swap.SendToken)
	if err != nil {
		return "", "", err
	}
	sendAddress, err := builder.Wallet.GetAddress(swap.Password, sendToken.Blockchain)
	if err != nil {
		return "", "", err
	}
	receiveToken, err := tokens.PatchToken(swap.ReceiveToken)
	if err != nil {
		return "", "", err
	}
	receiveAddress, err := builder.Wallet.GetAddress(swap.Password, receiveToken.Blockchain)
	if err != nil {
		return "", "", err
	}
	return sendAddress, receiveAddress, nil
}

func unmarshalSecretHash(secretHash string) ([32]byte, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(secretHash)
	if err != nil {
		return [32]byte{}, err
	}

	if len(hashBytes) != 32 {
		return [32]byte{}, fmt.Errorf("invalid secret hash")
	}

	hash32 := [32]byte{}
	copy(hash32[:], hashBytes)
	return hash32, nil
}
