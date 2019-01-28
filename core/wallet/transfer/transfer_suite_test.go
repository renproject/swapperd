package transfer_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/swapperd/core/wallet/transfer"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func TestTransfer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transfer Suite")
}

type MockStorage struct {
	Rand *rand.Rand
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Rand: rand.New(rand.New(rand.NewSource(time.Now().Unix()))),
	}
}

func (storage MockStorage) PutTransfer(receipt TransferReceipt) error {
	if storage.Rand.Int63()%2 == 1 {
		return fmt.Errorf("failed to store the transfer receipt")
	}
	return nil
}

func (storage MockStorage) Transfers() ([]TransferReceipt, error) {
	if storage.Rand.Int63()%2 == 1 {
		return nil, fmt.Errorf("failed to store the transfer receipt")
	}
	return nil, nil
}

type MockBlockhain struct {
	Rand *rand.Rand
}

func NewMockBlockhain() *MockBlockhain {
	return &MockBlockhain{
		Rand: rand.New(rand.New(rand.NewSource(time.Now().Unix()))),
	}
}

func (bc MockBlockhain) GetAddress(password string, blockchainName blockchain.BlockchainName) (string, error) {
	if bc.Rand.Int63()%2 == 1 {
		return "", fmt.Errorf("corrupted keystore")
	}
	return "", nil
}
func (bc MockBlockhain) Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error) {
	if bc.Rand.Int63()%2 == 1 {
		return "", fmt.Errorf("failed to connect")
	}
	return "", nil
}

func (bc MockBlockhain) Lookup(token blockchain.Token, txHash string) (UpdateReceipt, error) {
	if bc.Rand.Int63()%2 == 1 {
		return UpdateReceipt{}, fmt.Errorf("failed to connect")
	}
	return UpdateReceipt{}, nil
}
