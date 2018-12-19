package transfer

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	PutTransfer(receipt TransferReceipt) error
	Transfers() ([]TransferReceipt, error)
}

type Blockchain interface {
	GetAddress(blockchain blockchain.BlockchainName) (string, error)
	Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error)
	// Confirmations(txHash string) (int64, error)
}

type transfers struct {
	mu          *sync.RWMutex
	transferMap TransferReceiptMap
	logger      logrus.FieldLogger
	blockchain  Blockchain
	storage     Storage
}

func New(cap int, bc Blockchain, storage Storage, logger logrus.FieldLogger) tau.Task {
	return tau.New(tau.NewIO(cap), &transfers{new(sync.RWMutex), TransferReceiptMap{}, logger, bc, storage})
}

func (transfers *transfers) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	// case tau.Tick:
	// 	go transfers.update()
	// 	return nil
	case TransferRequest:
		from, err := transfers.blockchain.GetAddress(msg.Token.Blockchain)
		if err != nil {
			return tau.NewError(err)
		}
		txHash, err := transfers.blockchain.Transfer(msg.Password, msg.Token, msg.To, msg.Amount)
		if err != nil {
			return tau.NewError(err)
		}
		receipt := buildReceipt(msg, from, txHash)
		transfers.write(receipt)
		msg.Responder <- receipt
		return nil
	default:
		return tau.NewError(fmt.Errorf("invalid message type in transfers: %T", msg))
	}
}

func buildReceipt(req TransferRequest, from, txHash string) TransferReceipt {
	return TransferReceipt{
		Confirmations: 0,
		Timestamp:     time.Now().Unix(),
		TokenDetails: TokenDetails{
			To:     req.To,
			From:   from,
			Token:  req.Token,
			Amount: req.Amount.String(),
			Fee:    req.Fee.String(),
			TxHash: txHash,
		},
	}
}

// func (transfers *transfers) update() {
// 	balanceMap, err := transfers.blockchain.Balances()
// 	if err != nil {
// 		transfers.logger.Errorf("cannot update transfers: %v", err)
// 		return
// 	}
// 	transfers.mu.Lock()
// 	defer transfers.mu.Unlock()
// 	transfers.balanceMap = balanceMap
// }

func (transfers *transfers) write(receipt TransferReceipt) {
	transfers.mu.Lock()
	defer transfers.mu.Unlock()
	transfers.transferMap[receipt.TxHash] = receipt
}

func (transfers *transfers) read() TransferReceiptMap {
	transfers.mu.RLock()
	defer transfers.mu.RUnlock()
	return transfers.transferMap
}

type TransferRequest struct {
	Password string
	Token    blockchain.Token
	To       string
	Amount   *big.Int
	Fee      *big.Int

	Responder chan<- TransferReceipt
}

func (request TransferRequest) IsMessage() {
}

type TransferReceiptMap map[string]TransferReceipt

func (request TransferReceiptMap) IsMessage() {
}

type TransferReceipt struct {
	Confirmations int64 `json:"confirmations"`
	Timestamp     int64 `json:"timestamp"`
	TokenDetails
}

type TokenDetails struct {
	To     string           `json:"to"`
	From   string           `json:"from"`
	Token  blockchain.Token `json:"token"`
	Amount string           `json:"value"`
	Fee    string           `json:"fee"`
	TxHash string           `json:"txHash"`
}
