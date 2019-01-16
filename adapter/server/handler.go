package server

import (
	"bytes"
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

	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/core/swapper/status"
	"github.com/republicprotocol/swapperd/core/transfer"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

func NewErrBootloadRequired(msg string) error {
	return fmt.Errorf("please bootload before calling %s", msg)
}

type handler struct {
	bootloaded  map[string]bool
	swapperTask tau.Task
	walletTask  tau.Task
	wallet      wallet.Wallet
	logger      logrus.FieldLogger
}

// The Handler for swapperd requests
type Handler interface {
	GetID(password string) (GetIDResponse, error)
	GetInfo(password string) GetInfoResponse
	GetSwap(password string, id swap.SwapID) (GetSwapResponse, error)
	GetSwaps(password string) (GetSwapsResponse, error)
	GetBalances(password string) (GetBalancesResponse, error)
	GetAddresses(password string) (GetAddressesResponse, error)
	GetTransfers(password string) (GetTransfersResponse, error)
	GetJSONSignature(password string, message json.RawMessage) (GetSignatureResponseJSON, error)
	GetBase64Signature(password string, message string) (GetSignatureResponseString, error)
	GetHexSignature(password string, message string) (GetSignatureResponseString, error)
	PostTransfers(PostTransfersRequest) (PostTransfersResponse, error)
	PostSwaps(PostSwapRequest) (PostSwapResponse, error)
	PostDelayedSwaps(PostSwapRequest) error
	PostBootload(password string) error
}

func NewHandler(swapperTask, walletTask tau.Task, wallet wallet.Wallet, logger logrus.FieldLogger) Handler {
	return &handler{map[string]bool{}, swapperTask, walletTask, wallet, logger}
}

func (handler *handler) GetInfo(password string) GetInfoResponse {
	return GetInfoResponse{
		Version:         "0.3.0",
		Bootloaded:      handler.bootloaded[passwordHash(password)],
		SupportedTokens: handler.wallet.SupportedTokens(),
	}
}

func (handler *handler) GetAddresses(password string) (GetAddressesResponse, error) {
	return handler.wallet.Addresses(password)
}

func (handler *handler) GetSwaps(password string) (GetSwapsResponse, error) {
	resp := GetSwapsResponse{}
	if !handler.bootloaded[passwordHash(password)] {
		return resp, NewErrBootloadRequired("get swaps")
	}

	responder := make(chan map[swap.SwapID]swap.SwapReceipt)
	handler.swapperTask.IO().InputWriter() <- status.ReceiptQuery{Responder: responder}
	swapReceipts := <-responder

	for _, receipt := range swapReceipts {
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
	if !handler.bootloaded[passwordHash(password)] {
		return nil, NewErrBootloadRequired("get swaps")
	}
	responder := make(chan map[swap.SwapID]swap.SwapReceipt)
	handler.swapperTask.IO().InputWriter() <- status.ReceiptQuery{Responder: responder}
	swapReceipts := <-responder
	return swapReceipts, nil
}

func (handler *handler) GetBalances(password string) (GetBalancesResponse, error) {
	balanceMap, err := handler.wallet.Balances(password)
	return GetBalancesResponse(balanceMap), err
}

func (handler *handler) GetTransfers(password string) (GetTransfersResponse, error) {
	responder := make(chan transfer.TransferReceiptMap, 1)
	handler.walletTask.IO().InputWriter() <- transfer.TransferReceiptRequest{
		Responder: responder,
	}
	response := <-responder

	receiptMap := transfer.TransferReceiptMap{}
	for key, receipt := range response {
		passwordHash, err := base64.StdEncoding.DecodeString(receipt.PasswordHash)
		if receipt.PasswordHash != "" && err != nil {
			return GetTransfersResponse{}, err
		}

		if receipt.PasswordHash != "" && bcrypt.CompareHashAndPassword(passwordHash, []byte(password)) != nil {
			continue
		}

		receiptMap[key] = receipt
	}
	return MarshalGetTransfersResponse(receiptMap), nil
}

func (handler *handler) PostSwaps(swapReq PostSwapRequest) (PostSwapResponse, error) {
	if !handler.bootloaded[passwordHash(swapReq.Password)] {
		return PostSwapResponse{}, NewErrBootloadRequired("post swaps")
	}

	blob, err := handler.patchSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return PostSwapResponse{}, err
	}

	handler.swapperTask.IO().InputWriter() <- swapper.SwapRequest(blob)
	return handler.buildSwapResponse(blob)
}

func (handler *handler) PostDelayedSwaps(swapReq PostSwapRequest) error {
	if !handler.bootloaded[passwordHash(swapReq.Password)] {
		return NewErrBootloadRequired("post swaps")
	}

	blob, err := handler.patchDelayedSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return err
	}

	blob, err = handler.signDelayInfo(blob)
	if err != nil {
		return err
	}

	handler.swapperTask.IO().InputWriter() <- swapper.SwapRequest(swapReq)
	return nil
}

func (handler *handler) PostTransfers(req PostTransfersRequest) (PostTransfersResponse, error) {
	response := PostTransfersResponse{}
	token, err := blockchain.PatchToken(req.Token)
	if err != nil {
		return response, err
	}

	responder := make(chan transfer.TransferReceipt, 1)
	if err := handler.wallet.VerifyAddress(token.Blockchain, req.To); err != nil {
		return response, err
	}

	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return response, fmt.Errorf("invalid amount %s", req.Amount)
	}

	if err := handler.wallet.VerifyBalance(req.Password, token, amount); err != nil {
		return response, err
	}

	fee, err := handler.wallet.DefaultFee(token.Blockchain)
	if err != nil {
		return response, err
	}

	handler.walletTask.IO().InputWriter() <- transfer.NewTransferRequest(req.Password, token, req.To, amount, fee, responder)
	transferReceipt := <-responder
	return PostTransfersResponse(transferReceipt), nil
}

func (handler *handler) PostBootload(password string) error {
	if handler.bootloaded[passwordHash(password)] {
		return fmt.Errorf("already bootloaded")
	}
	handler.swapperTask.IO().InputWriter() <- swapper.Bootload{password}
	handler.bootloaded[passwordHash(password)] = true
	return nil
}

func (handler *handler) GetID(password string) (GetIDResponse, error) {
	id, err := handler.wallet.ID(password)
	if err != nil {
		return GetIDResponse{}, err
	}

	return GetIDResponse{
		PublicKey: id,
	}, nil
}

func (handler *handler) GetJSONSignature(password string, message json.RawMessage) (GetSignatureResponseJSON, error) {
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

func (handler *handler) patchSwap(swapBlob swap.SwapBlob) (swap.SwapBlob, error) {
	sendToken, err := blockchain.PatchToken(string(swapBlob.SendToken))
	if err != nil {
		return swapBlob, err
	}

	if err := handler.wallet.VerifyAddress(sendToken.Blockchain, swapBlob.SendTo); err != nil {
		return swapBlob, err
	}

	receiveToken, err := blockchain.PatchToken(string(swapBlob.ReceiveToken))
	if err != nil {
		return swapBlob, err
	}

	if err := handler.wallet.VerifyAddress(receiveToken.Blockchain, swapBlob.ReceiveFrom); err != nil {
		return swapBlob, err
	}

	if err := handler.verifySendAmount(swapBlob.Password, sendToken, swapBlob.SendAmount); err != nil {
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

	sendToken, err := blockchain.PatchToken(string(blob.SendToken))
	if err != nil {
		return blob, err
	}
	if err := handler.verifySendAmount(blob.Password, sendToken, blob.SendAmount); err != nil {
		return blob, err
	}

	receiveToken, err := blockchain.PatchToken(string(blob.ReceiveToken))
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

func (handler *handler) verifySendAmount(password string, token blockchain.Token, amount string) error {
	sendAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return fmt.Errorf("invalid send amount")
	}
	return handler.wallet.VerifyBalance(password, token, sendAmount)
}

func (handler *handler) verifyReceiveAmount(password string, token blockchain.Token) error {
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

	sendToken, err := blockchain.PatchToken(string(responseBlob.SendToken))
	if err != nil {
		return swapResponse, err
	}

	receiveToken, err := blockchain.PatchToken(string(responseBlob.ReceiveToken))
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

	blobBytes, err := json.Marshal(blob)
	if err != nil {
		return swapResponse, err
	}

	blobSig, err := handler.sign(blob.Password, blobBytes)
	if err != nil {
		return swapResponse, err
	}

	if blob.ShouldInitiateFirst {
		swapResponse.Swap = responseBlob
		swapResponse.Signature = base64.StdEncoding.EncodeToString(blobSig)
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

			handler.logger.Errorf("unexpected response", string(respBytes))
			return swapResponse, fmt.Errorf("unexpected status code while"+
				"posting to the response url: %d", resp.StatusCode)
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
