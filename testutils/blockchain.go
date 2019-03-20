package testutils

import (
	"errors"
	"math/big"
	"sync"

	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/tokens"
)

// MockBlockchain implements the `balance.Blockchain` interface.
type MockBlockchain struct {
	mu      *sync.Mutex
	balance map[tokens.Name]blockchain.Balance
}

// NewMockBlockchain creates a new `MockBlockchain`.
func NewMockBlockchain(balance map[tokens.Name]blockchain.Balance) *MockBlockchain {
	return &MockBlockchain{
		mu:      new(sync.Mutex),
		balance: copyBalanceMap(balance),
	}
}

// Balances implements the `balance.Blockchain` interface.
func (blockchain *MockBlockchain) Balances() (map[tokens.Name]blockchain.Balance, error) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	return copyBalanceMap(blockchain.balance), nil
}

// UpdateBalance with given data.
func (blockchain *MockBlockchain) UpdateBalance(balance map[tokens.Name]blockchain.Balance) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	blockchain.balance = copyBalanceMap(balance)
}

func (bc *MockBlockchain) GetAddress(password string, blockchainName tokens.BlockchainName) (string, error) {
	return "", nil
}
func (bc *MockBlockchain) Transfer(password string, token tokens.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error) {
	return "", blockchain.Cost{}, nil
}

func (bc *MockBlockchain) Lookup(token tokens.Token, txHash string) (transfer.UpdateReceipt, error) {
	return transfer.UpdateReceipt{}, nil
}

type FaultyBlockchain struct {
	balance map[tokens.Name]blockchain.Balance
	counter int
}

func NewFaultyBlockchain(balance map[tokens.Name]blockchain.Balance) *FaultyBlockchain {
	return &FaultyBlockchain{
		balance: balance,
		counter: 0,
	}
}

func (blockchain *FaultyBlockchain) Balances() (map[tokens.Name]blockchain.Balance, error) {
	blockchain.counter++
	if blockchain.counter%2 != 0 {
		return blockchain.balance, nil
	}
	return nil, errors.New("cannot get the balance")
}

func copyBalanceMap(balance map[tokens.Name]blockchain.Balance) map[tokens.Name]blockchain.Balance {
	copied := map[tokens.Name]blockchain.Balance{}
	for i, j := range balance {
		copied[i] = j
	}
	return copied
}
