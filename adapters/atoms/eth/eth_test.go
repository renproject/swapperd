package eth_test

import (
	"context"
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/atoms/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/drivers/eth/ganache"
	"github.com/republicprotocol/atom-go/services/atom"
)

var _ = Describe("ether", func() {

	var alice, bob *bind.TransactOpts
	var conn ethclient.Connection
	var swapID, swapIDFailed [32]byte
	var aliceOrderID, bobOrderID [32]byte
	var value *big.Int
	var validity int64
	var secret, secretHash [32]byte
	var err error
	var reqAtom, reqAtomFailed atom.RequestAtom
	var resAtom atom.ResponseAtom
	var data []byte
	var bobAddr common.Address

	BeforeSuite(func() {
		// Setup...
		conn, err = ganache.Connect("http://localhost:8545")
		Expect(err).ShouldNot(HaveOccurred())

		alice, _, err = ganache.NewAccount(conn, big.NewInt(1000000000000000000))
		Expect(err).ShouldNot(HaveOccurred())
		alice.GasLimit = 3000000
		bob, bobAddr, err = ganache.NewAccount(conn, big.NewInt(1000000000000000000))
		Expect(err).ShouldNot(HaveOccurred())
		bob.GasLimit = 3000000

		value = big.NewInt(10)
		validity = int64(time.Hour * 24)

		swapID[0] = 0x13
		swapIDFailed[0] = 0x23

		aliceOrderID[0] = 0x33
		bobOrderID[0] = 0x4a

		reqAtom, err = NewEthereumRequestAtom(context.Background(), conn, alice, bobAddr, swapID)
		Expect(err).ShouldNot(HaveOccurred())

		reqAtomFailed, err = NewEthereumRequestAtom(context.Background(), conn, alice, bobAddr, swapIDFailed)
		Expect(err).ShouldNot(HaveOccurred())

		resAtom, err = NewEthereumResponseAtom(context.Background(), conn, bob)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can initiate an eth atomic swap", func() {
		secret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(secret[:])
		err = reqAtom.Initiate(secretHash, value, validity)
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
		err = reqAtomFailed.Initiate(secretHash, value, 0)
		Expect(err).ShouldNot(HaveOccurred())
		err = reqAtomFailed.Refund()
		Expect(err).ShouldNot(HaveOccurred())
	})
})
