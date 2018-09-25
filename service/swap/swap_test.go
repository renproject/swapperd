package swap_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/renex-swapper-go/service/swap"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/swap"
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"

	configDriver "github.com/republicprotocol/renex-swapper-go/driver/config"
	loggerDriver "github.com/republicprotocol/renex-swapper-go/driver/logger"
	"github.com/republicprotocol/renex-swapper-go/driver/watchdog"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap", func() {

	// TODO: Fix the tests
	buildConfigs := func() (config.Config, keystore.Keystore, keystore.Keystore) {
		config := configDriver.New("", "nightly")
		keys := utils.LoadTestKeys("../../secrets/test.json")

		btcKeyA, err := keystore.NewBitcoinKey(keys.Alice.Bitcoin, "mainnet")
		Expect(err).ShouldNot(HaveOccurred())

		ethPrivKeyA, err := crypto.HexToECDSA(keys.Alice.Ethereum)
		Expect(err).ShouldNot(HaveOccurred())

		ethKeyA, err := keystore.NewEthereumKey(ethPrivKeyA, "kovan")
		Expect(err).ShouldNot(HaveOccurred())

		btcKeyB, err := keystore.NewBitcoinKey(keys.Bob.Bitcoin, "mainnet")
		Expect(err).ShouldNot(HaveOccurred())

		ethPrivKeyB, err := crypto.HexToECDSA(keys.Bob.Ethereum)
		Expect(err).ShouldNot(HaveOccurred())

		ethKeyB, err := keystore.NewEthereumKey(ethPrivKeyB, "kovan")
		Expect(err).ShouldNot(HaveOccurred())

		ksA := keystore.New(btcKeyA, ethKeyA)
		ksB := keystore.New(btcKeyB, ethKeyB)

		return config, ksB, ksA
	}

	buildRequests := func(ksA, ksB keystore.Keystore) (swapDomain.Request, swapDomain.Request) {

		var aliceUID, bobUID [32]byte
		var aliceSecret [32]byte
		rand.Read(aliceUID[:])
		rand.Read(bobUID[:])
		rand.Read(aliceSecret[:])
		aliceSecretHash := sha256.Sum256(aliceSecret[:])
		timelock := time.Now().Unix() + 48*60*60

		aliceReq := swapDomain.Request{
			UID:                aliceUID,
			TimeLock:           timelock,
			Secret:             aliceSecret,
			SecretHash:         aliceSecretHash,
			SendToAddress:      ksB.GetKey(token.ETH).(keystore.EthereumKey).Address.String(),
			ReceiveFromAddress: ksB.GetKey(token.BTC).(keystore.BitcoinKey).AddressString,
			SendValue:          big.NewInt(20000),
			ReceiveValue:       big.NewInt(20000),
			SendToken:          token.ETH,
			ReceiveToken:       token.BTC,
			GoesFirst:          true,
		}

		areq, err := json.Marshal(aliceReq)
		Expect(err).ShouldNot(HaveOccurred())
		fmt.Println("Alice Request: ", hex.EncodeToString(areq))

		bobReq := swapDomain.Request{
			UID:                bobUID,
			TimeLock:           timelock,
			Secret:             [32]byte{},
			SecretHash:         aliceSecretHash,
			SendToAddress:      ksA.GetKey(token.BTC).(keystore.BitcoinKey).AddressString,
			ReceiveFromAddress: ksA.GetKey(token.ETH).(keystore.EthereumKey).Address.String(),
			SendValue:          big.NewInt(20000),
			ReceiveValue:       big.NewInt(20000),
			SendToken:          token.BTC,
			ReceiveToken:       token.ETH,
			GoesFirst:          false,
		}

		breq, err := json.Marshal(bobReq)
		Expect(err).ShouldNot(HaveOccurred())
		fmt.Println("Bob Request: ", hex.EncodeToString(breq))

		return aliceReq, bobReq
	}

	buildSwappers := func(cfg config.Config, ksA, ksB keystore.Keystore) (Swapper, Swapper) {
		wd := watchdog.NewMock()
		loggr := loggerDriver.NewStdOut()
		swapperAlice := NewSwapper(swap.New(cfg, ksA, wd, loggr))
		swapperBob := NewSwapper(swap.New(cfg, ksB, wd, loggr))
		return swapperAlice, swapperBob
	}

	buildSwaps := func() (Swap, Swap) {
		conf, ksA, ksB := buildConfigs()
		reqA, reqB := buildRequests(ksA, ksB)
		aliceSwapper, bobSwapper := buildSwappers(conf, ksA, ksB)
		aliceSwap, err := aliceSwapper.NewSwap(reqA)
		Expect(err).Should(BeNil())
		bobSwap, err := bobSwapper.NewSwap(reqB)
		Expect(err).Should(BeNil())
		return aliceSwap, bobSwap
	}

	It("can do an eth - btc atomic swap", func() {
		aliceSwap, bobSwap := buildSwaps()

		// TODO: Use co library
		// co.ParBegin(
		// 	func() {
		// 		defer GinkgoRecover()
		// 		err := aliceSwap.Execute()
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 	},
		// 	func() {
		// 		defer GinkgoRecover()
		// 		err := bobSwap.Execute()
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 	})

		////

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := aliceSwap.Execute()
			Expect(err).ShouldNot(HaveOccurred())
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := bobSwap.Execute()
			Expect(err).ShouldNot(HaveOccurred())
		}()
		wg.Wait()
	})
})
