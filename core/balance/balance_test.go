package balance_test

import (
	"math/rand"
	"os"
	"reflect"
	"testing/quick"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	. "github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/testutils"
)

var Random *rand.Rand
var Logger *logrus.Logger

func init() {
	Random = rand.New(rand.NewSource(time.Now().Unix()))
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
}

var _ = Describe("Token balance on blockchain", func() {

	init := func(frequency time.Duration) (chan struct{}, map[blockchain.TokenName]blockchain.Balance, *testutils.MockBlockchain, Balances) {
		balance := randomBalance()
		done := make(chan struct{})
		bc := testutils.NewMockBlockchain(balance)
		balances := New(frequency, bc, Logger)
		return done, balance, bc, balances
	}

	Context("when querying balance", func() {

		It("should return the buffered balance", func() {
			done, balance, _, balances := init(time.Second)
			defer close(done)
			queries := make(chan BalanceQuery)
			go balances.Run(done, queries)
			time.Sleep(10 * time.Millisecond)

			responder := make(chan map[blockchain.TokenName]blockchain.Balance, 1)
			query := BalanceQuery{
				Response: responder,
			}
			queries <- query
			queriedBalance := <-responder
			Expect(reflect.DeepEqual(queriedBalance, balance)).Should(BeTrue())
		})

		It("should return the cached balance", func() {
			done, balance, bc, balances := init(2 * time.Second)
			defer close(done)
			queries := make(chan BalanceQuery)
			go balances.Run(done, queries)
			time.Sleep(10 * time.Millisecond)

			for i := 0; i < 10; i++ {
				newBalance := randomBalance()
				bc.UpdateBalance(newBalance)
				time.Sleep(100 * time.Millisecond)
				responder := make(chan map[blockchain.TokenName]blockchain.Balance, 1)
				query := BalanceQuery{
					Response: responder,
				}
				queries <- query
				queriedBalance := <-responder
				Expect(reflect.DeepEqual(queriedBalance, balance)).Should(BeTrue())
			}
		})

		It("should update the balance in the background with given frequency", func() {
			done, balance, bc, balances := init(200 * time.Millisecond)
			defer close(done)
			queries := make(chan BalanceQuery)
			go balances.Run(done, queries)
			time.Sleep(10 * time.Millisecond)

			for i := 0; i < 10; i++ {
				newBalance := randomBalance()
				bc.UpdateBalance(newBalance)
				time.Sleep(100 * time.Millisecond)
				responder := make(chan map[blockchain.TokenName]blockchain.Balance)
				query := BalanceQuery{
					Response: responder,
				}
				queries <- query
				queriedBalance := <-responder
				if i%2 == 0 {
					Expect(reflect.DeepEqual(queriedBalance, balance)).Should(BeTrue())
				} else {
					Expect(reflect.DeepEqual(queriedBalance, newBalance)).Should(BeTrue())
					balance = newBalance
				}
			}
		})

		It("should use the cached balance when something wrong with the blockchain", func() {
			balance := randomBalance()
			done := make(chan struct{})
			defer close(done)
			bc := testutils.NewFaultyBlockchain(balance)
			balances := New(200*time.Millisecond, bc, Logger)
			queries := make(chan BalanceQuery)
			go balances.Run(done, queries)
			time.Sleep(10 * time.Millisecond)

			for i := 0; i < 10; i++ {
				time.Sleep(100 * time.Millisecond)
				responder := make(chan map[blockchain.TokenName]blockchain.Balance)
				query := BalanceQuery{
					Response: responder,
				}
				queries <- query
				queriedBalance := <-responder
				Expect(reflect.DeepEqual(queriedBalance, balance)).Should(BeTrue())
			}
		})

		It("close the query channel should stop the Balances from running", func() {
			done, _, _, balances := init(200 * time.Millisecond)
			defer close(done)
			queries := make(chan BalanceQuery)
			go balances.Run(done, queries)
			time.Sleep(10 * time.Millisecond)
			Expect(func() {
				close(queries)
				time.Sleep(10 * time.Millisecond)
			}).ShouldNot(Panic())
		})
	})
})

func randomBalance() map[blockchain.TokenName]blockchain.Balance {
	balancePara, ok := quick.Value(reflect.TypeOf(map[blockchain.TokenName]blockchain.Balance{}), Random)
	Expect(ok).Should(BeTrue())
	balance := balancePara.Interface().(map[blockchain.TokenName]blockchain.Balance)
	return balance
}