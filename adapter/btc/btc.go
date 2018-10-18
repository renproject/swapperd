package btc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/core"
	"github.com/republicprotocol/swapperd/foundation"
)

type btcSwapContractBinder struct {
	scriptAddr string
	script     []byte
	req        foundation.Swap
	txVersion  int32
	fee        int64
	verify     bool
	ctx        context.Context
	core.Logger
	libbtc.Account
}

// NewBTCSwapContractBinder returns a new Bitcoin Atom instance
func NewBTCSwapContractBinder(account libbtc.Account, logger core.Logger, req foundation.Swap) (core.SwapContractBinder, error) {
	address, err := account.Address()
	if err != nil {
		return nil, err
	}
	script, scriptAddr, err := buildInitiateScript(address.EncodeAddress(), req, account.NetworkParams())
	if err != nil {
		return nil, err
	}
	logger.LogInfo(req.ID, fmt.Sprintf("Bitcoin atomic swap id: %s", scriptAddr))
	return &btcSwapContractBinder{
		scriptAddr: scriptAddr,
		script:     script,
		req:        req,
		txVersion:  2,
		fee:        10000,
		verify:     true,
		ctx:        context.Background(),
		Logger:     logger,
		Account:    account,
	}, nil
}

// Initiate the atomic swap by funding a HTLC on the Bitcoin blockchain.
func (atom *btcSwapContractBinder) Initiate() error {
	atom.LogInfo(atom.req.ID, "initiating on Bitcoin blockchain")
	scriptAddr, err := btcutil.DecodeAddress(atom.scriptAddr, atom.NetworkParams())
	if err != nil {
		return NewErrInitiate(err)
	}
	sendValue, err := strconv.ParseInt(hex.EncodeToString(atom.req.SendValue[:]), 16, 64)
	if err != nil {
		return NewErrInitiate(err)
	}
	if sendValue == 0 {
		return NewErrInitiate(fmt.Errorf("trying to send 0 Bitcoins"))
	}
	initiateScriptP2SHPKScript, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return NewErrInitiate(NewErrBuildScript(err))
	}

	// signing a transaction with the given private key
	return atom.SendTransaction(
		atom.ctx,
		nil,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			// checks whether the contract is funded, with given value
			funded, value, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, sendValue)
			if err != nil {
				return false
			}
			if funded {
				atom.LogInfo(atom.req.ID, fmt.Sprintf("Bitcoin swap initiated with send value %d", sendValue))
				return false
			}
			// creating unsigned transaction and adding transaction outputs
			tx.AddTxOut(wire.NewTxOut(sendValue-value, initiateScriptP2SHPKScript))
			return !funded
		},
		nil,
		func(tx *wire.MsgTx) bool {
			funded, _, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, sendValue)
			if err != nil {
				return false
			}
			if funded {
				atom.LogInfo(atom.req.ID, atom.FormatTransactionView("initiated the atomic swap on Bitcoin blockchain", tx.TxHash().String()))
			}
			return funded
		},
	)
}

func (atom *btcSwapContractBinder) Audit() error {
	receiveValue, err := strconv.ParseInt(hex.EncodeToString(atom.req.ReceiveValue[:]), 16, 64)
	if err != nil {
		return err
	}

	for {
		if funded, _, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, receiveValue); funded && err == nil {
			return nil
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return NewErrAudit(ErrTimedOut)
		}
		time.Sleep(15 * time.Second)
	}
}

// Redeem the Atomic Swap by revealing the secret and withdrawing funds from the
// HTLC.
func (atom *btcSwapContractBinder) Redeem(secret [32]byte) error {
	atom.LogInfo(atom.req.ID, "redeeming on Bitcoin blockchain")
	address, err := atom.Address()
	if err != nil {
		return NewErrRedeem(err)
	}
	payToAddrScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return NewErrRedeem(err)
	}
	return atom.SendTransaction(
		atom.ctx,
		atom.script,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			funded, val, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, 0)
			if err != nil {
				return false
			}
			if funded {
				tx.AddTxOut(wire.NewTxOut(val-atom.fee, payToAddrScript))
			}
			return funded
		},
		func(builder *txscript.ScriptBuilder) {
			builder.AddData(secret[:])
			builder.AddInt64(1)
		},
		func(tx *wire.MsgTx) bool {
			spent, err := atom.ScriptSpent(atom.ctx, atom.scriptAddr)
			if spent {
				atom.LogInfo(atom.req.ID, atom.FormatTransactionView("redeemed the atomic swap on Bitcoin blockchain", tx.TxHash().String()))
			}
			if err != nil {
				return false
			}
			return spent
		},
	)
}

func (atom *btcSwapContractBinder) AuditSecret() ([32]byte, error) {
	for {
		if spent, err := atom.ScriptSpent(atom.ctx, atom.scriptAddr); spent && err == nil {
			break
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return [32]byte{}, NewErrAuditSecret(ErrTimedOut)
		}
		time.Sleep(15 * time.Second)
	}

	sigScript, err := atom.GetScriptFromSpentP2SH(atom.ctx, atom.scriptAddr)
	if err != nil {
		return [32]byte{}, NewErrAuditSecret(err)
	}

	pushes, err := txscript.PushedData(sigScript)
	if err != nil {
		return [32]byte{}, NewErrAuditSecret(err)
	}
	for _, push := range pushes {
		if sha256.Sum256(push) == atom.req.SecretHash {
			var secret [32]byte
			copy(secret[:], push)
			atom.LogInfo(atom.req.ID, fmt.Sprintf("audit secret successful on Bitcoin blockchain %s", base64.StdEncoding.EncodeToString(secret[:])))
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// Refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *btcSwapContractBinder) Refund() error {
	atom.LogInfo(atom.req.ID, "refunding on Bitcoin blockchain")
	address, err := atom.Address()
	if err != nil {
		return NewErrRedeem(err)
	}
	payToAddrScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return NewErrRedeem(err)
	}
	return atom.SendTransaction(
		atom.ctx,
		atom.script,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			funded, val, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, 0)
			if err != nil {
				return false
			}
			if funded {
				tx.AddTxOut(wire.NewTxOut(val-atom.fee, payToAddrScript))
			}
			return funded
		},
		func(builder *txscript.ScriptBuilder) {
			builder.AddInt64(0)
		},
		func(tx *wire.MsgTx) bool {
			spent, err := atom.ScriptSpent(atom.ctx, atom.scriptAddr)
			if err != nil {
				return false
			}
			if spent {
				atom.LogInfo(atom.req.ID, atom.FormatTransactionView("refunded the atomic swap on Bitcoin blockchain", tx.TxHash().String()))
			}
			return spent
		},
	)
}
