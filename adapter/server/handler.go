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

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/wallet/balance"
	"github.com/republicprotocol/swapperd/core/wallet/transfer"
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
	bootloaded   bool
	passwordHash []byte
	walletIO     tau.IO
	wallet       wallet.Wallet
	storage      Storage
	logger       logrus.FieldLogger
}

// The Handler for swapperd requests
type Handler interface {
	GetID() GetIDResponse
	GetInfo() GetInfoResponse
	GetSwaps(chan<- status.ReceiptQuery) (GetSwapsResponse, error)
	GetBalances() GetBalancesResponse
	GetAddresses() (GetAddressesResponse, error)
	GetTransfers() GetTransfersResponse
	GetJSONSignature(password string, message json.RawMessage) (GetSignatureResponseJSON, error)
	GetBase64Signature(password string, message string) (GetSignatureResponseString, error)
	GetHexSignature(password string, message string) (GetSignatureResponseString, error)
	PostTransfers(PostTransfersRequest) (PostTransfersResponse, error)
	PostSwaps(PostSwapRequest, chan<- swap.SwapReceipt, chan<- swap.SwapBlob) (PostSwapResponse, error)
	PostDelayedSwaps(PostSwapRequest, chan<- swap.SwapReceipt, chan<- swap.SwapBlob) error
	PostBootload(password string, swaps, delayedSwaps chan<- swap.SwapBlob) error
	VerifyPassword(password string) bool
}

func NewHandler(passwordHash []byte, walletIO tau.IO, wallet wallet.Wallet, storage Storage, logger logrus.FieldLogger) Handler {
	return &handler{false, passwordHash, walletIO, wallet, storage, logger}
}

func (handler *handler) GetInfo() GetInfoResponse {
	return GetInfoResponse{
		Version:              "0.2.0",
		Bootloaded:           handler.bootloaded,
		SupportedBlockchains: handler.wallet.SupportedBlockchains(),
		SupportedTokens:      handler.wallet.SupportedTokens(),
	}
}

func (handler *handler) GetAddresses() (GetAddressesResponse, error) {
	return handler.wallet.Addresses()
}

func (handler *handler) GetSwaps(statuses chan<- status.ReceiptQuery) (GetSwapsResponse, error) {
	resp := GetSwapsResponse{}
	if !handler.bootloaded {
		return resp, NewErrBootloadRequired("get swaps")
	}
	responder := make(chan map[swap.SwapID]swap.SwapReceipt)
	statuses <- status.ReceiptQuery{Responder: responder}
	statusMap := <-responder
	for _, status := range statusMap {
		resp.Swaps = append(resp.Swaps, status)
	}
	return resp, nil
}

func (handler *handler) GetBalances() GetBalancesResponse {
	responder := make(chan balance.BalanceMap, 1)
	balanceReq := balance.BalanceRequest{
		Responder: responder,
	}
	handler.walletIO.InputWriter() <- balanceReq
	response := <-responder
	return GetBalancesResponse(response)
}

func (handler *handler) GetTransfers() GetTransfersResponse {
	responder := make(chan transfer.TransferReceiptMap, 1)
	transferReq := transfer.TransferReceiptRequest{
		Responder: responder,
	}
	handler.walletIO.InputWriter() <- transferReq
	response := <-responder
	return MarshalGetTransfersResponse(response)
}

func (handler *handler) PostSwaps(swapReq PostSwapRequest, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) (PostSwapResponse, error) {
	if !handler.bootloaded {
		return PostSwapResponse{}, NewErrBootloadRequired("post swaps")
	}
	password := swapReq.Password
	blob, err := handler.patchSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return PostSwapResponse{}, err
	}

	receipt := swap.NewSwapReceipt(blob)
	blob.Password = ""
	if err := handler.storage.PutSwap(blob); err != nil {
		return PostSwapResponse{}, err
	}

	blob.Password = password
	go func() {
		swaps <- blob
		receipts <- receipt
	}()
	return handler.BuildSwapResponse(blob)
}

func (handler *handler) PostDelayedSwaps(swapReq PostSwapRequest, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) error {
	if !handler.bootloaded {
		return NewErrBootloadRequired("post swaps")
	}
	password := swapReq.Password

	blob, err := handler.patchDelayedSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return err
	}

	blob, err = handler.signDelayInfo(blob)
	if err != nil {
		return err
	}

	blob.Password = ""
	receipt := swap.NewSwapReceipt(blob)
	if err := handler.storage.PutSwap(blob); err != nil {
		return err
	}

	blob.Password = password
	go func() {
		swaps <- blob
		receipts <- receipt
	}()

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

	if err := handler.wallet.VerifyBalance(token, amount); err != nil {
		return response, err
	}

	fee, err := handler.wallet.DefaultFee(token.Blockchain)
	if err != nil {
		return response, err
	}

	transferReq := transfer.NewTransferRequest(req.Password, token, req.To, amount, fee, responder)
	handler.walletIO.InputWriter() <- transferReq

	transferReceipt := <-responder
	return PostTransfersResponse(transferReceipt), nil
}

func (handler *handler) PostBootload(password string, swaps, delayedSwaps chan<- swap.SwapBlob) error {
	if handler.bootloaded {
		return fmt.Errorf("already bootloaded")
	}

	pendingSwaps, err := handler.storage.PendingSwaps()
	if err != nil {
		return err
	}

	handler.logger.Infof("loading %d pending atomic swaps", len(pendingSwaps))

	co.ParForAll(pendingSwaps, func(i int) {
		swap := pendingSwaps[i]
		swap.Password = password
		if swap.Delay {
			delayedSwaps <- swap
			return
		}
		swaps <- swap
	})

	handler.bootloaded = true
	return nil
}

func (handler *handler) GetID() GetIDResponse {
	return GetIDResponse{
		PublicKey: handler.wallet.ID(),
	}
}

func (handler *handler) GetJSONSignature(password string, message json.RawMessage) (GetSignatureResponseJSON, error) {
	sig, err := handler.Sign(password, message)
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

	sig, err := handler.Sign(password, msg)
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

	sig, err := handler.Sign(password, msg)
	if err != nil {
		return GetSignatureResponseString{}, err
	}

	return GetSignatureResponseString{
		Message:   message,
		Signature: hex.EncodeToString(sig),
	}, nil
}

func (handler *handler) VerifyPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword(handler.passwordHash, []byte(password)); err != nil {
		handler.logger.Info("password length", len(handler.passwordHash))
		handler.logger.Error(err)
	}
	return (bcrypt.CompareHashAndPassword(handler.passwordHash, []byte(password)) == nil)
}

func (handler *handler) patchSwap(swapBlob swap.SwapBlob) (swap.SwapBlob, error) {
	sendToken, err := blockchain.PatchToken(swapBlob.SendToken)
	if err != nil {
		return swapBlob, err
	}

	if err := handler.wallet.VerifyAddress(sendToken.Blockchain, swapBlob.SendTo); err != nil {
		return swapBlob, err
	}

	receiveToken, err := blockchain.PatchToken(swapBlob.ReceiveToken)
	if err != nil {
		return swapBlob, err
	}

	if err := handler.wallet.VerifyAddress(receiveToken.Blockchain, swapBlob.ReceiveFrom); err != nil {
		return swapBlob, err
	}

	if err := handler.verifySendAmount(sendToken, swapBlob.SendAmount); err != nil {
		return swapBlob, err
	}

	if err := handler.verifyReceiveAmount(receiveToken); err != nil {
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

	sendToken, err := blockchain.PatchToken(blob.SendToken)
	if err != nil {
		return blob, err
	}
	if err := handler.verifySendAmount(sendToken, blob.SendAmount); err != nil {
		return blob, err
	}

	receiveToken, err := blockchain.PatchToken(blob.ReceiveToken)
	if err != nil {
		return blob, err
	}
	if err := handler.verifyReceiveAmount(receiveToken); err != nil {
		return blob, err
	}

	secret := genereateSecret(blob.Password, blob.ID)
	secretHash := sha256.Sum256(secret[:])
	blob.SecretHash = base64.StdEncoding.EncodeToString(secretHash[:])
	blob.TimeLock = time.Now().Unix() + 3*swap.ExpiryUnit
	return blob, nil
}

func (handler *handler) verifySendAmount(token blockchain.Token, amount string) error {
	sendAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return fmt.Errorf("invalid send amount")
	}
	return handler.wallet.VerifyBalance(token, sendAmount)
}

func (handler *handler) verifyReceiveAmount(token blockchain.Token) error {
	return handler.wallet.VerifyBalance(token, nil)
}

func (handler *handler) signDelayInfo(blob swap.SwapBlob) (swap.SwapBlob, error) {
	signer, err := handler.wallet.ECDSASigner(blob.Password)
	if err != nil {
		return blob, fmt.Errorf("unable to load ecdsa signer: %v", err)
	}

	delayInfoHash := sha3.Sum256(blob.DelayInfo)
	delayInfoSig, err := signer.Sign(delayInfoHash[:])
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

func (handler *handler) BuildSwapResponse(blob swap.SwapBlob) (PostSwapResponse, error) {
	responseBlob := swap.SwapBlob{}
	responseBlob.SendToken = blob.ReceiveToken
	responseBlob.ReceiveToken = blob.SendToken
	responseBlob.SendAmount = blob.ReceiveAmount
	responseBlob.ReceiveAmount = blob.SendAmount
	swapResponse := PostSwapResponse{}

	sendToken, err := blockchain.PatchToken(responseBlob.SendToken)
	if err != nil {
		return swapResponse, err
	}

	receiveToken, err := blockchain.PatchToken(responseBlob.ReceiveToken)
	if err != nil {
		return swapResponse, err
	}

	sendTo, err := handler.wallet.GetAddress(sendToken.Blockchain)
	if err != nil {
		return swapResponse, err
	}

	receiveFrom, err := handler.wallet.GetAddress(receiveToken.Blockchain)
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

	blobSig, err := handler.Sign(blob.Password, blobBytes)
	if err != nil {
		return swapResponse, err
	}

	swapResponse.Swap = responseBlob
	swapResponse.Signature = base64.StdEncoding.EncodeToString(blobSig)

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
	return swapResponse, nil
}

func (handler *handler) Sign(password string, message []byte) ([]byte, error) {
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
