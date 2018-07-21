package watch_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/adapters/atoms"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/owner"
	"github.com/republicprotocol/atom-go/adapters/store/leveldb"
	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
	. "github.com/republicprotocol/atom-go/services/watch"

	btcKey "github.com/republicprotocol/atom-go/adapters/key/btc"
	ethKey "github.com/republicprotocol/atom-go/adapters/key/eth"

	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"

	wal "github.com/republicprotocol/atom-go/adapters/wallet/mock"
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
		var aliceOrder, bobOrder match.Match
		var aliceSendValue, bobSendValue *big.Int
		var aliceReceiveValue, bobReceiveValue *big.Int
		var aliceCurrency, bobCurrency uint32
		var alice, bob *ecdsa.PrivateKey
		var aliceEthKey, bobEthKey, aliceBtcKey, bobBtcKey swap.Key
		var swapID [32]byte

		rand.Read(swapID[:])

		aliceCurrency = 1
		bobCurrency = 0

		configAlice, err := config.LoadConfig(os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/secrets/local/configA.json")
		Expect(err).ShouldNot(HaveOccurred())

		configBob, err := config.LoadConfig(os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/secrets/local/configB.json")
		Expect(err).ShouldNot(HaveOccurred())

		ganache, err := ethclient.Connect(configAlice)
		Expect(err).ShouldNot(HaveOccurred())

		var ownPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/secrets/owner.json"

		own, err := owner.LoadOwner(ownPath)
		Expect(err).ShouldNot(HaveOccurred())

		pk, err := crypto.HexToECDSA(own.Ganache)
		Expect(err).ShouldNot(HaveOccurred())

		owner := bind.NewKeyedTransactor(pk)

		alice, err = crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		aliceAuth := bind.NewKeyedTransactor(alice)
		aliceEthKey, err = ethKey.NewEthereumKey(hex.EncodeToString(crypto.FromECDSA(alice)), "ganache")
		Expect(err).ShouldNot(HaveOccurred())
		aliceBtcKey, err = btcKey.NewBitcoinKey(hex.EncodeToString(crypto.FromECDSA(alice)), "regtest")
		Expect(err).ShouldNot(HaveOccurred())

		bob, err = crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		bobAuth := bind.NewKeyedTransactor(bob)
		bobEthKey, err = ethKey.NewEthereumKey(hex.EncodeToString(crypto.FromECDSA(bob)), "ganache")
		Expect(err).ShouldNot(HaveOccurred())
		bobBtcKey, err = btcKey.NewBitcoinKey(hex.EncodeToString(crypto.FromECDSA(bob)), "regtest")
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddrBytes, err := aliceEthKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		bobAddrBytes, err := bobEthKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())

		err = ganache.Transfer(common.BytesToAddress(aliceAddrBytes), owner, 1000000000000000000)
		Expect(err).ShouldNot(HaveOccurred())

		err = ganache.Transfer(common.BytesToAddress(bobAddrBytes), owner, 1000000000000000000)
		Expect(err).ShouldNot(HaveOccurred())

		time.Sleep(5 * time.Second)
		connection, err := btcclient.Connect(configBob)
		Expect(err).ShouldNot(HaveOccurred())

		aliceSendValue = big.NewInt(10000000)
		bobSendValue = big.NewInt(10000000)

		aliceReceiveValue = big.NewInt(99990000)
		bobReceiveValue = big.NewInt(8000000)

		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		time.Sleep(5 * time.Second)

		_AliceWIF, err := aliceBtcKey.GetKeyString()
		Expect(err).ShouldNot(HaveOccurred())

		AliceWIF, err := btcutil.DecodeWIF(_AliceWIF)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKeyLabel(AliceWIF, "alice")
		Expect(err).ShouldNot(HaveOccurred())

		_BobWIF, err := bobBtcKey.GetKeyString()
		Expect(err).ShouldNot(HaveOccurred())

		BobWIF, err := btcutil.DecodeWIF(_BobWIF)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKeyLabel(BobWIF, "bob")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = regtest.GetAddressForAccount(connection, "bob")
		Expect(err).ShouldNot(HaveOccurred())

		aliceNet, err = net.NewEthereumNetwork(ganache, aliceAuth)
		Expect(err).Should(BeNil())

		bobNet, err = net.NewEthereumNetwork(ganache, bobAuth)
		Expect(err).Should(BeNil())

		aliceInfo, err = ax.NewEthereumAtomInfo(ganache, aliceAuth)
		Expect(err).Should(BeNil())

		bobInfo, err = ax.NewEthereumAtomInfo(ganache, bobAuth)
		Expect(err).Should(BeNil())

		aliceOrder = match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, aliceReceiveValue, aliceCurrency, bobCurrency)
		bobOrder = match.NewMatch(bobOrderID, aliceOrderID, bobSendValue, bobReceiveValue, bobCurrency, aliceCurrency)

		aliceBtcAddrBytes, err := aliceBtcKey.GetAddress()
		Expect(err).Should(BeNil())

		bobEthAddrBytes, err := bobEthKey.GetAddress()
		Expect(err).Should(BeNil())

		aliceInfo.SetOwnerAddress(aliceOrderID, aliceBtcAddrBytes)
		bobInfo.SetOwnerAddress(bobOrderID, bobEthAddrBytes)

		aliceAtomBuilder := atoms.NewAtomBuilder(configAlice, []swap.Key{aliceEthKey, aliceBtcKey})
		bobAtomBuilder := atoms.NewAtomBuilder(configBob, []swap.Key{bobBtcKey, bobEthKey})

		aliceLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/atom-go/dbAlice")
		Expect(err).Should(BeNil())

		bobLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/atom-go/dbBob")
		Expect(err).Should(BeNil())

		aliceState := store.NewSwapState(aliceLDB)
		bobState := store.NewSwapState(bobLDB)

		mockWallet := wal.NewMockWallet()

		mockWallet.SetMatch(aliceOrderID, aliceOrder)
		mockWallet.SetMatch(bobOrderID, bobOrder)

		aliceWatch = NewWatch(aliceNet, aliceInfo, mockWallet, aliceAtomBuilder, aliceState)
		bobWatch = NewWatch(bobNet, bobInfo, mockWallet, bobAtomBuilder, bobState)
	})

	It("can do an eth - btc atomic swap (eth implementations)", func() {

		wg := &sync.WaitGroup{}
		doneChAlice := make(chan struct{}, 1)
		doneChBob := make(chan struct{}, 1)

		errChAlice := aliceWatch.Run(doneChAlice)
		errChBob := bobWatch.Run(doneChBob)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			for {
				select {
				case err, ok := <-errChAlice:
					if !ok {
						return
					}
					Expect(err).ShouldNot(HaveOccurred())
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			for {
				select {
				case err, ok := <-errChBob:
					if !ok {
						return
					}
					Expect(err).ShouldNot(HaveOccurred())
				}
			}
		}()

		Expect(aliceWatch.Add(aliceOrderID)).ShouldNot(HaveOccurred())
		Expect(bobWatch.Add(bobOrderID)).ShouldNot(HaveOccurred())

		aliceWatch.Notify()
		bobWatch.Notify()

		go func() {
			defer func() { doneChAlice <- struct{}{} }()
			defer func() { doneChBob <- struct{}{} }()

			for {
				if aliceWatch.Status(aliceOrderID) == swap.StatusRedeemed && bobWatch.Status(bobOrderID) == swap.StatusRedeemed {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}()

		wg.Wait()
	})

})
