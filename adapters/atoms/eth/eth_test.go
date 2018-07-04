package eth_test

import (
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/atoms/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/services/swap"
)

var _ = Describe("ether", func() {

	var bobAddr, aliceAddr common.Address
	var alice, bob *bind.TransactOpts
	var conn ethclient.Conn
	var aliceOrderID, bobOrderID [32]byte
	var value *big.Int
	var validity int64
	var secret, secretHash [32]byte
	var err error
	var reqAtom, reqAtomFailed swap.Atom
	var resAtom swap.Atom
	var data []byte
	var confPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/config.json"
	var ksPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/keystore.json"

	BeforeSuite(func() {
		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		keystore := keystore.NewKeystore(ksPath)
		conn, err = ethclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		pk, err := keystore.LoadKeypair("ethereum")
		Expect(err).ShouldNot(HaveOccurred())
		owner := bind.NewKeyedTransactor(pk)
		aliceAddr, alice, err = conn.NewAccount(1000000000000000000, owner)
		alice.GasLimit = 3000000
		bobAddr, bob, err = conn.NewAccount(1000000000000000000, owner)
		bob.GasLimit = 3000000
		Expect(err).ShouldNot(HaveOccurred())
		value = big.NewInt(10)
		validity = int64(time.Hour * 24)
		aliceOrderID[0] = 0x33
		bobOrderID[0] = 0x4a
		reqAtom, err = NewEthereumAtom(conn, alice)
		Expect(err).ShouldNot(HaveOccurred())
		reqAtomFailed, err = NewEthereumAtom(conn, alice)
		Expect(err).ShouldNot(HaveOccurred())
		resAtom, err = NewEthereumAtom(conn, bob)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can initiate an eth atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtom.Initiate(bobAddr.Bytes(), secretHash, value, validity)
		Expect(err).ShouldNot(HaveOccurred())
		data, err = reqAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit an eth atomic swap", func() {
		err = resAtom.Deserialize(data)
		Expect(err).ShouldNot(HaveOccurred())
		err = resAtom.Audit(secretHash, bobAddr.Bytes(), value, 60*60)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can redeem an eth atomic swap", func() {
		err = resAtom.Redeem(secret)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit secret after an eth atomic swap", func() {
		_secret, err := reqAtom.AuditSecret()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(_secret).Should(Equal(secret))
	})

	It("can refund an eth atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtomFailed.Initiate(bobAddr.Bytes(), secretHash, value, 0)
		Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Refund()
		Expect(err).ShouldNot(HaveOccurred())
	})
})
