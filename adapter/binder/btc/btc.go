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
	swap       foundation.Swap
	txVersion  int32
	fee        int64
	verify     bool
	swapper.Logger
	libbtc.Account
}

// NewBTCSwapContractBinder returns a new Bitcoin Atom instance
func NewBTCSwapContractBinder(account libbtc.Account, swap foundation.Swap, logger swapper.Logger) (swapper.Contract, error) {
	script, scriptAddr, err := buildInitiateScript(swap, account.NetworkParams())
	if err != nil {
		return nil, err
	}

	logger.LogInfo(swap.ID, fmt.Sprintf("BTC atomic swap = %s", scriptAddr))
	return &btcSwapContractBinder{
		scriptAddr: scriptAddr,
		script:     script,
		swap:       swap,
		txVersion:  2,
		fee:        10000,
		verify:     true,
		Logger:     logger,
		Account:    account,
	}, nil
}

// Initiate the atomic swap by funding a HTLC on the Bitcoin blockchain.
func (atom *btcSwapContractBinder) Initiate() error {
	atom.LogInfo(atom.swap.ID, "Initiating on Bitcoin blockchain for BTC")
	scriptAddr, err := btcutil.DecodeAddress(atom.scriptAddr, atom.NetworkParams())
	if err != nil {
		return NewErrInitiate(err)
	}
	initiateScriptP2SHPKScript, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return NewErrInitiate(NewErrBuildScript(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// signing a transaction with the given private key
	if err := atom.SendTransaction(
		ctx,
		nil,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			// checks whether the contract is funded, with given value
			funded, value, err := atom.ScriptFunded(ctx, atom.scriptAddr, atom.swap.Value.Int64())
			if err != nil {
				return false
			}
			if funded {
				atom.LogInfo(atom.swap.ID, fmt.Sprintf("Send value on Bitcoin blockchain = %d", atom.swap.Value.Int64()))
				return false
			}
			// creating unsigned transaction and adding transaction outputs
			tx.AddTxOut(wire.NewTxOut(atom.swap.Value.Int64()-value, initiateScriptP2SHPKScript))
			return !funded
		},
		nil,
		func(tx *wire.MsgTx) bool {
			funded, _, err := atom.ScriptFunded(ctx, atom.scriptAddr, atom.swap.Value.Int64())
			if err != nil {
				return false
			}
			if funded {
				atom.LogInfo(atom.swap.ID, atom.FormatTransactionView("Initiated on Bitcoin blockchain", tx.TxHash().String()))
			}
			return funded
		},
	); err != nil && err != libbtc.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

func (atom *btcSwapContractBinder) Audit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	for {
		if funded, _, err := atom.ScriptFunded(ctx, atom.scriptAddr, atom.swap.Value.Int64()); funded && err == nil {
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
	atom.LogInfo(atom.swap.ID, "Redeeming on Bitcoin blockchain")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	address, err := atom.Address()
	if err != nil {
		return NewErrRedeem(err)
	}
	payToAddrScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return NewErrRedeem(err)
	}
	if err := atom.SendTransaction(
		ctx,
		atom.script,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			funded, val, err := atom.ScriptFunded(ctx, atom.scriptAddr, 0)
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
			spent, err := atom.ScriptSpent(ctx, atom.scriptAddr)
			if spent {
				atom.LogInfo(atom.swap.ID, atom.FormatTransactionView("Redeemed on Bitcoin blockchain", tx.TxHash().String()))
			}
			if err != nil {
				return false
			}
			return spent
		},
	); err != nil && err != libbtc.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

func (atom *btcSwapContractBinder) AuditSecret() ([32]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	for {
		atom.LogInfo(atom.swap.ID, "Auditing secret on Bitcoin blockchain")
		if spent, err := atom.ScriptSpent(ctx, atom.scriptAddr); spent && err == nil {
			break
		}
		if time.Now().Unix() > atom.swap.TimeLock {
			return [32]byte{}, NewErrAuditSecret(ErrTimedOut)
		}
		time.Sleep(time.Minute)
	}

	sigScript, err := atom.GetScriptFromSpentP2SH(ctx, atom.scriptAddr)
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
			atom.LogInfo(atom.swap.ID, fmt.Sprintf("Audit succeeded on Bitcoin blockchain secret = %s", base64.StdEncoding.EncodeToString(secret[:])))
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// Refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *btcSwapContractBinder) Refund() error {
	atom.LogInfo(atom.swap.ID, "Refunding on Bitcoin blockchain")
	address, err := atom.Address()
	if err != nil {
		return NewErrRedeem(err)
	}
	payToAddrScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return NewErrRedeem(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := atom.SendTransaction(
		ctx,
		atom.script,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			funded, val, err := atom.ScriptFunded(ctx, atom.scriptAddr, 0)
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
			spent, err := atom.ScriptSpent(ctx, atom.scriptAddr)
			if err != nil {
				return false
			}
			if spent {
				atom.LogInfo(atom.swap.ID, atom.FormatTransactionView("Refunded on Bitcoin blockchain", tx.TxHash().String()))
			}
			return spent
		},
	); err != nil && err != libbtc.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}
