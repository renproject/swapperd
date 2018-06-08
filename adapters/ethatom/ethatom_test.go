package ethatom_test

import (
	"context"
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/atom-go/adapters/ethatom"
	"github.com/republicprotocol/atom-go/adapters/ethclient"
	"github.com/republicprotocol/atom-go/drivers/eth/ganache"
	"github.com/republicprotocol/atom-go/services/atom"
)

var _ = Describe("ether", func() {

	var alice, bob *bind.TransactOpts
	var aliceAddr, bobAddr common.Address
	var conn ethclient.Connection
	var aliceSwapID [32]byte
	var bobSwapID [32]byte
	var aliceOrderID [32]byte
	var bobOrderID [32]byte
	var value *big.Int
	var validity int64
	var aliceSecret [32]byte
	var secretHash [32]byte
	var aliceAtom, bobAliceAtom atom.Atom
	var bobAtom, aliceBobAtom atom.Atom
	var aliceData, bobData []byte
	var bobAliceData, aliceBobData []byte
	var err error

	BeforeSuite(func() {
		// Setup...
		conn, err = ganache.Connect("http://localhost:8545")
		Expect(err).ShouldNot(HaveOccurred())

		alice, aliceAddr, err = ganache.NewAccount(conn, big.NewInt(1000000000000000000))
		Expect(err).ShouldNot(HaveOccurred())
		alice.GasLimit = 3000000
		bob, bobAddr, err = ganache.NewAccount(conn, big.NewInt(1000000000000000000))
		Expect(err).ShouldNot(HaveOccurred())
		bob.GasLimit = 3000000

		value = big.NewInt(10)
		validity = int64(time.Hour * 24)

		aliceSwapID[0] = 0x13
		bobSwapID[0] = 0x1a

		aliceOrderID[0] = 0x33
		bobOrderID[0] = 0x4a

		bobAliceAtom, err = NewEthereumAtom(context.Background(), conn, bob, [32]byte{})
		bobAliceData, err = bobAliceAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())

		aliceBobAtom, err = NewEthereumAtom(context.Background(), conn, alice, [32]byte{})
		aliceBobData, err = aliceBobAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can initiate an eth atomic swap", func() {
		aliceAtom, err = NewEthereumAtom(context.Background(), conn, alice, aliceSwapID)
		Expect(err).ShouldNot(HaveOccurred())
		err = aliceAtom.Deserialize(bobAliceData)
		Expect(err).ShouldNot(HaveOccurred())
		aliceSecret = [32]byte{1, 3, 3, 7}
		secretHash = sha256.Sum256(aliceSecret[:])
		err = aliceAtom.Initiate(secretHash, value, validity)
		Expect(err).ShouldNot(HaveOccurred())
		aliceData, err = aliceAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit and initiate an atomic swap", func() {
		err = bobAliceAtom.Deserialize(aliceData)
		Expect(err).ShouldNot(HaveOccurred())
		_secretHash, _from, _to, _value, _expiry, err := bobAliceAtom.Audit()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(_secretHash).Should(Equal(secretHash))
		Expect(_to).Should(Equal(bobAddr.Bytes()))
		Expect(_from).Should(Equal(aliceAddr.Bytes()))
		Expect(_value).Should(Equal(value))
		Expect(_expiry).Should(Equal(validity))

		bobAtom, err = NewEthereumAtom(context.Background(), conn, bob, bobSwapID)
		Expect(err).ShouldNot(HaveOccurred())
		err = bobAtom.Deserialize(aliceBobData)
		Expect(err).ShouldNot(HaveOccurred())
		err = bobAtom.Initiate(_secretHash, value, validity)
		Expect(err).ShouldNot(HaveOccurred())
		bobData, err = bobAtom.Serialize()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can audit atom details and reveal the secret", func() {
		err = aliceBobAtom.Deserialize(bobData)
		Expect(err).ShouldNot(HaveOccurred())
		_secretHash, _from, _to, _value, _expiry, err := aliceBobAtom.Audit()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(_secretHash).Should(Equal(secretHash))
		Expect(_from).Should(Equal(bobAddr.Bytes()))
		Expect(_to).Should(Equal(aliceAddr.Bytes()))
		Expect(_value).Should(Equal(value))
		Expect(_expiry).Should(Equal(validity))

		err = aliceBobAtom.Redeem(aliceSecret)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can retrieve the secret from his contract and complete the swap", func() {
		secret, err := bobAtom.AuditSecret()
		Expect(err).ShouldNot(HaveOccurred())

		err = bobAliceAtom.Redeem(secret)
		Expect(err).ShouldNot(HaveOccurred())
	})
})
