package binder

import (
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/adapter/binder/btc"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/adapter/binder/eth"
	"github.com/republicprotocol/swapperd/adapter/fund"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/sirupsen/logrus"
)

type builder struct {
	fund.Manager
	logrus.FieldLogger
}

func NewBuilder(manager fund.Manager, logger logrus.FieldLogger) swapper.ContractBuilder {
	return &builder{
		manager,
		logger,
	}
}

func (builder *builder) BuildSwapContracts(swap foundation.SwapRequest) (swapper.Contract, swapper.Contract, error) {
	native, foreign, err := builder.buildComplementarySwaps(swap)
	if err != nil {
		return nil, nil, err
	}
	nativeBinder, err := builder.buildBinder(native, swap.Password)
	if err != nil {
		return nil, nil, err
	}
	foreignBinder, err := builder.buildBinder(foreign, swap.Password)
	if err != nil {
		return nil, nil, err
	}
	return nativeBinder, foreignBinder, nil
}

func (builder *builder) buildBinder(swap foundation.Swap, password string) (swapper.Contract, error) {
	switch swap.Token {
	case foundation.TokenBTC:
		btcAccount, err := builder.BitcoinAccount(password)
		if err != nil {
			return nil, err
		}
		return btc.NewBTCSwapContractBinder(btcAccount, swap, builder.FieldLogger)
	case foundation.TokenETH:
		ethAccount, err := builder.EthereumAccount(password)
		if err != nil {
			return nil, err
		}
		return eth.NewETHSwapContractBinder(ethAccount, swap, builder.FieldLogger)
	case foundation.TokenWBTC:
		ethAccount, err := builder.EthereumAccount(password)
		if err != nil {
			return nil, err
		}
		return erc20.NewERC20SwapContractBinder(ethAccount, swap, builder.FieldLogger)
	default:
		return nil, foundation.NewErrUnsupportedToken(swap.Token.Name)
	}
}

func (builder *builder) buildComplementarySwaps(swap foundation.SwapRequest) (foundation.Swap, foundation.Swap, error) {
	fundingAddr, spendingAddr, err := builder.calculateAddresses(swap)
	if err != nil {
		return foundation.Swap{}, foundation.Swap{}, err
	}

	nativeExpiry, foreignExpiry := builder.calculateTimeLocks(swap.SwapBlob)

	nativeSwap, err := builder.buildNativeSwap(swap.SwapBlob, nativeExpiry, fundingAddr)
	if err != nil {
		return foundation.Swap{}, foundation.Swap{}, err
	}
	foreignSwap, err := builder.buildForeignSwap(swap.SwapBlob, foreignExpiry, spendingAddr)
	if err != nil {
		return foundation.Swap{}, foundation.Swap{}, err
	}
	return nativeSwap, foreignSwap, nil
}

func (builder *builder) buildNativeSwap(swap foundation.SwapBlob, timelock int64, fundingAddress string) (foundation.Swap, error) {
	token, err := foundation.PatchToken(swap.SendToken)
	if err != nil {
		return foundation.Swap{}, err
	}
	value, ok := big.NewInt(0).SetString(swap.SendAmount, 10)
	if !ok {
		return foundation.Swap{}, fmt.Errorf("corrupted send value: %v", swap.SendAmount)
	}

	fee, ok := big.NewInt(0).SetString(swap.SendFee, 10)
	if !ok {
		return foundation.Swap{}, fmt.Errorf("corrupted send fee: %v", swap.SendFee)
	}

	brokerFee := new(big.Int).Div(new(big.Int).Mul(value, big.NewInt(swap.BrokerFee)), big.NewInt(10000))

	secretHash, err := unmarshalSecretHash(swap.SecretHash)
	if err != nil {
		return foundation.Swap{}, err
	}

	return foundation.Swap{
		ID:              swap.ID,
		Token:           token,
		Value:           value,
		Fee:             fee,
		SecretHash:      secretHash,
		TimeLock:        swap.TimeLock,
		SpendingAddress: swap.SendTo,
		FundingAddress:  fundingAddress,
		BrokerAddress:   swap.BrokerSendTokenAddr,
		BrokerFee:       brokerFee,
	}, nil
}

func (builder *builder) buildForeignSwap(swap foundation.SwapBlob, timelock int64, spendingAddress string) (foundation.Swap, error) {
	token, err := foundation.PatchToken(swap.ReceiveToken)
	if err != nil {
		return foundation.Swap{}, err
	}

	value, ok := big.NewInt(0).SetString(swap.ReceiveAmount, 10)
	if !ok {
		return foundation.Swap{}, fmt.Errorf("corrupted receive value: %v", swap.ReceiveAmount)
	}

	fee, ok := big.NewInt(0).SetString(swap.ReceiveFee, 10)
	if !ok {
		return foundation.Swap{}, fmt.Errorf("corrupted receive fee: %v", swap.ReceiveFee)
	}

	brokerFee := new(big.Int).Div(new(big.Int).Mul(value, big.NewInt(swap.BrokerFee)), big.NewInt(10000))

	secretHash, err := unmarshalSecretHash(swap.SecretHash)
	if err != nil {
		return foundation.Swap{}, err
	}

	return foundation.Swap{
		ID:              swap.ID,
		Token:           token,
		Value:           value,
		Fee:             fee,
		SecretHash:      secretHash,
		TimeLock:        swap.TimeLock,
		SpendingAddress: spendingAddress,
		FundingAddress:  swap.ReceiveFrom,
		BrokerAddress:   swap.BrokerReceiveTokenAddr,
		BrokerFee:       brokerFee,
	}, nil
}

func (builder *builder) calculateTimeLocks(swap foundation.SwapBlob) (native, foreign int64) {
	if swap.ShouldInitiateFirst {
		native = swap.TimeLock
		foreign = swap.TimeLock - 24*60*60
		return
	}
	native = swap.TimeLock - 24*60*60
	foreign = swap.TimeLock
	return
}

func (builder *builder) calculateAddresses(swap foundation.SwapRequest) (string, string, error) {
	sendToken, err := foundation.PatchToken(swap.SendToken)
	if err != nil {
		return "", "", err
	}

	receiveToken, err := foundation.PatchToken(swap.ReceiveToken)
	if err != nil {
		return "", "", err
	}

	ethAccount, err := builder.EthereumAccount(swap.Password)
	if err != nil {
		return "", "", err
	}

	btcAccount, err := builder.BitcoinAccount(swap.Password)
	if err != nil {
		return "", "", err
	}

	ethAddress := ethAccount.Address()
	btcAddress, err := btcAccount.Address()
	if err != nil {
		return "", "", err
	}

	if sendToken.Blockchain == foundation.Ethereum && receiveToken.Blockchain == foundation.Bitcoin {
		return ethAddress.String(), btcAddress.EncodeAddress(), nil
	}

	if sendToken.Blockchain == foundation.Bitcoin && receiveToken.Blockchain == foundation.Ethereum {
		return btcAddress.EncodeAddress(), ethAddress.String(), nil
	}

	if sendToken.Blockchain == foundation.Ethereum && receiveToken.Blockchain == foundation.Ethereum {
		return ethAddress.String(), ethAddress.String(), nil
	}

	return "", "", fmt.Errorf("unsupported blockchain pairing: %s <=> %s", sendToken.Blockchain, receiveToken.Blockchain)
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
