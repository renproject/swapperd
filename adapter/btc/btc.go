package btc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type Conn interface {
	// SignTransaction should sign the transaction using the given bitcoin key.
	SignTransaction(tx *wire.MsgTx, key keystore.BitcoinKey, fee int64) (*wire.MsgTx, bool, error)

	// PublishTransaction should publish a signed transaction to the Bitcoin
	// blockchain.
	PublishTransaction(signedTransaction *wire.MsgTx, postCon func() bool) error

	// Net should return the network information of the underlying
	// Bitcoin blockchain.
	Net() *chaincfg.Params

	// GetScriptFromSpentP2SH for the given script address.
	GetScriptFromSpentP2SH(address string) ([]byte, error)

	// Balance of the given address on Bitcoin blockchain.
	Balance(address string, confirmations int64) int64

	// Withdraw given value or the whole balance to the given address.
	Withdraw(addr string, key keystore.BitcoinKey, value, fee int64) error

	// SpendBalance of the given address on Bitcoin blockchain.
	SpendBalance(address string) (*wire.MsgTx, []byte, []int64, error)

	// ScriptFunded checks whether a script is funded.
	ScriptFunded(address string, value int64) (bool, int64)

	// ScriptSpent checks whether a script is spent.
	ScriptSpent(address string) bool
}

type bitcoinAtom struct {
	scriptAddr string
	script     []byte
	key        keystore.BitcoinKey
	req        swapDomain.Request
	txVersion  int32
	fee        int64
	verify     bool
	logger.Logger
	Conn
}

// NewBitcoinAtom returns a new Bitcoin Atom instance
func NewBitcoinAtom(conf config.BitcoinNetwork, key keystore.BitcoinKey, logger logger.Logger, req swapDomain.Request) (swap.Atom, error) {
	conn := NewConnWithConfig(conf)

	script, scriptAddr, err := buildInitiateScript(key.AddressString, req, conn.Net())
	if err != nil {
		return nil, err
	}
	logger.LogInfo(req.UID, fmt.Sprintf("Bitcoin Atomic Swap ID: %s", scriptAddr))
	return &bitcoinAtom{
		scriptAddr: scriptAddr,
		script:     script,
		key:        key,
		req:        req,
		txVersion:  2,
		fee:        10000,
		verify:     true,
		Logger:     logger,
		Conn:       conn,
	}, nil
}

// initiate the atomic swap by funding a HTLC on the Bitcoin blockchain.
func (atom *bitcoinAtom) Initiate() error {
	atom.LogInfo(atom.req.UID, "Initiating on bitcoin blockchain")
	scriptAddr, err := btcutil.DecodeAddress(atom.scriptAddr, atom.Net())
	if err != nil {
		return err
	}

	sendValue, err := strconv.ParseInt(atom.req.SendValue, 10, 64)
	if err != nil {
		return err
	}

	if sendValue == 0 {
		atom.LogError(atom.req.UID, "Trying to send 0 Bitcoins")
	}

	if funded, value := atom.ScriptFunded(atom.scriptAddr, sendValue); value > 0 || funded {
		if funded {
			atom.LogDebug(atom.req.UID, fmt.Sprintf("Bitcoin swap initiated with send value %d", sendValue))
			return swap.ErrSwapAlreadyInitiated
		}
		sendValue = sendValue - value
	}

	initiateScriptP2SHPkScript, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return NewErrInitiate(NewErrBuildScript(err))
	}

	// creating unsigned transaction and adding transaction outputs
	unsignedTx := wire.NewMsgTx(atom.txVersion)
	unsignedTx.AddTxOut(wire.NewTxOut(sendValue, initiateScriptP2SHPkScript))

	// signing a transaction with the given private key
	stx, complete, err := atom.Conn.SignTransaction(unsignedTx, atom.key, atom.fee)
	if err != nil {
		return NewErrInitiate(NewErrSignTransaction(err))
	}
	if !complete {
		return NewErrInitiate(ErrCompleteSignTransaction)
	}

	if err := atom.Conn.PublishTransaction(stx,
		func() bool {
			success, _ := atom.ScriptFunded(atom.scriptAddr, sendValue)
			return success
		},
	); err != nil {
		return err
	}

	atom.LogInfo(atom.req.UID, "Initiated on bitcoin blockchain")
	return nil
}

func (atom *bitcoinAtom) Audit() error {
	receiveValue, err := strconv.ParseInt(atom.req.ReceiveValue, 10, 64)
	if err != nil {
		return err
	}

	for {
		if funded, _ := atom.ScriptFunded(atom.scriptAddr, receiveValue); funded {
			return nil
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return NewErrAudit(ErrTimedOut)
		}
		time.Sleep(15 * time.Second)
	}
}

// redeem the Atomic Swap by revealing the secret and withdrawing funds from the
// HTLC.
func (atom *bitcoinAtom) Redeem(secret [32]byte) error {
	atom.LogInfo(atom.req.UID, "Redeeming on bitcoin blockchain")
	if spent := atom.ScriptSpent(atom.scriptAddr); spent {
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}

	redeemTx, scriptPubKey, inputValues, err := atom.SpendBalance(atom.scriptAddr)
	if err != nil {
		return NewErrRedeem(err)
	}

	receiveValue, err := strconv.ParseInt(atom.req.ReceiveValue, 10, 64)
	if err != nil {
		return NewErrRedeem(err)
	}

	// create bitcoin script to pay to the user's personal address
	payToAddrScript, err := txscript.PayToAddrScript(atom.key.Address)
	if err != nil {
		return NewErrRedeem(err)
	}

	redeemTx.AddTxOut(wire.NewTxOut(receiveValue-atom.fee, payToAddrScript))

	for i, txIn := range redeemTx.TxIn {
		// sign transaction
		redeemSig, redeemPubKey, err := sign(redeemTx, int(txIn.PreviousOutPoint.Index), atom.script, atom.key)
		if err != nil {
			return NewErrRedeem(err)
		}

		// build signature script
		redeemSigScript, err := newRedeemScript(atom.script, redeemSig, redeemPubKey, secret)
		if err != nil {
			return NewErrRedeem(err)
		}

		txIn.SignatureScript = redeemSigScript

		verifyTransaction(scriptPubKey, redeemTx, int(txIn.PreviousOutPoint.Index), inputValues[i])
	}

	if err := atom.PublishTransaction(redeemTx,
		func() bool {
			return atom.ScriptSpent(atom.scriptAddr)
		},
	); err != nil {
		return err
	}

	atom.LogInfo(atom.req.UID, "Redeemed on bitcoin blockchain")
	return nil
}

func (atom *bitcoinAtom) AuditSecret() ([32]byte, error) {
	for {
		if spent := atom.ScriptSpent(atom.scriptAddr); spent {
			break
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return [32]byte{}, NewErrAuditSecret(ErrTimedOut)
		}
		time.Sleep(15 * time.Second)
	}

	sigScript, err := atom.Conn.GetScriptFromSpentP2SH(atom.scriptAddr)
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
			atom.LogInfo(atom.req.UID, fmt.Sprintf("Audit secret successful on Bitcoin blockchain %s", base64.StdEncoding.EncodeToString(secret[:])))
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *bitcoinAtom) Refund() error {
	atom.LogInfo(atom.req.UID, "Refunding on bitcoin blockchain")
	if spent := atom.ScriptSpent(atom.scriptAddr); spent {
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}

	refundTx, scriptPubKey, inputValues, err := atom.SpendBalance(atom.scriptAddr)
	if err != nil {
		return NewErrRefund(err)
	}

	receiveValue, err := strconv.ParseInt(atom.req.ReceiveValue, 10, 64)
	if err != nil {
		return NewErrRefund(err)
	}

	// create bitcoin script to pay to the user's personal address
	payToAddrScript, err := txscript.PayToAddrScript(atom.key.Address)
	if err != nil {
		return NewErrRefund(err)
	}

	refundTx.AddTxOut(wire.NewTxOut(receiveValue-atom.fee, payToAddrScript))

	for i, txIn := range refundTx.TxIn {
		// sign transaction
		refundSig, refundPubKey, err := sign(refundTx, int(txIn.PreviousOutPoint.Index), atom.script, atom.key)
		if err != nil {
			return NewErrRefund(err)
		}

		// build signature script
		refundSigScript, err := newRefundScript(atom.script, refundSig, refundPubKey)
		if err != nil {
			return NewErrRefund(err)
		}

		txIn.SignatureScript = refundSigScript

		verifyTransaction(scriptPubKey, refundTx, int(txIn.PreviousOutPoint.Index), inputValues[i])
	}

	if err := atom.PublishTransaction(refundTx,
		func() bool {
			return atom.ScriptSpent(atom.scriptAddr)
		},
	); err != nil {
		return err
	}

	atom.LogInfo(atom.req.UID, "Refunded on bitcoin blockchain")
	return nil
}
