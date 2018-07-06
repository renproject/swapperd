package swap_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/atom-go/domains/match"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	. "github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/atom-go/services/watch"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"

	ethKey "github.com/republicprotocol/atom-go/adapters/key/eth"

	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"

	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"
	"github.com/republicprotocol/atom-go/adapters/store/leveldb"
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
		var alice, bob *ecdsa.PrivateKey
		var aliceKey, bobKey swap.Key
		var aliceBitcoinAddress, bobBitcoinAddress string
		var swapID [32]byte

		rand.Read(aliceOrderID[:])
		rand.Read(bobOrderID[:])

		rand.Read(swapID[:])

		aliceCurrency = 1
		bobCurrency = 0

		var confPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/config.json"
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())

		ganache, err := ethclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		pk, err := crypto.HexToECDSA("2aba04ee8a322b8648af2a784144181a0c793f1a2e80519418f3d20bbfb22249")
		Expect(err).ShouldNot(HaveOccurred())
		owner := bind.NewKeyedTransactor(pk)

		alice, err = crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		aliceKey, err = ethKey.NewEthereumKey(hex.EncodeToString(crypto.FromECDSA(alice)), "ganache")
		Expect(err).ShouldNot(HaveOccurred())

		bob, err = crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		bobKey, err = ethKey.NewEthereumKey(hex.EncodeToString(crypto.FromECDSA(bob)), "ganache")
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddrBytes, err := aliceKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		bobAddrBytes, err := bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())

		err = ganache.Transfer(common.BytesToAddress(aliceAddrBytes), owner, 1000000000000000000)
		Expect(err).ShouldNot(HaveOccurred())

		err = ganache.Transfer(common.BytesToAddress(bobAddrBytes), owner, 1000000000000000000)
		Expect(err).ShouldNot(HaveOccurred())

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

		aliceNet, err = net.NewEthereumNetwork(ganache, aliceAuth)
		Expect(err).Should(BeNil())

		bobNet, err = net.NewEthereumNetwork(ganache, bobAuth)
		Expect(err).Should(BeNil())

		aliceInfo, err = ax.NewEtereumAtomInfo(ganache, aliceAuth)
		Expect(err).Should(BeNil())

		bobInfo, err = ax.NewEtereumAtomInfo(ganache, bobAuth)
		Expect(err).Should(BeNil())

		aliceOrder = match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, aliceRecieveValue, aliceCurrency, bobCurrency)
		bobOrder = match.NewMatch(bobOrderID, aliceOrderID, bobSendValue, bobRecieveValue, bobCurrency, aliceCurrency)

		aliceInfo.SetOwnerAddress(aliceOrderID, []byte(aliceBitcoinAddress))
		bobInfo.SetOwnerAddress(bobOrderID, bob.From.Bytes())

		reqAlice, err := eth.NewEthereumAtom(ganache, alice)
		Expect(err).Should(BeNil())

		reqBob := btc.NewBitcoinAtom(connection, bobBitcoinAddress)
		resAlice := btc.NewBitcoinAtom(connection, aliceBitcoinAddress)

		resBob, err := eth.NewEthereumAtom(ganache, bob)
		Expect(err).Should(BeNil())

		aliceStr := NewSwapStore(leveldb.NewLDBStore("/db"))
		bobStr := NewSwapStore(leveldb.NewLDBStore("/db"))

		aliceSwap = NewSwap(reqAlice, resAlice, aliceInfo, aliceOrder, aliceNet, aliceStr)
		bobSwap = NewSwap(reqBob, resBob, bobInfo, bobOrder, bobNet, bobStr)
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
