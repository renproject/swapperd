package watch_test

import (
	"crypto/rand"
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

	"github.com/republicprotocol/renex-swapper-go/adapter/atoms"
	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/eth"
	binder "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/binder"
	btcclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/btc"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/configs/network"
	"github.com/republicprotocol/renex-swapper-go/adapter/configs/owner"
	"github.com/republicprotocol/renex-swapper-go/adapter/store/leveldb"
	"github.com/republicprotocol/renex-swapper-go/domains/match"
	"github.com/republicprotocol/renex-swapper-go/drivers/btc/regtest"
	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/republicprotocol/renex-swapper-go/services/swap"
	. "github.com/republicprotocol/renex-swapper-go/services/watch"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap using Watch", func() {

	var aliceWatch, bobWatch Watch
	var aliceOrderID, bobOrderID [32]byte

	rand.Read(aliceOrderID[:])
	rand.Read(bobOrderID[:])

	BeforeSuite(func() {
		netConf, aliceKS, bobKS := LoadConfigs()
		ethConn, aliceBinder, bobBinder := SetupEthereumNetwork(netConf, aliceKS, bobKS)
		btcConn := SetupBitcoinNetwork(netConf, aliceKS, bobKS)
		aliceMatch, bobMatch := GetMatches()
		SendAddresses(aliceMatch.PersonalOrderID(), bobMatch.PersonalOrderID(), aliceKS, bobKS, aliceBinder, bobBinder)
		aliceWatch, bobWatch := buildWatchers()

	})

	It("can do an eth - btc atomic swap (eth implementations)", func() {

		wg := &sync.WaitGroup{}

		errChAlice := aliceWatch.Start()
		errChBob := bobWatch.Start()

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
			defer aliceWatch.Stop()
			defer bobWatch.Stop()

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

func LoadConfigs() (network.Config, keystore.Keystore, keystore.Keystore) {
	var confPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/renex-swapper-go/secrets/local/networkA.json"
	config, err := network.LoadNetwork(confPath)
	Expect(err).ShouldNot(HaveOccurred())

	var ksPathA = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/renex-swapper-go/secrets/local/keystoreA.json"
	ksA, err := keystore.LoadKeystore(ksPathA)
	Expect(err).ShouldNot(HaveOccurred())

	var ksPathB = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/renex-swapper-go/secrets/local/keystoreB.json"
	ksB, err := keystore.LoadKeystore(ksPathB)
	Expect(err).ShouldNot(HaveOccurred())

	return config, ksA, ksB
}

func SetupEthereumNetwork(netConf network.Config, ksA keystore.Keystore, ksB keystore.Keystore) (ethclient.Conn, binder.Binder, binder.Binder) {
	ganache, err := ethclient.Connect(netConf)
	Expect(err).ShouldNot(HaveOccurred())

	aliceKey := ksA.EthereumKey
	aliceBinder, err := binder.NewBinder(aliceKey.GetKey(), ganache)

	bobKey := ksB.EthereumKey
	bobBinder, err := binder.NewBinder(bobKey.GetKey(), ganache)

	aliceAddrBytes, err := aliceKey.GetAddress()
	Expect(err).ShouldNot(HaveOccurred())
	bobAddrBytes, err := bobKey.GetAddress()
	Expect(err).ShouldNot(HaveOccurred())

	var ownPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/renex-swapper-go/secrets/owner.json"

	own, err := owner.LoadOwner(ownPath)
	Expect(err).ShouldNot(HaveOccurred())

	pk, err := crypto.HexToECDSA(own.Ganache)
	Expect(err).ShouldNot(HaveOccurred())

	owner := bind.NewKeyedTransactor(pk)

	err = ganache.Transfer(common.BytesToAddress(aliceAddrBytes), owner, 1000000000000000000)
	Expect(err).ShouldNot(HaveOccurred())

	err = ganache.Transfer(common.BytesToAddress(bobAddrBytes), owner, 1000000000000000000)
	Expect(err).ShouldNot(HaveOccurred())

	return ganache, aliceBinder, bobBinder
}

func SetupBitcoinNetwork(netConf network.Config, ksA, ksB keystore.Keystore) btcclient.Conn {
	time.Sleep(5 * time.Second)
	connection, err := btcclient.Connect(netConf)
	Expect(err).ShouldNot(HaveOccurred())

	go func() {
		err = regtest.Mine(connection)
		Expect(err).ShouldNot(HaveOccurred())
	}()
	time.Sleep(5 * time.Second)

	_AliceWIF, err := ksA.BitcoinKey.GetKeyString()
	Expect(err).ShouldNot(HaveOccurred())

	AliceWIF, err := btcutil.DecodeWIF(_AliceWIF)
	Expect(err).ShouldNot(HaveOccurred())

	err = connection.Client.ImportPrivKeyLabel(AliceWIF, "alice")
	Expect(err).ShouldNot(HaveOccurred())

	_BobWIF, err := ksB.BitcoinKey.GetKeyString()
	Expect(err).ShouldNot(HaveOccurred())

	BobWIF, err := btcutil.DecodeWIF(_BobWIF)
	Expect(err).ShouldNot(HaveOccurred())

	err = connection.Client.ImportPrivKeyLabel(BobWIF, "bob")
	Expect(err).ShouldNot(HaveOccurred())

	_, err = regtest.GetAddressForAccount(connection, "bob")
	Expect(err).ShouldNot(HaveOccurred())

	return connection
}

func GetMatches() (match.Match, match.Match) {
	rand.Read(aliceOrderID[:])
	rand.Read(bobOrderID[:])

	aliceCurrency := uint32(1)
	bobCurrency := uint32(0)

	aliceSendValue := big.NewInt(10000000)
	bobSendValue := big.NewInt(10000000)

	aliceReceiveValue := big.NewInt(99990000)
	bobReceiveValue := big.NewInt(8000000)

	aliceOrder := match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, aliceReceiveValue, aliceCurrency, bobCurrency)
	bobOrder := match.NewMatch(bobOrderID, aliceOrderID, bobSendValue, bobReceiveValue, bobCurrency, aliceCurrency)

	return aliceOrder, bobOrder
}

func SendAddresses(aliceOrderID, bobOrderID [32]byte, aliceKS, bobKS keystore.Keystore, aliceBinder, bobBinder binder.Binder) {

	err := aliceBinder.SubmitBuyOrder(aliceOrderID)
	Expect(err).Should(BeNil())
	err = bobBinder.SubmitSellOrder(bobOrderID)
	Expect(err).Should(BeNil())
	err = aliceBinder.AuthorizeAtomBox()
	Expect(err).Should(BeNil())
	err = bobBinder.AuthorizeAtomBox()
	Expect(err).Should(BeNil())

	aliceBtcAddrBytes, err := aliceKS.BitcoinKey.GetAddress()
	Expect(err).Should(BeNil())

	bobEthAddrBytes, err := bobKS.EthereumKey.GetAddress()
	Expect(err).Should(BeNil())

	err = aliceBinder.SendOwnerAddress(aliceOrderID, aliceBtcAddrBytes)
	Expect(err).Should(BeNil())
	err = bobBinder.SendOwnerAddress(bobOrderID, bobEthAddrBytes)
	Expect(err).Should(BeNil())
}

func SetupWatchers(ethConn ethclient.Conn, btcConn btcclient.Conn, aliceMatch, bobMatch match.Match, aliceKS, bobKS keystore.Keystore, aliceBinder, bobBinder binder.Binder) (Swap, Swap) {
	aliceAtomBuilder, err := atoms.NewAtomBuilder(aliceBinder, configAlice, aliceKS)
	Expect(err).Should(BeNil())
	bobAtomBuilder, err := atoms.NewAtomBuilder(bobBinder, configBob, bobKS)
	Expect(err).Should(BeNil())

	aliceLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/renex-swapper-go/dbAlice")
	Expect(err).Should(BeNil())

	bobLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/renex-swapper-go/dbBob")
	Expect(err).Should(BeNil())

	aliceState := store.NewState(aliceLDB)
	bobState := store.NewState(bobLDB)

	mockWallet := wal.NewMockWallet()

	mockWallet.SetMatch(aliceOrderID, aliceOrder)
	mockWallet.SetMatch(bobOrderID, bobOrder)

	aliceWatch = NewWatch(aliceNet, aliceInfo, mockWallet, aliceAtomBuilder, aliceState)
	bobWatch = NewWatch(bobNet, bobInfo, mockWallet, bobAtomBuilder, bobState)

	reqAlice, err := eth.NewEthereumAtom(ethConn, &aliceKS.EthereumKey, aliceMatch.PersonalOrderID())
	Expect(err).Should(BeNil())

	reqBob := btc.NewBitcoinAtom(btcConn, &bobKS.BitcoinKey, bobMatch.PersonalOrderID())
	resAlice := btc.NewBitcoinAtom(btcConn, &aliceKS.BitcoinKey, aliceMatch.ForeignOrderID())

	resBob, err := eth.NewEthereumAtom(ethConn, &bobKS.EthereumKey, bobMatch.ForeignOrderID())
	Expect(err).Should(BeNil())

	aliceLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/renex-swapper-go/temp/dbAlice")
	Expect(err).Should(BeNil())

	bobLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/renex-swapper-go/temp/dbBob")
	Expect(err).Should(BeNil())

	aliceState := store.NewState(aliceLDB)
	bobState := store.NewState(bobLDB)

	aliceState.PutStatus(aliceMatch.PersonalOrderID(), StatusInfoSubmitted)
	bobState.PutStatus(bobMatch.PersonalOrderID(), StatusInfoSubmitted)

	aliceSwap := NewSwap(reqAlice, resAlice, aliceMatch, &aliceBinder, aliceState)
	bobSwap := NewSwap(reqBob, resBob, bobMatch, &bobBinder, bobState)

	return aliceSwap, bobSwap
}
