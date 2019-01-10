package binder

import (
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/adapter/binder/btc"
	"github.com/republicprotocol/swapperd/adapter/binder/erc20"
	"github.com/republicprotocol/swapperd/adapter/binder/eth"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
)

type builder struct {
	wallet.Wallet
	logrus.FieldLogger
}

func NewBuilder(wallet wallet.Wallet, logger logrus.FieldLogger) swapper.ContractBuilder {
	return &builder{
		wallet,
		logger,
	}
}

func (builder *builder) BuildSwapContracts(req swapper.SwapRequest) (swapper.Contract, swapper.Contract, error) {
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

func (builder *builder) buildBinder(swap swap.Swap, cost blockchain.Cost, password string) (swapper.Contract, error) {
	switch swap.Token {
	case blockchain.TokenBTC:
		btcAccount, err := builder.BitcoinAccount(password)
		if err != nil {
			return nil, err
		}
		return btc.NewBTCSwapContractBinder(btcAccount, swap, cost, builder.FieldLogger)
	case blockchain.TokenETH:
		ethAccount, err := builder.EthereumAccount(password)
		if err != nil {
			return nil, err
		}
		return eth.NewETHSwapContractBinder(ethAccount, swap, cost, builder.FieldLogger)
	case blockchain.TokenWBTC, blockchain.TokenDGX, blockchain.TokenREN,
		blockchain.TokenTUSD, blockchain.TokenOMG, blockchain.TokenZRX:
		ethAccount, err := builder.EthereumAccount(password)
		if err != nil {
			return nil, err
		}
		return erc20.NewERC20SwapContractBinder(ethAccount, swap, cost, builder.FieldLogger)
	default:
		return nil, blockchain.NewErrUnsupportedToken(swap.Token.Name)
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
	token, err := blockchain.PatchToken(blob.SendToken)
	if err != nil {
		return swap.Swap{}, err
	}
	value, ok := new(big.Int).SetString(blob.SendAmount, 10)
	if !ok {
		return swap.Swap{}, fmt.Errorf("corrupted send value: %v", blob.SendAmount)
	}

	fee, ok := new(big.Int).SetString(blob.SendFee, 10)
	if !ok {
		fee, err = builder.Wallet.DefaultFee(token.Blockchain)
		if err != nil {
			return swap.Swap{}, fmt.Errorf("failed to get default fee: %v", err)
		}
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
		Fee:             fee,
		SecretHash:      secretHash,
		TimeLock:        blob.TimeLock,
		SpendingAddress: blob.SendTo,
		FundingAddress:  fundingAddress,
		BrokerAddress:   blob.BrokerSendTokenAddr,
		BrokerFee:       brokerFee,
	}, nil
}

func (builder *builder) buildForeignSwap(blob swap.SwapBlob, timelock int64, spendingAddress string) (swap.Swap, error) {
	token, err := blockchain.PatchToken(blob.ReceiveToken)
	if err != nil {
		return swap.Swap{}, err
	}

	value, ok := new(big.Int).SetString(blob.ReceiveAmount, 10)
	if !ok {
		return swap.Swap{}, fmt.Errorf("corrupted receive value: %v", blob.ReceiveAmount)
	}

	fee, ok := new(big.Int).SetString(blob.ReceiveFee, 10)
	if !ok {
		fee, err = builder.Wallet.DefaultFee(token.Blockchain)
		if err != nil {
			return swap.Swap{}, fmt.Errorf("failed to get default fee: %v", err)
		}
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

	return swap.Swap{
		ID:              blob.ID,
		Token:           token,
		Value:           value,
		Fee:             fee,
		SecretHash:      secretHash,
		TimeLock:        blob.TimeLock,
		SpendingAddress: spendingAddress,
		FundingAddress:  blob.ReceiveFrom,
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
	sendToken, err := blockchain.PatchToken(swap.SendToken)
	if err != nil {
		return "", "", err
	}

	receiveToken, err := blockchain.PatchToken(swap.ReceiveToken)
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

	if sendToken.Blockchain == blockchain.Ethereum && receiveToken.Blockchain == blockchain.Bitcoin {
		return ethAddress.String(), btcAddress.EncodeAddress(), nil
	}

	if sendToken.Blockchain == blockchain.Bitcoin && receiveToken.Blockchain == blockchain.Ethereum {
		return btcAddress.EncodeAddress(), ethAddress.String(), nil
	}

	if sendToken.Blockchain == blockchain.Ethereum && receiveToken.Blockchain == blockchain.Ethereum {
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
