package binder_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/swapperd/adapter/binder"
	"github.com/renproject/swapperd/adapter/wallet"
	"github.com/renproject/swapperd/core/wallet/swapper/immediate"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

var _ = Describe("Binder Tests", func() {
	logger := logrus.StandardLogger()

	buildContractBuilder := func() (wallet.Wallet, immediate.ContractBuilder) {
		config := wallet.Testnet
		config.Mnemonic = "weird" // os.Getenv("MNEMONIC")
		wallet := wallet.New(config, logger)
		builder := NewBuilder(wallet, logger)
		return wallet, builder
	}

	buildSwaps := func(wallet wallet.Wallet, sendToken, receiveToken tokens.Token) (immediate.SwapRequest, immediate.SwapRequest) {
		password := "hello" // os.Getenv("PASSWORD")
		sendTokenAddr, err := wallet.GetAddress(password, sendToken.Blockchain)
		if err != nil {
			panic(err)
		}
		receiveTokenAddr, err := wallet.GetAddress(password, receiveToken.Blockchain)
		if err != nil {
			panic(err)
		}

		initID := [32]byte{}
		redeemID := [32]byte{}
		rand.Read(initID[:])
		rand.Read(redeemID[:])
		initSwapID := swap.SwapID(base64.StdEncoding.EncodeToString(initID[:]))
		redeemSwapID := swap.SwapID(base64.StdEncoding.EncodeToString(redeemID[:]))

		secret := sha3.Sum256(append([]byte(password), []byte(initSwapID)...))
		secretHash32 := sha256.Sum256(secret[:])
		secretHash := base64.StdEncoding.EncodeToString(secretHash32[:])
		timeLock := time.Now().Unix() + 3600

		val, err := rand.Int(rand.Reader, big.NewInt(40000))
		if err != nil {
			panic(err)
		}

		sndAmt := new(big.Int).Add(big.NewInt(20000), val)
		rcvAmt := new(big.Int).Add(big.NewInt(20000), val)

		initSwapBlob := swap.SwapBlob{
			ID:                  initSwapID,
			SendToken:           sendToken.Name,
			ReceiveToken:        receiveToken.Name,
			SendTo:              sendTokenAddr,
			ReceiveFrom:         receiveTokenAddr,
			SendAmount:          sndAmt.String(),
			ReceiveAmount:       rcvAmt.String(),
			SecretHash:          secretHash,
			TimeLock:            timeLock,
			Password:            password,
			ShouldInitiateFirst: true,
		}

		respSwapBlob := swap.SwapBlob{
			ID:                  redeemSwapID,
			SendToken:           receiveToken.Name,
			ReceiveToken:        sendToken.Name,
			SendTo:              receiveTokenAddr,
			ReceiveFrom:         sendTokenAddr,
			SendAmount:          rcvAmt.String(),
			ReceiveAmount:       sndAmt.String(),
			SecretHash:          secretHash,
			TimeLock:            timeLock,
			Password:            password,
			ShouldInitiateFirst: false,
		}

		wallet.LockBalance(receiveToken.Name, rcvAmt.String())
		wallet.LockBalance(sendToken.Name, sndAmt.String())

		return immediate.NewSwapRequest(initSwapBlob, blockchain.Cost{}, blockchain.Cost{}), immediate.NewSwapRequest(respSwapBlob, blockchain.Cost{}, blockchain.Cost{})
	}

	// for _, sendToken := range tokens.SupportedTokens {
	// 	for _, receiveToken := range tokens.SupportedTokens {
	// 		if sendToken.Name == receiveToken.Name {
	// 			continue
	// 		}
	Context(fmt.Sprintf("when swapping between %s and %s", "ZEC", "ETH"), func() {
		It("should successfully do an atomic swap", func() {
			wallet, builder := buildContractBuilder()
			sndReq, rcvReq := buildSwaps(wallet, tokens.ZEC, tokens.ETH)
			swapper := immediate.New(2048, builder, wallet, logger)
			swapper.Send(sndReq)
			swapper.Send(rcvReq)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			counter := 0
			tau.New(tau.NewIO(2048), tau.ReduceFunc(func(msg tau.Message) tau.Message {
				switch msg := msg.(type) {
				case immediate.SwapRequest:
					swapper.Send(msg)
				case tau.Error:
					logger.Error(msg)
				case immediate.DeleteSwap:
					counter++
					if counter == 2 {
						cancel()
					}
				default:
					swapper.Send(tau.NewTick(time.Now()))
				}
				return nil
			}), swapper).Run(ctx.Done())
			Expect(counter).Should(Equal(2))
		})
	})
	// 	}
	// }
})
