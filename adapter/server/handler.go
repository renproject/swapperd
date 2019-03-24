package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/renproject/swapperd/adapter/wallet"
	coreWallet "github.com/renproject/swapperd/core/wallet"
	"github.com/renproject/swapperd/core/wallet/swapper"
	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

var ErrHandlerIsShuttingDown = fmt.Errorf("Http handler is shutting down")

type handler struct {
	version    string
	bootloaded map[string]bool
	wallet     wallet.Wallet
	storage    Storage
	receiver   *Receiver
}

// The Handler for swapperd requests
type Handler interface {
	GetID(password string, idType string) (string, error)
	GetInfo(password string) GetInfoResponse
	GetSwap(password string, id swap.SwapID) (GetSwapResponse, error)
	GetSwaps(password string) (GetSwapsResponse, error)
	GetBalances(password string) (GetBalancesResponse, error)
	GetBalance(password string, token tokens.Token) (GetBalanceResponse, error)
	GetAddresses(password string) (GetAddressesResponse, error)
	GetAddress(password string, token tokens.Token) (GetAddressResponse, error)
	GetTransfers(password string) (GetTransfersResponse, error)
	GetJSONSignature(password string, message json.RawMessage) (GetSignatureResponseJSON, error)
	GetBase64Signature(password string, message string) (GetSignatureResponseString, error)
	GetHexSignature(password string, message string) (GetSignatureResponseString, error)
	PostTransfers(PostTransfersRequest) error
	PostSwaps(PostSwapRequest) (PostSwapResponse, error)
	PostDelayedSwaps(PostSwapRequest) error
	Shutdown()
}

func NewHandler(cap int, version string, wallet wallet.Wallet, storage Storage, receiver *Receiver) Handler {
	return &handler{
		version:    version,
		bootloaded: map[string]bool{},
		wallet:     wallet,
		storage:    storage,
		receiver:   receiver,
	}
}

func (handler *handler) GetInfo(password string) GetInfoResponse {
	handler.bootload(password)
	return GetInfoResponse{
		Version:         handler.version,
		Bootloaded:      handler.bootloaded[passwordHash(password)],
		SupportedTokens: handler.wallet.SupportedTokens(),
	}
}

func (handler *handler) Shutdown() {
	handler.receiver.Shutdown()
}

func (handler *handler) GetAddresses(password string) (GetAddressesResponse, error) {
	handler.bootload(password)
	return handler.wallet.Addresses(password)
}

func (handler *handler) GetAddress(password string, token tokens.Token) (GetAddressResponse, error) {
	handler.bootload(password)
	address, err := handler.wallet.GetAddress(password, token.Blockchain)
	return GetAddressResponse(address), err
}

func (handler *handler) GetSwaps(password string) (GetSwapsResponse, error) {
	handler.bootload(password)
	resp := GetSwapsResponse{}

	receipts, err := handler.getSwapReceipts(password)
	if err != nil {
		return resp, err
	}

	// swapReceipts := handler.getSwapReceipts(password)
	for _, receipt := range receipts {
		passwordHash, err := base64.StdEncoding.DecodeString(receipt.PasswordHash)
		if receipt.PasswordHash != "" && err != nil {
			return resp, fmt.Errorf("corrupted password")
		}

		if receipt.PasswordHash != "" && bcrypt.CompareHashAndPassword(passwordHash, []byte(password)) != nil {
			continue
		}
		receipt.PasswordHash = ""
		resp.Swaps = append(resp.Swaps, receipt)
	}

	return resp, nil
}

func (handler *handler) GetSwap(password string, id swap.SwapID) (GetSwapResponse, error) {
	handler.bootload(password)
	swapReceipts, err := handler.getSwapReceipts(password)
	if err != nil {
		return GetSwapResponse{}, err
	}
	receipt, ok := swapReceipts[id]
	if !ok {
		return GetSwapResponse{}, fmt.Errorf("swap receipt not found")
	}
	return GetSwapResponse(receipt), nil
}

func (handler *handler) getSwapReceipts(password string) (map[swap.SwapID]swap.SwapReceipt, error) {
	receiptMap := map[swap.SwapID]swap.SwapReceipt{}

	receipts, err := handler.storage.Receipts()
	if err != nil {
		return nil, err
	}

	for _, receipt := range receipts {
		receiptMap[receipt.ID] = receipt
	}

	return receiptMap, nil
}

func (handler *handler) GetBalances(password string) (GetBalancesResponse, error) {
	handler.bootload(password)
	balanceMap, err := handler.wallet.Balances(password)
	return GetBalancesResponse(balanceMap), err
}

func (handler *handler) GetBalance(password string, token tokens.Token) (GetBalanceResponse, error) {
	handler.bootload(password)
	balance, err := handler.wallet.Balance(password, token)
	return GetBalanceResponse(balance), err
}

func (handler *handler) GetTransfers(password string) (GetTransfersResponse, error) {
	handler.bootload(password)
	transfers, err := handler.storage.Transfers()
	if err != nil {
		return GetTransfersResponse{}, err
	}

	receiptMap := transfer.TransferReceiptMap{}
	for _, receipt := range transfers {
		passwordHash, err := base64.StdEncoding.DecodeString(receipt.PasswordHash)
		if receipt.PasswordHash != "" && err != nil {
			return GetTransfersResponse{}, err
		}

		if receipt.PasswordHash != "" && bcrypt.CompareHashAndPassword(passwordHash, []byte(password)) != nil {
			continue
		}

		update, err := handler.wallet.Lookup(receipt.Token, receipt.TxHash)
		if err != nil {
			return GetTransfersResponse{}, fmt.Errorf("Failed to lookup tx with txHash (%s) on %s blockchain", receipt.TxHash, receipt.Token.Blockchain)
		}

		update.Update(&receipt)
		receiptMap[receipt.TxHash] = receipt
	}
	return MarshalGetTransfersResponse(receiptMap), nil
}

func (handler *handler) PostSwaps(swapReq PostSwapRequest) (PostSwapResponse, error) {
	handler.bootload(swapReq.Password)

	blob, err := handler.patchSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return PostSwapResponse{}, err
	}

	if err := handler.Write(swapper.SwapRequest(blob)); err != nil {
		return PostSwapResponse{}, err
	}
	return handler.buildSwapResponse(blob)
}

func (handler *handler) PostDelayedSwaps(swapReq PostSwapRequest) error {
	handler.bootload(swapReq.Password)

	blob, err := handler.patchDelayedSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return err
	}

	blob, err = handler.signDelayInfo(blob)
	if err != nil {
		return err
	}

	return handler.Write(swapper.SwapRequest(blob))
}

func (handler *handler) PostTransfers(req PostTransfersRequest) error {
	handler.bootload(req.Password)
	if req.Speed == blockchain.Nil {
		req.Speed = blockchain.Fast
	}

	token, err := tokens.PatchToken(req.Token)
	if err != nil {
		return err
	}
	if err := handler.wallet.VerifyAddress(token.Blockchain, req.To); err != nil {
		return err
	}
	if req.SendAll {
		balance, err := handler.wallet.Balance(req.Password, token)
		if err != nil {
			return err
		}
		amount, ok := new(big.Int).SetString(balance.Amount, 10)
		if !ok {
			return fmt.Errorf("unable to decode balance: %s", balance.Amount)
		}
		return handler.Write(transfer.NewTransferRequest(req.Password, token, req.To, amount, req.Speed, true))
	}

	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount %s", req.Amount)
	}

	if err := handler.wallet.VerifyBalance(req.Password, token, amount); err != nil {
		return err
	}

	return handler.Write(transfer.NewTransferRequest(req.Password, token, req.To, amount, req.Speed, false))
}

func (handler *handler) GetID(password, idType string) (string, error) {
	handler.bootload(password)
	id, err := handler.wallet.ID(password, idType)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (handler *handler) GetJSONSignature(password string, message json.RawMessage) (GetSignatureResponseJSON, error) {
	handler.bootload(password)
	sig, err := handler.sign(password, message)
	if err != nil {
		return GetSignatureResponseJSON{}, err
	}
	return GetSignatureResponseJSON{
		Message:   message,
		Signature: base64.StdEncoding.EncodeToString(sig),
	}, nil
}

func (handler *handler) GetBase64Signature(password string, message string) (GetSignatureResponseString, error) {
	handler.bootload(password)
	msg, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return GetSignatureResponseString{}, err
	}

	sig, err := handler.sign(password, msg)
	if err != nil {
		return GetSignatureResponseString{}, err
	}

	return GetSignatureResponseString{
		Message:   message,
		Signature: base64.StdEncoding.EncodeToString(sig),
	}, nil
}

func (handler *handler) GetHexSignature(password string, message string) (GetSignatureResponseString, error) {
	handler.bootload(password)
	if len(message) > 2 && message[:2] == "0x" {
		message = message[2:]
	}
	msg, err := hex.DecodeString(message)
	if err != nil {
		return GetSignatureResponseString{}, err
	}

	sig, err := handler.sign(password, msg)
	if err != nil {
		return GetSignatureResponseString{}, err
	}

	return GetSignatureResponseString{
		Message:   message,
		Signature: hex.EncodeToString(sig),
	}, nil
}

func (handler *handler) Write(msg tau.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return handler.receiver.Write(ctx, msg)
}

func (handler *handler) bootload(password string) {
	if !handler.bootloaded[passwordHash(password)] {
		if err := handler.Write(coreWallet.Bootload{password}); err != nil {
			return
		}
		handler.bootloaded[passwordHash(password)] = true
	}
}

func (handler *handler) patchSwap(swapBlob swap.SwapBlob) (swap.SwapBlob, error) {
	sendToken, err := tokens.PatchToken(swapBlob.SendToken)
	if err != nil {
		return swapBlob, err
	}

	if err := handler.wallet.VerifyAddress(sendToken.Blockchain, swapBlob.SendTo); err != nil {
		return swapBlob, err
	}

	receiveToken, err := tokens.PatchToken(swapBlob.ReceiveToken)
	if err != nil {
		return swapBlob, err
	}

	if err := handler.wallet.VerifyAddress(receiveToken.Blockchain, swapBlob.ReceiveFrom); err != nil {
		return swapBlob, err
	}

	if swapBlob.WithdrawAddress != "" {
		if err := handler.wallet.VerifyAddress(receiveToken.Blockchain, swapBlob.WithdrawAddress); err != nil {
			return swapBlob, err
		}
	}

	if err := handler.verifySendAmount(swapBlob.Password, sendToken, swapBlob.SendAmount, swapBlob.BrokerFee); err != nil {
		return swapBlob, err
	}

	if err := handler.verifyReceiveAmount(swapBlob.Password, receiveToken); err != nil {
		return swapBlob, err
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])
	swapBlob.ID = swap.SwapID(base64.StdEncoding.EncodeToString(swapID[:]))
	secret := [32]byte{}
	if swapBlob.ShouldInitiateFirst {
		swapBlob.TimeLock = time.Now().Unix() + 3*swap.ExpiryUnit
		secret = genereateSecret(swapBlob.Password, swapBlob.ID)
		hash := sha256.Sum256(secret[:])
		swapBlob.SecretHash = base64.StdEncoding.EncodeToString(hash[:])
		return swapBlob, nil
	}

	secretHash, err := base64.StdEncoding.DecodeString(swapBlob.SecretHash)
	if len(secretHash) != 32 || err != nil {
		return swapBlob, fmt.Errorf("invalid secret hash")
	}
	if time.Now().Unix()+2*swap.ExpiryUnit > swapBlob.TimeLock {
		return swapBlob, fmt.Errorf("not enough time to do the atomic swap")
	}
	return swapBlob, nil
}

func (handler *handler) patchDelayedSwap(blob swap.SwapBlob) (swap.SwapBlob, error) {
	if blob.DelayCallbackURL == "" {
		return blob, fmt.Errorf("delay url cannot be empty")
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])
	blob.ID = swap.SwapID(base64.StdEncoding.EncodeToString(swapID[:]))

	sendToken, err := tokens.PatchToken(string(blob.SendToken))
	if err != nil {
		return blob, err
	}
	if err := handler.verifySendAmount(blob.Password, sendToken, blob.SendAmount, blob.BrokerFee); err != nil {
		return blob, err
	}

	receiveToken, err := tokens.PatchToken(string(blob.ReceiveToken))
	if err != nil {
		return blob, err
	}
	if err := handler.verifyReceiveAmount(blob.Password, receiveToken); err != nil {
		return blob, err
	}

	secret := genereateSecret(blob.Password, blob.ID)
	secretHash := sha256.Sum256(secret[:])
	blob.SecretHash = base64.StdEncoding.EncodeToString(secretHash[:])
	blob.TimeLock = time.Now().Unix() + 3*swap.ExpiryUnit
	return blob, nil
}

func (handler *handler) verifySendAmount(password string, token tokens.Token, amount string, fee int64) error {
	sendAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return fmt.Errorf("invalid send amount")
	}
	return handler.wallet.VerifyBalance(password, token, withFees(sendAmount, fee))
}

func (handler *handler) verifyReceiveAmount(password string, token tokens.Token) error {
	return handler.wallet.VerifyBalance(password, token, nil)
}

func (handler *handler) signDelayInfo(blob swap.SwapBlob) (swap.SwapBlob, error) {
	delayInfoSig, err := handler.sign(blob.Password, blob.DelayInfo)
	if err != nil {
		return blob, fmt.Errorf("failed to sign delay info: %v", err)
	}

	signedDelayInfo, err := json.Marshal(struct {
		Message   json.RawMessage `json:"message"`
		Signature string          `json:"signature"`
	}{
		Message:   blob.DelayInfo,
		Signature: base64.StdEncoding.EncodeToString(delayInfoSig),
	})
	if err != nil {
		return blob, fmt.Errorf("unable to marshal signed delay info: %v", err)
	}

	blob.DelayInfo = signedDelayInfo
	return blob, nil
}

func (handler *handler) buildSwapResponse(blob swap.SwapBlob) (PostSwapResponse, error) {
	responseBlob := swap.SwapBlob{}
	responseBlob.SendToken = blob.ReceiveToken
	responseBlob.ReceiveToken = blob.SendToken
	responseBlob.SendAmount = blob.ReceiveAmount
	responseBlob.ReceiveAmount = blob.SendAmount
	swapResponse := PostSwapResponse{}

	sendToken, err := tokens.PatchToken(responseBlob.SendToken)
	if err != nil {
		return swapResponse, err
	}

	receiveToken, err := tokens.PatchToken(responseBlob.ReceiveToken)
	if err != nil {
		return swapResponse, err
	}

	sendTo, err := handler.wallet.GetAddress(blob.Password, sendToken.Blockchain)
	if err != nil {
		return swapResponse, err
	}

	receiveFrom, err := handler.wallet.GetAddress(blob.Password, receiveToken.Blockchain)
	if err != nil {
		return swapResponse, err
	}

	responseBlob.SendTo = sendTo
	responseBlob.ReceiveFrom = receiveFrom
	responseBlob.SecretHash = blob.SecretHash
	responseBlob.TimeLock = blob.TimeLock

	responseBlob.BrokerFee = blob.BrokerFee
	responseBlob.BrokerSendTokenAddr = blob.BrokerReceiveTokenAddr
	responseBlob.BrokerReceiveTokenAddr = blob.BrokerSendTokenAddr

	responseBlobBytes, err := json.Marshal(responseBlob)
	if err != nil {
		return swapResponse, err
	}

	responseBlobSig, err := handler.sign(blob.Password, responseBlobBytes)
	if err != nil {
		return swapResponse, err
	}

	if blob.ShouldInitiateFirst {
		swapResponse.Swap = responseBlob
		swapResponse.Signature = base64.StdEncoding.EncodeToString(responseBlobSig)
	}

	if blob.ResponseURL != "" {
		data, err := json.MarshalIndent(swapResponse, "", "  ")
		if err != nil {
			return swapResponse, err
		}
		buf := bytes.NewBuffer(data)

		resp, err := http.Post(blob.ResponseURL, "application/json", buf)
		if err != nil {
			return swapResponse, err
		}

		if resp.StatusCode != 200 {
			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return swapResponse, err
			}

			return swapResponse, fmt.Errorf("unexpected status code (%d) while"+
				"posting to the response url: %s", resp.StatusCode, respBytes)
		}
	}
	swapResponse.ID = blob.ID
	return swapResponse, nil
}

func (handler *handler) sign(password string, message []byte) ([]byte, error) {
	signer, err := handler.wallet.ECDSASigner(password)
	if err != nil {
		return nil, fmt.Errorf("unable to load ecdsa signer: %v", err)
	}
	hash := sha3.Sum256(message)
	sig, err := signer.Sign(hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign swap response: %v", err)
	}
	return sig, nil
}

func genereateSecret(password string, id swap.SwapID) [32]byte {
	return sha3.Sum256(append([]byte(password), []byte(id)...))
}

func passwordHash(password string) string {
	passwordHash32 := sha3.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(passwordHash32[:])
}

func withFees(amount *big.Int, bips int64) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(bips+10000)), big.NewInt(10000))
}
