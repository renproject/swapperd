package zec

import (
	"encoding/hex"
	"fmt"
)

var ErrCompleteSignTransaction = NewErrSignTransaction(fmt.Errorf("incomplete signature"))
var ErrContractOutput = fmt.Errorf("transaction does not contain a contract output")
var ErrInitiated = fmt.Errorf("atomic swap already initiated")
var ErrMalformedRedeemTx = fmt.Errorf("redeem transaction returned by the ZCash blockchain is malformed")
var ErrMalformedInitiateTx = fmt.Errorf("initiate transaction returned by the ZCash blockchain is malformed")
var ErrUnknownMessageType = fmt.Errorf("unknown message type")
var ErrTimedOut = fmt.Errorf("timed out")

func NewErrDecodeAddress(addr string, err error) error {
	return fmt.Errorf("failed to decode address (%s): %v", addr, err)
}

func NewErrDecodeScript(script []byte, err error) error {
	return fmt.Errorf("failed to decode script (%s): %v", script, err)
}

func NewErrSignTransaction(err error) error {
	return fmt.Errorf("failed to sign transaction: %v", err)
}

func NewErrPublishTransaction(err error) error {
	return fmt.Errorf("failed to publish signed transaction: %v", err)
}

func NewErrBuildScript(err error) error {
	return fmt.Errorf("failed to build zcash script: %v", err)
}

func NewErrDecodeTransaction(txBytes []byte, err error) error {
	return fmt.Errorf("failed to decode contract transaction: %s %v", hex.EncodeToString(txBytes), err)
}

func NewErrScriptExec(err error) error {
	return fmt.Errorf("script execution error: %v", err)
}

func NewErrInitiate(err error) error {
	return fmt.Errorf("failed to initiate: %v", err)
}

func NewErrRedeem(err error) error {
	return fmt.Errorf("failed to redeem: %v", err)
}

func NewErrRefund(err error) error {
	return fmt.Errorf("failed to refund: %v", err)
}

func NewErrAudit(err error) error {
	return fmt.Errorf("failed to audit: %v", err)
}

func NewErrAuditSecret(err error) error {
	return fmt.Errorf("failed to audit secret: %v", err)
}
