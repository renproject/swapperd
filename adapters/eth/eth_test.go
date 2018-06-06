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
	. "github.com/republicprotocol/atom-go/adapters/eth"
	"github.com/republicprotocol/atom-go/drivers/eth/ganache"
)

var _ = Describe("ether", func() {

	var alice, bob *bind.TransactOpts
	var aliceAddr, bobAddr common.Address
	var conn Connection
	var aliceSwapID [32]byte
	var bobSwapID [32]byte
	var aliceOrderID [32]byte
	var bobOrderID [32]byte
	var value *big.Int
	var validity int64

	BeforeEach(func() {
		// Setup...
		var err error
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
	})

	It("can perform ETH-ETH atom swap", func() {

		var aliceSecret [32]byte
		var secretHash [32]byte

		{ // Alice can initiate swap
			aliceAtom, err := NewEthereumAtom(context.Background(), conn, alice, aliceSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			aliceSecret = [32]byte{1, 3, 3, 7}
			secretHash = sha256.Sum256(aliceSecret[:])
			err = aliceAtom.Initiate(secretHash, aliceAddr.Bytes(), bobAddr.Bytes(), value, validity)
			Expect(err).ShouldNot(HaveOccurred())
			err = aliceAtom.Store(aliceOrderID)
			Expect(err).ShouldNot(HaveOccurred())
		}

		{ // Bob can audit Alice's contract and upload his own
			aliceAtom, err := NewEthereumAtom(context.Background(), conn, bob, [32]byte{})
			Expect(err).ShouldNot(HaveOccurred())
			err = aliceAtom.Retrieve(aliceOrderID)
			Expect(err).ShouldNot(HaveOccurred())

			_secretHash, _from, _to, _value, _expiry, err := aliceAtom.Audit()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(_secretHash).Should(Equal(secretHash))
			Expect(_from).Should(Equal(aliceAddr.Bytes()))
			Expect(_to).Should(Equal(bobAddr.Bytes()))
			Expect(_value).Should(Equal(value))
			Expect(_expiry).Should(Equal(validity))

			bobAtom, err := NewEthereumAtom(context.Background(), conn, bob, bobSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobAtom.Initiate(_secretHash, bobAddr.Bytes(), aliceAddr.Bytes(), value, validity)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobAtom.Store(bobOrderID)
			Expect(err).ShouldNot(HaveOccurred())

		}

		{ // Alice can audit Bob's contract and reveal the secret
			bobAtom, err := NewEthereumAtom(context.Background(), conn, alice, bobSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobAtom.Retrieve(bobOrderID)
			Expect(err).ShouldNot(HaveOccurred())

			_secretHash, _from, _to, _value, _expiry, err := bobAtom.Audit()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(_secretHash).Should(Equal(secretHash))
			Expect(_from).Should(Equal(bobAddr.Bytes()))
			Expect(_to).Should(Equal(aliceAddr.Bytes()))
			Expect(_value).Should(Equal(value))
			Expect(_expiry).Should(Equal(validity))

			err = bobAtom.Redeem(aliceSecret)
			Expect(err).ShouldNot(HaveOccurred())
		}

		{ // Bob can retrieve the secret from his contract and complete the swap
			bobAtom, err := NewEthereumAtom(context.Background(), conn, bob, bobSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobAtom.Retrieve(bobOrderID)
			Expect(err).ShouldNot(HaveOccurred())

			secret, err := bobAtom.AuditSecret()
			Expect(err).ShouldNot(HaveOccurred())

			aliceAtom, err := NewEthereumAtom(context.Background(), conn, bob, aliceSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = aliceAtom.Retrieve(aliceOrderID)
			Expect(err).ShouldNot(HaveOccurred())

			err = aliceAtom.Redeem(secret)
			Expect(err).ShouldNot(HaveOccurred())
		}
	})
})
