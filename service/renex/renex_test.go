package renex_test

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	renexAdapter "github.com/republicprotocol/renex-swapper-go/adapter/renex"
	stateAdapter "github.com/republicprotocol/renex-swapper-go/adapter/state"

	configDriver "github.com/republicprotocol/renex-swapper-go/driver/config"
	loggerDriver "github.com/republicprotocol/renex-swapper-go/driver/logger"
	networkDriver "github.com/republicprotocol/renex-swapper-go/driver/network"
	storeDriver "github.com/republicprotocol/renex-swapper-go/driver/store"
	"github.com/republicprotocol/renex-swapper-go/driver/watchdog"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap", func() {

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

		return config, ksA, ksB
	}

	buildMatches := func() (match.Match, match.Match) {
		var aliceOrderID, bobOrderID [32]byte
		rand.Read(aliceOrderID[:])
		rand.Read(bobOrderID[:])

		aliceCurrency := token.ETH
		bobCurrency := token.BTC

		aliceSendValue := big.NewInt(100000)
		bobSendValue := big.NewInt(100000)
		aliceReceiveValue := big.NewInt(100000)
		bobReceiveValue := big.NewInt(100000)

		aliceOrder := match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, aliceReceiveValue, aliceCurrency, bobCurrency)
		bobOrder := match.NewMatch(bobOrderID, aliceOrderID, bobSendValue, bobReceiveValue, bobCurrency, aliceCurrency)

		return aliceOrder, bobOrder
	}

	sendAddresses := func(aliceOrderID, bobOrderID [32]byte, aliceKS, bobKS keystore.Keystore, net network.Network) {
		AliceBtcKey := aliceKS.GetKey(token.BTC).(keystore.BitcoinKey)
		BobEthKey := bobKS.GetKey(token.ETH).(keystore.EthereumKey)
		err := net.SendOwnerAddress(aliceOrderID, []byte(AliceBtcKey.AddressString))
		Expect(err).Should(BeNil())
		err = net.SendOwnerAddress(bobOrderID, []byte(BobEthKey.Address.String()))
		Expect(err).Should(BeNil())
	}

	buildRenExWatchers := func(aliceMatch, bobMatch match.Match, cfg config.Config, ksA, ksB keystore.Keystore, net network.Network) (RenEx, RenEx) {
		wd := watchdog.NewMock()
		loggr := loggerDriver.NewStdOut()

		aliceLDB, err := storeDriver.NewLevelDB("/Users/susruth/go/src/github.com/republicprotocol/renex-swapper-go/temp/dbAlice")
		Expect(err).Should(BeNil())

		bobLDB, err := storeDriver.NewLevelDB("/Users/susruth/go/src/github.com/republicprotocol/renex-swapper-go/temp/dbBob")
		Expect(err).Should(BeNil())

		aliceState := state.NewState(stateAdapter.New(aliceLDB, loggr))
		bobState := state.NewState(stateAdapter.New(bobLDB, loggr))

		mockBinder, err := renexAdapter.NewMockBinder(aliceMatch, bobMatch)
		Expect(err).Should(BeNil())

		aliceRenEx := renexAdapter.New(cfg, ksA, net, wd, aliceState, loggr, mockBinder)
		bobRenEx := renexAdapter.New(cfg, ksB, net, wd, bobState, loggr, mockBinder)

		return NewRenEx(aliceRenEx), NewRenEx(bobRenEx)
	}

	buildRenExSwappers := func() ([32]byte, [32]byte, RenEx, RenEx) {
		conf, ksA, ksB := buildConfigs()
		matchA, matchB := buildMatches()
		net := networkDriver.NewMock()
		sendAddresses(matchA.PersonalOrderID(), matchB.PersonalOrderID(), ksA, ksB, net)
		alice, bob := buildRenExWatchers(matchA, matchB, conf, ksA, ksB, net)
		return matchA.PersonalOrderID(), matchB.PersonalOrderID(), alice, bob
	}

	It("can do an eth - btc atomic swap using renex", func() {
		aliceID, bobID, alice, bob := buildRenExSwappers()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			err := <-alice.Start()
			Expect(err).ShouldNot(HaveOccurred())
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			err := <-bob.Start()
			Expect(err).ShouldNot(HaveOccurred())
		}()

		alice.Add(aliceID)
		alice.Notify()
		bob.Add(bobID)
		bob.Notify()

		wg.Wait()
	})
})
