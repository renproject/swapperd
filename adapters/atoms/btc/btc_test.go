package btc_test

import (
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/btcsuite/btcutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/atoms/btc"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	"github.com/republicprotocol/atom-go/services/swap"
)

var _ = Describe("bitcoin", func() {

	var connection btcclient.Conn
	// var cmd *exec.Cmd
	var aliceAddr, bobAddr string // btcutil.Address
	var _aliceAddr, _bobAddr btcutil.Address

	var value *big.Int
	var validity int64
	var secret, secretHash [32]byte
	var err error
	var reqAtom, reqAtomFailed swap.AtomRequester
	var resAtom swap.AtomResponder
	var data []byte
	var confPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/config.json"

	BeforeSuite(func() {
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		connection, err = btcclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		time.Sleep(5 * time.Second)

		_aliceAddr, err = connection.Client.GetAccountAddress("alice")
		Expect(err).ShouldNot(HaveOccurred())
		aliceAddr = _aliceAddr.EncodeAddress()

		_bobAddr, err = connection.Client.GetAccountAddress("bob")
		Expect(err).ShouldNot(HaveOccurred())
		bobAddr = _bobAddr.EncodeAddress()

		reqAtom = NewBitcoinAtomRequester(connection, aliceAddr)
		reqAtomFailed = NewBitcoinAtomRequester(connection, aliceAddr)
		resAtom = NewBitcoinAtomResponder(connection, bobAddr)

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
		before, err := connection.Client.GetReceivedByAddress(_bobAddr)
		Expect(err).ShouldNot(HaveOccurred())
		err = resAtom.Redeem(secret)
		Expect(err).ShouldNot(HaveOccurred())
		after, err := connection.Client.GetReceivedByAddress(_bobAddr)
		Expect(err).ShouldNot(HaveOccurred())
		data, err = resAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(after - before).Should(Equal(btcutil.Amount(990000)))
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
		before, err := connection.Client.GetReceivedByAddress(_aliceAddr)
		Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Refund()
		Expect(err).ShouldNot(HaveOccurred())
		after, err := connection.Client.GetReceivedByAddress(_aliceAddr)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(after - before).Should(Equal(btcutil.Amount(990000)))
	})
})
