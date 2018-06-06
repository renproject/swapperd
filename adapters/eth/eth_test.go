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
	})

	It("can perform ETH-ETH arc swap", func() {

		var aliceArcData, bobArcData []byte

		var aliceSecret [32]byte
		var secretHash [32]byte

		{ // Alice can initiate swap
			aliceArc, err := NewEthereumArc(context.Background(), conn, alice, aliceSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			aliceSecret = [32]byte{1, 3, 3, 7}
			secretHash = sha256.Sum256(aliceSecret[:])
			err = aliceArc.Initiate(secretHash, aliceAddr.Bytes(), bobAddr.Bytes(), value, validity)
			Expect(err).ShouldNot(HaveOccurred())
			aliceArcData, err = aliceArc.Serialize()
			Expect(err).ShouldNot(HaveOccurred())
		}

		{ // Bob can audit Alice's contract and upload his own
			aliceArc, err := NewEthereumArc(context.Background(), conn, bob, [32]byte{})
			Expect(err).ShouldNot(HaveOccurred())
			err = aliceArc.Deserialize(aliceArcData)
			Expect(err).ShouldNot(HaveOccurred())

			_secretHash, _from, _to, _value, _expiry, err := aliceArc.Audit()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(_secretHash).Should(Equal(secretHash))
			Expect(_from).Should(Equal(aliceAddr.Bytes()))
			Expect(_to).Should(Equal(bobAddr.Bytes()))
			Expect(_value).Should(Equal(value))
			Expect(_expiry).Should(Equal(validity))

			bobArc, err := NewEthereumArc(context.Background(), conn, bob, bobSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobArc.Initiate(_secretHash, bobAddr.Bytes(), aliceAddr.Bytes(), value, validity)
			Expect(err).ShouldNot(HaveOccurred())
			bobArcData, err = bobArc.Serialize()
			Expect(err).ShouldNot(HaveOccurred())

		}

		{ // Alice can audit Bob's contract and reveal the secret
			bobArc, err := NewEthereumArc(context.Background(), conn, alice, bobSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobArc.Deserialize(bobArcData)
			Expect(err).ShouldNot(HaveOccurred())

			_secretHash, _from, _to, _value, _expiry, err := bobArc.Audit()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(_secretHash).Should(Equal(secretHash))
			Expect(_from).Should(Equal(bobAddr.Bytes()))
			Expect(_to).Should(Equal(aliceAddr.Bytes()))
			Expect(_value).Should(Equal(value))
			Expect(_expiry).Should(Equal(validity))

			err = bobArc.Redeem(aliceSecret)
			Expect(err).ShouldNot(HaveOccurred())
		}

		{ // Bob can retrieve the secret from his contract and complete the swap
			bobArc, err := NewEthereumArc(context.Background(), conn, bob, bobSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobArc.Deserialize(bobArcData)
			Expect(err).ShouldNot(HaveOccurred())

			secret, err := bobArc.AuditSecret()
			Expect(err).ShouldNot(HaveOccurred())

			aliceArc, err := NewEthereumArc(context.Background(), conn, bob, aliceSwapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = aliceArc.Deserialize(aliceArcData)
			Expect(err).ShouldNot(HaveOccurred())

			err = aliceArc.Redeem(secret)
			Expect(err).ShouldNot(HaveOccurred())
		}
	})
})
