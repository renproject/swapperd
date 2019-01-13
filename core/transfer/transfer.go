package transfer

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	PutTransfer(receipt TransferReceipt) error
	Transfers() ([]TransferReceipt, error)
}

type Blockchain interface {
	GetAddress(password string, blockchainName blockchain.BlockchainName) (string, error)
	Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error)
	Lookup(token blockchain.Token, txHash string) (UpdateReceipt, error)
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
	case Bootload:
		return transfers.handleBootload()
	case TransferReceiptRequest:
		return transfers.handleTransferReceiptRequest(msg)
	case TransferRequest:
		return transfers.handleTransferRequest(msg)
	case tau.Tick:
		return transfers.handleTick()
	default:
		return tau.NewError(fmt.Errorf("invalid message type in transfers: %T", msg))
	}
}

func (transfers *transfers) handleTick() tau.Message {
	transfers.update()
	return nil
}

func (transfers *transfers) handleBootload() tau.Message {
	transferReceipts, err := transfers.storage.Transfers()
	if err != nil {
		return tau.NewError(err)
	}
	for _, transferReceipt := range transferReceipts {
		transfers.write(transferReceipt)
	}
	transfers.update()
	return nil
}

func (transfers *transfers) handleTransferReceiptRequest(msg TransferReceiptRequest) tau.Message {
	msg.Responder <- transfers.read()
	return nil
}

func (transfers *transfers) handleTransferRequest(msg TransferRequest) tau.Message {
	from, err := transfers.blockchain.GetAddress(msg.Password, msg.Token.Blockchain)
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
	if err := transfers.storage.PutTransfer(receipt); err != nil {
		return tau.NewError(err)
	}
	return nil
}

func (transfers *transfers) update() {
	updatedTransferMap := TransferReceiptMap{}
	for txHash, receipt := range transfers.transferMap {
		update, err := transfers.blockchain.Lookup(receipt.Token, txHash)
		if err != nil {
			transfers.logger.Error(err)
			continue
		}
		update.Update(&receipt)
		updatedTransferMap[txHash] = receipt
	}
	transfers.mu.Lock()
	defer transfers.mu.Unlock()
	transfers.transferMap = updatedTransferMap
}

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

func buildReceipt(req TransferRequest, from, txHash string) TransferReceipt {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	return TransferReceipt{
		Confirmations: 0,
		Timestamp:     time.Now().Unix(),
		PasswordHash:  base64.StdEncoding.EncodeToString(passwordHash),
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

type TransferRequest struct {
	Password string
	Token    blockchain.Token
	To       string
	Amount   *big.Int
	Fee      *big.Int

	Responder chan<- TransferReceipt
}

func NewTransferRequest(password string, token blockchain.Token, to string, amount, fee *big.Int, responder chan<- TransferReceipt) TransferRequest {
	return TransferRequest{password, token, to, amount, fee, responder}
}

func (request TransferRequest) IsMessage() {
}

type TransferReceiptMap map[string]TransferReceipt

func (request TransferReceiptMap) IsMessage() {
}

type TransferReceipt struct {
	Confirmations int64  `json:"confirmations"`
	Timestamp     int64  `json:"timestamp"`
	PasswordHash  string `json:"passwordHash,omitempty"`
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

type UpdateReceipt struct {
	TxHash string
	Update func(*TransferReceipt)
}

func NewUpdateReceipt(txHash string, update func(*TransferReceipt)) UpdateReceipt {
	return UpdateReceipt{txHash, update}
}

type TransferReceiptRequest struct {
	Responder chan<- TransferReceiptMap
}

func (request TransferReceiptRequest) IsMessage() {
}

type Bootload struct {
}

func (request Bootload) IsMessage() {
}
