package zec

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/renproject/libzec-go"
	"github.com/renproject/swapperd/core/wallet/swapper/immediate"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/sirupsen/logrus"
)

type zecSwapContractBinder struct {
	scriptAddr string
	script     []byte
	swap       swap.Swap
	speed      blockchain.TxExecutionSpeed
	cost       blockchain.Cost
	logrus.FieldLogger
	libzec.Account
}

// NewZECSwapContractBinder returns a new Zcash Atom instance
func NewZECSwapContractBinder(account libzec.Account, swap swap.Swap, cost blockchain.Cost, logger logrus.FieldLogger) (immediate.Contract, error) {
	script, scriptAddr, err := buildInitiateScript(swap, account.NetworkParams())
	if err != nil {
		return nil, err
	}

	fields := logrus.Fields{}
	fields["SwapID"] = swap.ID
	fields["ContractID"] = scriptAddr
	fields["Token"] = swap.Token.Name
	logger = logger.WithFields(fields)

	if _, ok := cost[tokens.NameZEC]; !ok {
		cost[tokens.NameZEC] = big.NewInt(0)
	}

	if swap.BrokerFee.Int64() != 0 && swap.BrokerFee.Int64() < 600 {
		swap.BrokerFee = big.NewInt(600)
	}

	swap.Value = new(big.Int).Add(swap.Value, swap.BrokerFee)

	logger.Info(swap.ID, fmt.Sprintf("ZEC atomic swap = %s", scriptAddr))
	return &zecSwapContractBinder{
		scriptAddr:  scriptAddr,
		script:      script,
		swap:        swap,
		speed:       swap.Speed,
		FieldLogger: logger,
		Account:     account,
		cost:        cost,
	}, nil
}

// Initiate the atomic swap by funding a HTLC on the Zcash blockchain.
func (atom *zecSwapContractBinder) Initiate() error {
	atom.Info("Initiating on Zcash blockchain for ZEC")
	scriptAddr, err := libzec.DecodeAddress(atom.scriptAddr, atom.NetworkParams())
	if err != nil {
		return NewErrInitiate(err)
	}
	initiateScriptP2SHPKScript, err := libzec.PayToAddrScript(scriptAddr)
	if err != nil {
		return NewErrInitiate(NewErrBuildScript(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// signing a transaction with the given private key
	txHash, txFee, err := atom.SendTransaction(
		ctx,
		nil,
		libzec.Fast,
		nil,
		func(tx *wire.MsgTx) bool {
			// checks whether the contract is funded, with given value
			funded, value, err := atom.ScriptFunded(ctx, atom.scriptAddr, atom.swap.Value.Int64())
			if err != nil {
				return false
			}
			if funded {
				atom.Info(fmt.Sprintf("Send value on Zcash blockchain = %d", atom.swap.Value.Int64()))
				return false
			}

			if atom.swap.Value.Int64()-value < 600 {
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
			return funded
		},
		false,
	)
	if err != nil {
		if err == libzec.ErrPreConditionCheckFailed {
			return immediate.ErrAlreadyInitiated
		}
		return err
	}

	atom.cost[tokens.NameZEC] = new(big.Int).Add(big.NewInt(txFee), atom.cost[tokens.NameZEC])
	atom.cost[tokens.NameZEC] = new(big.Int).Add(atom.swap.BrokerFee, atom.cost[tokens.NameZEC])
	atom.Info(atom.FormatTransactionView("Initiated on Zcash blockchain", txHash))
	return nil
}

func (atom *zecSwapContractBinder) Audit() error {
	if funded, amount, err := atom.ScriptFunded(context.Background(), atom.scriptAddr, atom.swap.Value.Int64()); funded && err == nil {
		value := new(big.Int).Sub(atom.swap.Value, atom.swap.BrokerFee)
		if amount < value.Int64() {
			return fmt.Errorf("Audit Failed")
		}
		return nil
	}

	if time.Now().Unix() > atom.swap.TimeLock {
		return immediate.ErrSwapExpired
	}
	return immediate.ErrAuditPending
}

// Redeem the Atomic Swap by revealing the secret and withdrawing funds from the
// HTLC.
func (atom *zecSwapContractBinder) Redeem(secret [32]byte) error {
	atom.Info("Redeeming on Zcash blockchain")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	address, err := libzec.DecodeAddress(atom.swap.WithdrawAddress, atom.Account.NetworkParams())
	if err != nil {
		return NewErrRedeem(err)
	}

	payToAddrScript, err := libzec.PayToAddrScript(address)
	if err != nil {
		return NewErrRedeem(err)
	}

	var feeAddrScript []byte
	if atom.swap.BrokerFee.Int64() != 0 {
		feeAddress, err := libzec.DecodeAddress(atom.swap.BrokerAddress, atom.NetworkParams())
		if err != nil {
			return NewErrRedeem(err)
		}

		feeAddrScript, err = libzec.PayToAddrScript(feeAddress)
		if err != nil {
			return NewErrRedeem(err)
		}
	}

	txHash, txFee, err := atom.SendTransaction(
		ctx,
		atom.script,
		libzec.Fast,
		nil,
		func(tx *wire.MsgTx) bool {
			redeemed, val, err := atom.ScriptRedeemed(ctx, atom.scriptAddr, 0)
			if err != nil {
				return false
			}
			if !redeemed {
				if val-atom.swap.BrokerFee.Int64() < 600 {
					return false
				}
				if atom.swap.BrokerFee.Int64() >= 600 {
					tx.AddTxOut(wire.NewTxOut(atom.swap.BrokerFee.Int64(), feeAddrScript))
					tx.AddTxOut(wire.NewTxOut(val-atom.swap.BrokerFee.Int64(), payToAddrScript))
					return true
				}
				tx.AddTxOut(wire.NewTxOut(val, payToAddrScript))
				return true
			}
			return !redeemed
		},
		func(builder *txscript.ScriptBuilder) {
			builder.AddData(secret[:])
			builder.AddInt64(1)
		},
		func(tx *wire.MsgTx) bool {
			spent, _, err := atom.ScriptSpent(ctx, atom.scriptAddr, atom.swap.SpendingAddress)
			if err != nil {
				return false
			}
			return spent
		},
		true,
	)
	if err != nil {
		if err != libzec.ErrPreConditionCheckFailed {
			return nil
		}
		return err
	}
	atom.cost[tokens.NameZEC] = new(big.Int).Add(big.NewInt(txFee), atom.cost[tokens.NameZEC])
	atom.Info(atom.FormatTransactionView("Redeemed on Zcash blockchain", txHash))
	return nil
}

func (atom *zecSwapContractBinder) AuditSecret() ([32]byte, error) {
	atom.Info("Auditing secret on Zcash blockchain")
	spent, sigScript, err := atom.ScriptSpent(context.Background(), atom.scriptAddr, atom.swap.SpendingAddress)
	if !spent || err != nil {
		if time.Now().Unix() > atom.swap.TimeLock {
			return [32]byte{}, immediate.ErrSwapExpired
		}
		return [32]byte{}, immediate.ErrAuditPending
	}
	sigScriptBytes, err := hex.DecodeString(sigScript)
	if err != nil {
		return [32]byte{}, NewErrAuditSecret(err)
	}

	pushes, err := txscript.PushedData(sigScriptBytes)
	if err != nil {
		return [32]byte{}, NewErrAuditSecret(err)
	}
	for _, push := range pushes {
		if sha256.Sum256(push) == atom.swap.SecretHash {
			var secret [32]byte
			copy(secret[:], push)
			atom.Info(fmt.Sprintf("Audit succeeded on Zcash blockchain secret = %s", base64.StdEncoding.EncodeToString(secret[:])))
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// Refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *zecSwapContractBinder) Refund() error {
	atom.Info("Refunding on Zcash blockchain")
	address, err := atom.Address()
	if err != nil {
		return NewErrRefund(err)
	}
	payToAddrScript, err := libzec.PayToAddrScript(address)
	if err != nil {
		return NewErrRefund(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	txHash, txFee, err := atom.SendTransaction(
		ctx,
		atom.script,
		libzec.Fast,
		func(txIn *wire.TxIn) {
			txIn.Sequence = 0
		},
		func(tx *wire.MsgTx) bool {
			funded, val, err := atom.ScriptFunded(ctx, atom.scriptAddr, 0)
			if err != nil {
				return false
			}
			if funded {
				tx.AddTxOut(wire.NewTxOut(val, payToAddrScript))
			}
			tx.LockTime = uint32(atom.swap.TimeLock)
			return funded
		},
		func(builder *txscript.ScriptBuilder) {
			builder.AddInt64(0)
		},
		func(tx *wire.MsgTx) bool {
			spent, _, err := atom.ScriptSpent(ctx, atom.scriptAddr, atom.swap.SpendingAddress)
			if err != nil {
				return false
			}
			return spent
		},
		true,
	)

	if err != nil {
		if err != libzec.ErrPreConditionCheckFailed {
			return NewErrRefund(err)
		}
		return nil
	}
	atom.cost[tokens.NameZEC] = new(big.Int).Add(big.NewInt(txFee), atom.cost[tokens.NameZEC])
	atom.cost[tokens.NameZEC] = new(big.Int).Sub(atom.cost[tokens.NameZEC], atom.swap.BrokerFee)
	atom.Info(atom.FormatTransactionView("Refunded on Zcash blockchain", txHash))
	return nil
}

func (atom *zecSwapContractBinder) Cost() blockchain.Cost {
	return atom.cost
}
