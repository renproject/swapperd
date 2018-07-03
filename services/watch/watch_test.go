package watch_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	"github.com/republicprotocol/atom-go/services/swap"
	. "github.com/republicprotocol/atom-go/services/watch"

	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"
	wal "github.com/republicprotocol/atom-go/adapters/wallet/eth"
	"github.com/republicprotocol/atom-go/domains/match"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap using Watch", func() {

	var aliceWatch, bobWatch Watch
	var aliceOrderID, bobOrderID [32]byte

	rand.Read(aliceOrderID[:])
	rand.Read(bobOrderID[:])

	BeforeSuite(func() {
		var aliceInfo, bobInfo swap.Info
		var aliceNet, bobNet swap.Network
		var aliceSendValue, bobSendValue *big.Int
		var aliceCurrency, bobCurrency uint32
		var alice, bob *bind.TransactOpts
		var aliceBitcoinAddress, bobBitcoinAddress string
		var swapID [32]byte

		rand.Read(swapID[:])

		aliceCurrency = 1
		bobCurrency = 0

		var confPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/config.json"
		var ksPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/keystore.json"
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		keystore := keystore.NewKeystore(ksPath)

		ganache, err := ethclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		ownerECDSA, err := keystore.LoadKeypair("ethereum")
		Expect(err).ShouldNot(HaveOccurred())
		owner := bind.NewKeyedTransactor(ownerECDSA)

		_, alice, err = ganache.NewAccount(1000000000000000000, owner)
		Expect(err).ShouldNot(HaveOccurred())
		alice.GasLimit = 3000000

		_, bob, err = ganache.NewAccount(1000000000000000000, owner)
		Expect(err).ShouldNot(HaveOccurred())
		bob.GasLimit = 3000000

		time.Sleep(5 * time.Second)
		connection, err := btcclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		aliceSendValue = big.NewInt(10000000)
		bobSendValue = big.NewInt(10000000)

		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		time.Sleep(5 * time.Second)

		aliceAddr, err := regtest.GetAddressForAccount(connection, "alice")
		Expect(err).ShouldNot(HaveOccurred())
		aliceBitcoinAddress = aliceAddr.EncodeAddress()

		bobAddr, err := regtest.GetAddressForAccount(connection, "bob")
		Expect(err).ShouldNot(HaveOccurred())
		bobBitcoinAddress = bobAddr.EncodeAddress()
		Expect(err).Should(BeNil())

		aliceNet, err = net.NewEthereumNetwork(ganache, alice)
		Expect(err).Should(BeNil())

		bobNet, err = net.NewEthereumNetwork(ganache, bob)
		Expect(err).Should(BeNil())

		aliceInfo, err = ax.NewEtereumAtomInfo(ganache, alice)
		Expect(err).Should(BeNil())

		bobInfo, err = ax.NewEtereumAtomInfo(ganache, bob)
		Expect(err).Should(BeNil())

		atomMatch := match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, bobSendValue, aliceCurrency, bobCurrency)
		mockWallet, err := wal.NewEthereumWallet(ganache, *owner)
		Expect(err).Should(BeNil())

		err = mockWallet.SetMatch(atomMatch)
		Expect(err).Should(BeNil())

		aliceInfo.SetOwnerAddress(aliceOrderID, []byte(aliceBitcoinAddress))
		bobInfo.SetOwnerAddress(bobOrderID, bob.From.Bytes())

		reqAlice, err := eth.NewEthereumRequestAtom(ganache, alice)
		Expect(err).Should(BeNil())

		reqBob := btc.NewBitcoinAtomRequester(connection, bobBitcoinAddress)
		resAlice := btc.NewBitcoinAtomResponder(connection, aliceBitcoinAddress)

		resBob, err := eth.NewEthereumResponseAtom(ganache, bob)
		Expect(err).Should(BeNil())

		reqAlice, err = eth.NewEthereumRequestAtom(ganache, alice)
		Expect(err).Should(BeNil())

		reqBob = btc.NewBitcoinAtomRequester(connection, bobBitcoinAddress)
		resAlice = btc.NewBitcoinAtomResponder(connection, aliceBitcoinAddress)

		resBob, err = eth.NewEthereumResponseAtom(ganache, bob)
		Expect(err).Should(BeNil())

		aliceWatch = NewWatch(aliceNet, aliceInfo, mockWallet, reqAlice, resAlice)
		bobWatch = NewWatch(bobNet, bobInfo, mockWallet, reqBob, resBob)
	})

	It("can do an eth - btc atomic swap", func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := aliceWatch.Run(aliceOrderID)
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 1")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := bobWatch.Run(bobOrderID)
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 2")
		}()

		wg.Wait()
	})
})
