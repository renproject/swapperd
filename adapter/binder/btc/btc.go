package btc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type btcSwapContractBinder struct {
	scriptAddr string
	script     []byte
	swap       foundation.SwapTry
	txVersion  int32
	fee        int64
	verify     bool
	ctx        context.Context
	swapper.Logger
	libbtc.Account
}

// NewBTCSwapContractBinder returns a new Bitcoin Atom instance
func NewBTCSwapContractBinder(account libbtc.Account, swap foundation.SwapTry, logger swapper.Logger) (swapper.SwapContractBinder, error) {
	script, scriptAddr, err := buildInitiateScript(swap, account.NetworkParams())
	if err != nil {
		return nil, err
	}

	logger.LogInfo(swap.ID, fmt.Sprintf("Bitcoin atomic swap id: %s", scriptAddr))
	return &btcSwapContractBinder{
		scriptAddr: scriptAddr,
		script:     script,
		swap:       swap,
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
	atom.LogInfo(atom.swap.ID, "initiating on Bitcoin blockchain")
	scriptAddr, err := btcutil.DecodeAddress(atom.scriptAddr, atom.NetworkParams())
	if err != nil {
		return NewErrInitiate(err)
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
			funded, value, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, atom.swap.Value.Int64())
			if err != nil {
				return false
			}
			if funded {
				atom.LogInfo(atom.swap.ID, fmt.Sprintf("Bitcoin swap initiated with send value %d", atom.swap.Value.Int64()))
				return false
			}
			// creating unsigned transaction and adding transaction outputs
			tx.AddTxOut(wire.NewTxOut(atom.swap.Value.Int64()-value, initiateScriptP2SHPKScript))
			return !funded
		},
		nil,
		func(tx *wire.MsgTx) bool {
			funded, _, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, atom.swap.Value.Int64())
			if err != nil {
				return false
			}
			if funded {
				atom.LogInfo(atom.swap.ID, atom.FormatTransactionView("initiated the atomic swap on Bitcoin blockchain", tx.TxHash().String()))
			}
			return funded
		},
	)
}

func (atom *btcSwapContractBinder) Audit() error {
	for {
		if funded, _, err := atom.ScriptFunded(atom.ctx, atom.scriptAddr, atom.swap.Value.Int64()); funded && err == nil {
			return nil
		}
		if time.Now().Unix() > atom.swap.TimeLock {
			return NewErrAudit(ErrTimedOut)
		}
		time.Sleep(15 * time.Second)
	}
}

// Redeem the Atomic Swap by revealing the secret and withdrawing funds from the
// HTLC.
func (atom *btcSwapContractBinder) Redeem(secret [32]byte) error {
	atom.LogInfo(atom.swap.ID, "redeeming on Bitcoin blockchain")
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
				atom.LogInfo(atom.swap.ID, atom.FormatTransactionView("redeemed the atomic swap on Bitcoin blockchain", tx.TxHash().String()))
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
		if time.Now().Unix() > atom.swap.TimeLock {
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
		if sha256.Sum256(push) == atom.swap.SecretHash {
			var secret [32]byte
			copy(secret[:], push)
			atom.LogInfo(atom.swap.ID, fmt.Sprintf("audit secret successful on Bitcoin blockchain %s", base64.StdEncoding.EncodeToString(secret[:])))
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// Refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *btcSwapContractBinder) Refund() error {
	atom.LogInfo(atom.swap.ID, "refunding on Bitcoin blockchain")
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
				atom.LogInfo(atom.swap.ID, atom.FormatTransactionView("refunded the atomic swap on Bitcoin blockchain", tx.TxHash().String()))
			}
			return spent
		},
	)
}
