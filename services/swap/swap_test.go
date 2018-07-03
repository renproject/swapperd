package swap_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/atom-go/domains/match"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	. "github.com/republicprotocol/atom-go/services/swap"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/keystore"

	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap", func() {

	var aliceSwap, bobSwap Swap

	BeforeSuite(func() {

		var aliceInfo, bobInfo Info
		var aliceNet, bobNet Network
		var aliceOrder, bobOrder match.Match
		var aliceOrderID, bobOrderID [32]byte
		var aliceSendValue, bobSendValue *big.Int
		var aliceRecieveValue, bobRecieveValue *big.Int
		var aliceCurrency, bobCurrency uint32
		var alice, bob *bind.TransactOpts
		var aliceBitcoinAddress, bobBitcoinAddress string
		var swapID [32]byte

		rand.Read(aliceOrderID[:])
		rand.Read(bobOrderID[:])

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

		aliceRecieveValue = big.NewInt(99990000)
		bobRecieveValue = big.NewInt(8000000)

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

		aliceOrder = match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, aliceRecieveValue, aliceCurrency, bobCurrency)
		bobOrder = match.NewMatch(bobOrderID, aliceOrderID, bobSendValue, bobRecieveValue, bobCurrency, aliceCurrency)

		aliceInfo.SetOwnerAddress(aliceOrderID, []byte(aliceBitcoinAddress))
		bobInfo.SetOwnerAddress(bobOrderID, bob.From.Bytes())

		reqAlice, err := eth.NewEthereumRequestAtom(ganache, alice)
		Expect(err).Should(BeNil())

		reqBob := btc.NewBitcoinAtomRequester(connection, bobBitcoinAddress)
		resAlice := btc.NewBitcoinAtomResponder(connection, aliceBitcoinAddress)

		resBob, err := eth.NewEthereumResponseAtom(ganache, bob)
		Expect(err).Should(BeNil())

		aliceSwap = NewSwap(reqAlice, resAlice, aliceInfo, aliceOrder, aliceNet)
		bobSwap = NewSwap(reqBob, resBob, bobInfo, bobOrder, bobNet)
	})

	It("can do an eth - btc atomic swap", func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := aliceSwap.Execute()
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 1")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := bobSwap.Execute()
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 2")
		}()

		wg.Wait()
	})
})
