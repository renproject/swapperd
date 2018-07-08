package btc_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"os"
	"time"

	"github.com/btcsuite/btcutil"

	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/atoms/btc"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/key/btc"
	"github.com/republicprotocol/atom-go/services/swap"
)

var _ = Describe("bitcoin", func() {

	var connection btcclient.Conn
	// var cmd *exec.Cmd
	var aliceAddr, bobAddr string // btcutil.Address
	var aliceAddrBytes, bobAddrBytes []byte
	var _aliceAddr, _bobAddr btcutil.Address

	var value *big.Int
	var validity int64
	var secret, secretHash [32]byte
	var err error
	var reqAtom, reqAtomFailed swap.Atom
	var resAtom swap.Atom
	var data []byte
	//	var confLocal = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/configLocal.json" // Bitcoin Regtest
	var confPath = os.Getenv("HOME") + "/go/src/github.com/republicprotocol/atom-go/secrets/configTestnet.json" // Bitcoin Testnet

	BeforeSuite(func() {
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		connection, err = btcclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		// go func() {
		// 	err = regtest.Mine(connection)
		// 	Expect(err).ShouldNot(HaveOccurred())
		// }()
		// time.Sleep(5 * time.Second)

		alice, err := crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		aliceKey, err := btc.NewBitcoinKey(hex.EncodeToString(crypto.FromECDSA(alice)), "testnet")
		Expect(err).ShouldNot(HaveOccurred())

		bob, err := crypto.GenerateKey()
		Expect(err).ShouldNot(HaveOccurred())
		bobKey, err := btc.NewBitcoinKey(hex.EncodeToString(crypto.FromECDSA(bob)), "testnet")
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddrBytes, err = aliceKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		bobAddrBytes, err = bobKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddr = string(aliceAddrBytes)
		bobAddr = string(bobAddrBytes)

		_aliceAddr, err = btcutil.DecodeAddress(aliceAddr, connection.ChainParams)
		Expect(err).ShouldNot(HaveOccurred())
		_bobAddr, err = btcutil.DecodeAddress(bobAddr, connection.ChainParams)
		Expect(err).ShouldNot(HaveOccurred())

		btcvalue, err := btcutil.NewAmount(0.05)
		Expect(err).ShouldNot(HaveOccurred())

		connection.Client.SendToAddress(_aliceAddr, btcvalue)
		connection.Client.SendToAddress(_bobAddr, btcvalue)

		_aliceWIF, err := aliceKey.GetKeyString()
		Expect(err).ShouldNot(HaveOccurred())

		aliceWIF, err := btcutil.DecodeWIF(_aliceWIF)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKey(aliceWIF)
		Expect(err).ShouldNot(HaveOccurred())

		_bobWIF, err := bobKey.GetKeyString()
		Expect(err).ShouldNot(HaveOccurred())

		bobWIF, err := btcutil.DecodeWIF(_bobWIF)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKey(bobWIF)
		Expect(err).ShouldNot(HaveOccurred())

		reqAtom = NewBitcoinAtom(connection, aliceKey)
		reqAtomFailed = NewBitcoinAtom(connection, aliceKey)
		resAtom = NewBitcoinAtom(connection, bobKey)

		value = big.NewInt(1000000)
		validity = time.Now().Unix() + 48*60*60
	})

	It("can initiate a btc atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtom.Initiate([]byte(bobAddr), secretHash, value, validity)
		Expect(err).ShouldNot(HaveOccurred())
		data, err = reqAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit a btc atomic swap", func() {
		err = resAtom.Deserialize(data)
		Expect(err).ShouldNot(HaveOccurred())
		err = resAtom.Audit(secretHash, []byte(bobAddr), value, 60*60)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can redeem a btc atomic swap", func() {
		// before, err := connection.Client.GetReceivedByAddress(_bobAddr)
		// Expect(err).ShouldNot(HaveOccurred())
		err = resAtom.Redeem(secret)
		Expect(err).ShouldNot(HaveOccurred())
		// after, err := connection.Client.GetReceivedByAddress(_bobAddr)
		// Expect(err).ShouldNot(HaveOccurred())
		data, err = resAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
		// Expect(after - before).Should(Equal(btcutil.Amount(990000)))
	})

	It("can audit secret after a btc atomic swap", func() {
		err = reqAtom.Deserialize(data)
		Expect(err).ShouldNot(HaveOccurred())
		_secret, err := reqAtom.AuditSecret()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(_secret).Should(Equal(secret))
	})

	It("can refund a btc atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtomFailed.Initiate([]byte(aliceAddr), secretHash, value, 0)
		Expect(err).ShouldNot(HaveOccurred())
		// before, err := connection.Client.GetReceivedByAddress(_aliceAddr)
		// Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Refund()
		Expect(err).ShouldNot(HaveOccurred())
		// after, err := connection.Client.GetReceivedByAddress(_aliceAddr)
		// Expect(err).ShouldNot(HaveOccurred())
		// Expect(after - before).Should(Equal(btcutil.Amount(990000)))
	})
})
