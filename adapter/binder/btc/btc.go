package btc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
)

type btcSwapContractBinder struct {
	scriptAddr string
	script     []byte
	swap       swap.Swap
	txVersion  int32
	fee        int64
	verify     bool
	cost       blockchain.Cost
	logrus.FieldLogger
	libbtc.Account
}

// NewBTCSwapContractBinder returns a new Bitcoin Atom instance
func NewBTCSwapContractBinder(account libbtc.Account, swap swap.Swap, cost blockchain.Cost, logger logrus.FieldLogger) (swapper.Contract, error) {
	script, scriptAddr, err := buildInitiateScript(swap, account.NetworkParams())
	if err != nil {
		return nil, err
	}

	fields := logrus.Fields{}
	fields["SwapID"] = swap.ID
	fields["ContractID"] = scriptAddr
	fields["Token"] = swap.Token.Name
	logger = logger.WithFields(fields)

	if _, ok := cost[blockchain.BTC]; !ok {
		cost[blockchain.BTC] = big.NewInt(0)
	}

	logger.Info(swap.ID, fmt.Sprintf("BTC atomic swap = %s", scriptAddr))
	return &btcSwapContractBinder{
		scriptAddr:  scriptAddr,
		script:      script,
		swap:        swap,
		txVersion:   2,
		fee:         swap.Fee.Int64(),
		verify:      true,
		FieldLogger: logger,
		Account:     account,
		cost:        cost,
	}, nil
}

// Initiate the atomic swap by funding a HTLC on the Bitcoin blockchain.
func (atom *btcSwapContractBinder) Initiate() error {
	atom.Info("Initiating on Bitcoin blockchain for BTC")
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
	depositValue := atom.swap.Value.Int64() + atom.swap.BrokerFee.Int64()

	// signing a transaction with the given private key
	if err := atom.SendTransaction(
		ctx,
		nil,
		atom.fee,
		func(tx *wire.MsgTx) bool {
			// checks whether the contract is funded, with given value
			funded, value, err := atom.ScriptFunded(ctx, atom.scriptAddr, depositValue)
			if err != nil {
				return false
			}
			if funded {
				atom.Info(fmt.Sprintf("Send value on Bitcoin blockchain = %d", depositValue))
				return false
			}
			// creating unsigned transaction and adding transaction outputs
			tx.AddTxOut(wire.NewTxOut(depositValue-value, initiateScriptP2SHPKScript))
			return !funded
		},
		nil,
		func(tx *wire.MsgTx) bool {
			funded, _, err := atom.ScriptFunded(ctx, atom.scriptAddr, depositValue)
			if err != nil {
				return false
			}
			if funded {
				atom.Info(atom.FormatTransactionView("Initiated on Bitcoin blockchain", tx.TxHash().String()))
			}
			return funded
		},
	); err != nil && err != libbtc.ErrPreConditionCheckFailed {
		return err
	}
	atom.cost[blockchain.BTC] = new(big.Int).Add(big.NewInt(atom.fee), atom.cost[blockchain.BTC])
	return nil
}

func (atom *btcSwapContractBinder) Audit() error {
	depositValue := atom.swap.Value.Int64() + atom.swap.BrokerFee.Int64()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	for {
		if funded, _, err := atom.ScriptFunded(ctx, atom.scriptAddr, depositValue); funded && err == nil {
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
	atom.Info("Redeeming on Bitcoin blockchain")
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

	feeAddress, err := btcutil.DecodeAddress(atom.swap.BrokerAddress, atom.NetworkParams())
	if err != nil {
		return NewErrRedeem(err)
	}

	feeAddrScript, err := txscript.PayToAddrScript(feeAddress)
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
				tx.AddTxOut(wire.NewTxOut(atom.swap.BrokerFee.Int64(), feeAddrScript))
				tx.AddTxOut(wire.NewTxOut(val-atom.swap.BrokerFee.Int64()-atom.fee, payToAddrScript))
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
				atom.Info(atom.FormatTransactionView("Redeemed on Bitcoin blockchain", tx.TxHash().String()))
			}
			if err != nil {
				return false
			}
			return spent
		},
	); err != nil && err != libbtc.ErrPreConditionCheckFailed {
		return err
	}
	atom.cost[blockchain.BTC] = new(big.Int).Add(big.NewInt(atom.fee), atom.cost[blockchain.BTC])
	return nil
}

func (atom *btcSwapContractBinder) AuditSecret() ([32]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	for {
		atom.Info("Auditing secret on Bitcoin blockchain")
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
			atom.Info(fmt.Sprintf("Audit succeeded on Bitcoin blockchain secret = %s", base64.StdEncoding.EncodeToString(secret[:])))
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// Refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *btcSwapContractBinder) Refund() error {
	atom.Info("Refunding on Bitcoin blockchain")
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
				atom.Info(atom.FormatTransactionView("Refunded on Bitcoin blockchain", tx.TxHash().String()))
			}
			return spent
		},
	); err != nil && err != libbtc.ErrPreConditionCheckFailed {
		return err
	}
	atom.cost[blockchain.BTC] = new(big.Int).Add(big.NewInt(atom.fee), atom.cost[blockchain.BTC])
	return nil
}

func (atom *btcSwapContractBinder) Cost() blockchain.Cost {
	return atom.cost
}
