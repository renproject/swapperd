package testutils

import (
	"errors"
	"math/big"
	"sync"

	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
)

// MockBlockchain implements the `balance.Blockchain` interface.
type MockBlockchain struct {
	mu      *sync.Mutex
	balance map[blockchain.TokenName]blockchain.Balance
}

// NewMockBlockchain creates a new `MockBlockchain`.
func NewMockBlockchain(balance map[blockchain.TokenName]blockchain.Balance) *MockBlockchain {
	return &MockBlockchain{
		mu:      new(sync.Mutex),
		balance: copyBalanceMap(balance),
	}
}

// Balances implements the `balance.Blockchain` interface.
func (blockchain *MockBlockchain) Balances() (map[blockchain.TokenName]blockchain.Balance, error) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	return copyBalanceMap(blockchain.balance), nil
}

// UpdateBalance with given data.
func (blockchain *MockBlockchain) UpdateBalance(balance map[blockchain.TokenName]blockchain.Balance) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	blockchain.balance = copyBalanceMap(balance)
}

func (bc *MockBlockchain) GetAddress(password string, blockchainName blockchain.BlockchainName) (string, error) {
	return "", nil
}
func (bc *MockBlockchain) Transfer(password string, token blockchain.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error) {
	return "", blockchain.Cost{}, nil
}

func (bc *MockBlockchain) Lookup(token blockchain.Token, txHash string) (transfer.UpdateReceipt, error) {
	return transfer.UpdateReceipt{}, nil
}

type FaultyBlockchain struct {
	balance map[blockchain.TokenName]blockchain.Balance
	counter int
}

func NewFaultyBlockchain(balance map[blockchain.TokenName]blockchain.Balance) *FaultyBlockchain {
	return &FaultyBlockchain{
		balance: balance,
		counter: 0,
	}
}

func (blockchain *FaultyBlockchain) Balances() (map[blockchain.TokenName]blockchain.Balance, error) {
	blockchain.counter++
	if blockchain.counter%2 != 0 {
		return blockchain.balance, nil
	}
	return nil, errors.New("cannot get the balance")
}

func copyBalanceMap(balance map[blockchain.TokenName]blockchain.Balance) map[blockchain.TokenName]blockchain.Balance {
	copied := map[blockchain.TokenName]blockchain.Balance{}
	for i, j := range balance {
		copied[i] = j
	}
	return copied
}
