package swap_test

import (
	"crypto/rand"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/atom-go/services/store"

	"github.com/republicprotocol/atom-go/domains/match"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	. "github.com/republicprotocol/atom-go/services/swap"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	"github.com/republicprotocol/atom-go/adapters/configs/keystore"
	"github.com/republicprotocol/atom-go/adapters/configs/owner"

	btcclient "github.com/republicprotocol/atom-go/adapters/blockchain/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/blockchain/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/configs/network"

	"github.com/republicprotocol/atom-go/adapters/blockchain/binder"
	"github.com/republicprotocol/atom-go/adapters/store/leveldb"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap", func() {

	var aliceSwap, bobSwap Swap

	BeforeSuite(func() {
		netConf, aliceKS, bobKS := LoadConfigs()
		ethConn, aliceBinder, bobBinder := SetupEthereumNetwork(netConf, aliceKS, bobKS)
		btcConn := SetupBitcoinNetwork(netConf, aliceKS, bobKS)
		aliceMatch, bobMatch := GetMatches()
		SendAddresses(aliceMatch.PersonalOrderID(), bobMatch.PersonalOrderID(), aliceKS, bobKS, aliceBinder, bobBinder)
		aliceSwap, bobSwap = SetupSwaps(ethConn, btcConn, aliceMatch, bobMatch, aliceKS, bobKS, aliceBinder, bobBinder)
	})

	It("can do an eth - btc atomic swap", func() {
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

func LoadConfigs() (network.Config, keystore.Keystore, keystore.Keystore) {
	var confPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/local_secrets/local/networkA.json"
	config, err := network.LoadNetwork(confPath)
	Expect(err).ShouldNot(HaveOccurred())

	var ksPathA = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/secrets/local.alice.json"
	ksA, err := keystore.Load(ksPathA)
	Expect(err).ShouldNot(HaveOccurred())

	var ksPathB = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/secrets/local.bob.json"
	ksB, err := keystore.Load(ksPathB)
	Expect(err).ShouldNot(HaveOccurred())

	return config, ksA, ksB
}

func SetupEthereumNetwork(netConf network.Config, ksA keystore.Keystore, ksB keystore.Keystore) (ethclient.Conn, binder.Binder, binder.Binder) {
	ganache, err := ethclient.Connect(netConf)
	Expect(err).ShouldNot(HaveOccurred())

	aliceKey, err := ksA.GetKey(1, 0)
	Expect(err).ShouldNot(HaveOccurred())
	aliceEthKey, err := aliceKey.GetKey()
	Expect(err).ShouldNot(HaveOccurred())
	aliceBinder, err := binder.NewBinder(aliceEthKey, ganache)
	Expect(err).ShouldNot(HaveOccurred())

	bobKey, err := ksB.GetKey(1, 0)
	Expect(err).ShouldNot(HaveOccurred())
	bobEthKey, err := bobKey.GetKey()
	Expect(err).ShouldNot(HaveOccurred())
	bobBinder, err := binder.NewBinder(bobEthKey, ganache)

	aliceAddrBytes, err := aliceKey.GetAddress()
	Expect(err).ShouldNot(HaveOccurred())
	bobAddrBytes, err := bobKey.GetAddress()
	Expect(err).ShouldNot(HaveOccurred())

	var ownPath = os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/local_secrets/owner.json"

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

	AliceBtcKey, err := ksA.GetKey(0, 0)
	Expect(err).ShouldNot(HaveOccurred())

	_AliceWIF := AliceBtcKey.GetKeyString()
	AliceWIF, err := btcutil.DecodeWIF(_AliceWIF)
	Expect(err).ShouldNot(HaveOccurred())

	err = connection.Client.ImportPrivKeyLabel(AliceWIF, "alice")
	Expect(err).ShouldNot(HaveOccurred())

	BobBtcKey, err := ksB.GetKey(0, 0)
	Expect(err).ShouldNot(HaveOccurred())
	_BobWIF := BobBtcKey.GetKeyString()

	BobWIF, err := btcutil.DecodeWIF(_BobWIF)
	Expect(err).ShouldNot(HaveOccurred())

	err = connection.Client.ImportPrivKeyLabel(BobWIF, "bob")
	Expect(err).ShouldNot(HaveOccurred())

	_, err = regtest.GetAddressForAccount(connection, "bob")
	Expect(err).ShouldNot(HaveOccurred())

	return connection
}

func GetMatches() (match.Match, match.Match) {
	var aliceOrderID, bobOrderID [32]byte
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

	AliceBtcKey, err := aliceKS.GetKey(0, 0)
	Expect(err).ShouldNot(HaveOccurred())

	BobEthKey, err := bobKS.GetKey(1, 0)
	Expect(err).ShouldNot(HaveOccurred())

	aliceBtcAddrBytes, err := AliceBtcKey.GetAddress()
	Expect(err).Should(BeNil())

	bobEthAddrBytes, err := BobEthKey.GetAddress()
	Expect(err).Should(BeNil())

	err = aliceBinder.SendOwnerAddress(aliceOrderID, aliceBtcAddrBytes)
	Expect(err).Should(BeNil())
	err = bobBinder.SendOwnerAddress(bobOrderID, bobEthAddrBytes)
	Expect(err).Should(BeNil())
}

func SetupSwaps(ethConn ethclient.Conn, btcConn btcclient.Conn, aliceMatch, bobMatch match.Match, aliceKS, bobKS keystore.Keystore, aliceBinder, bobBinder binder.Binder) (Swap, Swap) {
	aliceEthKey, err := aliceKS.GetKey(1, 0)
	Expect(err).Should(BeNil())
	aliceBtcKey, err := aliceKS.GetKey(0, 0)
	Expect(err).Should(BeNil())
	bobEthKey, err := bobKS.GetKey(1, 0)
	Expect(err).Should(BeNil())
	bobBtcKey, err := bobKS.GetKey(0, 0)
	Expect(err).Should(BeNil())

	reqAlice, err := eth.NewEthereumAtom(&aliceBinder, ethConn, aliceEthKey, aliceMatch.PersonalOrderID())
	Expect(err).Should(BeNil())

	reqBob := btc.NewBitcoinAtom(&bobBinder, btcConn, bobBtcKey, bobMatch.PersonalOrderID())
	resAlice := btc.NewBitcoinAtom(&aliceBinder, btcConn, aliceBtcKey, bobMatch.PersonalOrderID())

	resBob, err := eth.NewEthereumAtom(&bobBinder, ethConn, bobEthKey, aliceMatch.PersonalOrderID())
	Expect(err).Should(BeNil())

	aliceLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/atom-go/temp/dbAlice")
	Expect(err).Should(BeNil())

	bobLDB, err := leveldb.NewLDBStore("/Users/susruth/go/src/github.com/republicprotocol/atom-go/temp/dbBob")
	Expect(err).Should(BeNil())

	aliceState := store.NewState(aliceLDB)
	bobState := store.NewState(bobLDB)

	aliceState.PutStatus(aliceMatch.PersonalOrderID(), StatusInfoSubmitted)
	bobState.PutStatus(bobMatch.PersonalOrderID(), StatusInfoSubmitted)

	aliceSwap := NewSwap(reqAlice, resAlice, aliceMatch, &aliceBinder, aliceState)
	bobSwap := NewSwap(reqBob, resBob, bobMatch, &bobBinder, bobState)

	return aliceSwap, bobSwap
}
