package transfer

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/tokens"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	PutTransfer(receipt TransferReceipt) error
	Transfers() ([]TransferReceipt, error)
}

type Blockchain interface {
	GetAddress(password string, blockchainName tokens.BlockchainName) (string, error)
	Transfer(password string, token tokens.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) (string, blockchain.Cost, error)
	Lookup(token tokens.Token, txHash string) (UpdateReceipt, error)
}

type transfers struct {
	blockchain Blockchain
	storage    Storage
	logger     logrus.FieldLogger
}

func New(cap int, bc Blockchain, storage Storage, logger logrus.FieldLogger) tau.Task {
	return tau.New(tau.NewIO(cap), &transfers{bc, storage, logger})
}

func (transfers *transfers) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case TransferRequest:
		return transfers.handleTransferRequest(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in transfers: %T", msg))
	}
}

func (transfers *transfers) handleTransferRequest(msg TransferRequest) tau.Message {
	from, err := transfers.blockchain.GetAddress(msg.Password, msg.Token.Blockchain)
	if err != nil {
		return tau.NewError(err)
	}
	txHash, txCost, err := transfers.blockchain.Transfer(msg.Password, msg.Token, msg.To, msg.Amount, msg.Speed, msg.SendAll)
	if err != nil {
		return tau.NewError(err)
	}
	receipt := buildReceipt(msg, from, txHash, txCost)
	if err := transfers.storage.PutTransfer(receipt); err != nil {
		return tau.NewError(err)
	}
	return nil
}

func buildReceipt(req TransferRequest, from, txHash string, txCost blockchain.Cost) TransferReceipt {
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
			TxCost: blockchain.CostToCostBlob(txCost),
			TxHash: txHash,
		},
	}
}

type TransferRequest struct {
	Password string
	Token    tokens.Token
	To       string
	Amount   *big.Int
	Speed    blockchain.TxExecutionSpeed
	SendAll  bool
}

func NewTransferRequest(password string, token tokens.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, sendAll bool) TransferRequest {
	return TransferRequest{password, token, to, amount, speed, sendAll}
}

func (request TransferRequest) IsMessage() {
}

type TransferReceiptMap map[string]TransferReceipt

type TransferReceipt struct {
	Confirmations int64  `json:"confirmations"`
	Timestamp     int64  `json:"timestamp"`
	PasswordHash  string `json:"passwordHash,omitempty"`
	TokenDetails
}

type TokenDetails struct {
	To     string              `json:"to"`
	From   string              `json:"from"`
	Token  tokens.Token        `json:"token"`
	Amount string              `json:"value"`
	TxCost blockchain.CostBlob `json:"txCost"`
	TxHash string              `json:"txHash"`
}

type UpdateReceipt struct {
	TxHash string
	Update func(*TransferReceipt)
}

func NewUpdateReceipt(txHash string, update func(*TransferReceipt)) UpdateReceipt {
	return UpdateReceipt{txHash, update}
}
