package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

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
	PublishTransaction(signedTransaction []byte, postCon func() (bool, error)) error

	// Net should return the network information of the underlying
	// Bitcoin blockchain.
	Net() *chaincfg.Params

	// GetScriptFromSpentP2SH for the given script address.
	GetScriptFromSpentP2SH(address string) ([]byte, error)

	// Balance of the given address on Bitcoin blockchain.
	Balance(address string) (int64, error)

	// GetUnspentOutputs of the given address on Bitcoin blockchain.
	GetUnspentOutputs(address string) (UnspentOutputs, error)

	ScriptFunded(address string, value int64) (bool, int64, error)
	ScriptSpent(address string) (bool, error)
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
	conn, err := NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}
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

	if funded, value, err := atom.ScriptFunded(atom.scriptAddr, sendValue); value > 0 || funded || err != nil {
		if err != nil {
			return err
		}
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

	// marshal signed transaction information
	var initiateTxBuffer bytes.Buffer
	initiateTxBuffer.Grow(stx.SerializeSize())
	if err := stx.Serialize(&initiateTxBuffer); err != nil {
		return NewErrInitiate(err)
	}

	if err := atom.Conn.PublishTransaction(initiateTxBuffer.Bytes(),
		func() (bool, error) {
			success, _, err := atom.ScriptFunded(atom.scriptAddr, sendValue)
			return success, err
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
		if funded, _, err := atom.ScriptFunded(atom.scriptAddr, receiveValue); funded || err != nil {
			return err
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
	output := UnspentOutput{}
	receiveValue, err := strconv.ParseInt(atom.req.ReceiveValue, 10, 64)
	if err != nil {
		return err
	}
	for {
		outs, err := atom.Conn.GetUnspentOutputs(atom.scriptAddr)
		if err != nil {
			return NewErrRedeem(err)
		}
		if len(outs.Outputs) != 0 {
			output = outs.Outputs[0]
			break
		}
	}

	if spent, err := atom.ScriptSpent(atom.scriptAddr); spent || err != nil {
		if err != nil {
			return err
		}
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}

	// create bitcoin script to pay to the user's personal address
	payToAddrScript, err := txscript.PayToAddrScript(atom.key.Address)
	if err != nil {
		return NewErrRedeem(err)
	}

	// build transaction
	hashBytes, err := hex.DecodeString(output.TxHash)
	if err != nil {
		return NewErrRedeem(err)
	}
	txHash, err := chainhash.NewHash(hashBytes)
	if err != nil {
		return NewErrRedeem(err)
	}
	redeemTx := wire.NewMsgTx(atom.txVersion)
	redeemTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txHash, output.Vout), nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(receiveValue-atom.fee, payToAddrScript))

	// sign transaction
	redeemSig, redeemPubKey, err := sign(redeemTx, int(output.Vout), atom.script, atom.key)
	if err != nil {
		return NewErrRedeem(err)
	}

	// build signature script
	redeemSigScript, err := newRedeemScript(atom.script, redeemSig, redeemPubKey, secret)
	if err != nil {
		return NewErrRedeem(err)
	}
	redeemTx.TxIn[output.Vout].SignatureScript = redeemSigScript

	if atom.verify {
		// verifying the redeem script
		scriptPubKey, err := hex.DecodeString(output.ScriptPubKey)
		if err != nil {
			return NewErrRefund(err)
		}
		e, err := txscript.NewEngine(scriptPubKey, redeemTx, int(output.Vout),
			txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), receiveValue)
		if err != nil {
			return NewErrRedeem(err)
		}
		err = e.Execute()
		if err != nil {
			return NewErrRedeem(NewErrScriptExec(err))
		}
	}

	// marshal signed transaction information
	var redeemTxBuffer bytes.Buffer
	redeemTxBuffer.Grow(redeemTx.SerializeSize())
	if err := redeemTx.Serialize(&redeemTxBuffer); err != nil {
		return NewErrRedeem(err)
	}

	if err := atom.PublishTransaction(redeemTxBuffer.Bytes(),
		func() (bool, error) {
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
		if spent, err := atom.ScriptSpent(atom.scriptAddr); spent || err != nil {
			if err != nil {
				return [32]byte{}, err
			}
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
	sendValue, err := strconv.ParseInt(atom.req.SendValue, 10, 64)
	if err != nil {
		return err
	}
	// TODO: Use all the outputs
	outs, err := atom.Conn.GetUnspentOutputs(atom.scriptAddr)
	if err != nil {
		return NewErrRefund(err)
	}
	output := outs.Outputs[0]

	if spent, err := atom.ScriptSpent(atom.scriptAddr); spent || err != nil {
		if err != nil {
			return err
		}
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}

	// create bitcoin script to pay to the user's personal address
	payToAddrScript, err := txscript.PayToAddrScript(atom.key.Address)
	if err != nil {
		return NewErrRefund(err)
	}

	// build transaction
	hashBytes, err := hex.DecodeString(output.TxHash)
	if err != nil {
		return NewErrRedeem(err)
	}
	txHash, err := chainhash.NewHash(hashBytes)
	if err != nil {
		return NewErrRedeem(err)
	}
	refundTx := wire.NewMsgTx(atom.txVersion)
	refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txHash, output.Vout), nil, nil))
	refundTx.AddTxOut(wire.NewTxOut(sendValue-atom.fee, payToAddrScript))

	// sign transaction
	refundSig, refundPubKey, err := sign(refundTx, int(output.Vout), atom.script, atom.key)
	if err != nil {
		return NewErrRefund(err)
	}

	// build signature script
	refundSigScript, err := newRefundScript(atom.script, refundSig, refundPubKey)
	if err != nil {
		return NewErrRefund(err)
	}
	refundTx.TxIn[int(output.Vout)].SignatureScript = refundSigScript

	scriptPubKey, err := hex.DecodeString(output.ScriptPubKey)
	if err != nil {
		return NewErrRefund(err)
	}

	// TODO: Remove this verify and force to verify before publishing Tx.
	if atom.verify {
		// verifying the refund script
		e, err := txscript.NewEngine(scriptPubKey, refundTx, int(output.Vout),
			txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), sendValue)
		if err != nil {
			return NewErrRefund(err)
		}
		err = e.Execute()
		if err != nil {
			return NewErrRefund(NewErrScriptExec(err))
		}
	}

	// marshal signed transaction information
	var refundTxBuffer bytes.Buffer
	refundTxBuffer.Grow(refundTx.SerializeSize())
	if err := refundTx.Serialize(&refundTxBuffer); err != nil {
		return NewErrRefund(err)
	}

	if err := atom.PublishTransaction(refundTxBuffer.Bytes(),
		func() (bool, error) {
			return atom.ScriptSpent(atom.scriptAddr)
		},
	); err != nil {
		return err
	}

	atom.LogInfo(atom.req.UID, "Refunded the transaction on Bitcoin blockchain")
	return nil
}
