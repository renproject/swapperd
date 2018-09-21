package btc

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type Conn interface {
	// SignTransaction should sign the transaction using the given bitcoin key.
	SignTransaction(tx *wire.MsgTx, key keystore.BitcoinKey, fee int64) (*wire.MsgTx, bool, error)

	// PublishTransaction should publish a signed transaction to the Bitcoin
	// blockchain.
	PublishTransaction(signedTransaction []byte) error

	// Net should return the network information of the underlying
	// Bitcoin blockchain.
	Net() *chaincfg.Params

	// GetScriptFromSpentP2SH for the given script address.
	GetScriptFromSpentP2SH(address string) ([]byte, error)

	// GetFundingTransaction for the given script address.
	GetFundingTransaction(address string) ([]byte, error)

	// Balance of the given address on Bitcoin blockchain.
	Balance(address string) (int64, error)

	// TransactionCount of the given address
	TransactionCount(address string) (int64, error)
}

type bitcoinAtom struct {
	scriptAddr string
	script     []byte
	key        keystore.BitcoinKey
	req        swap.Request
	txVersion  int32
	fee        int64
	verify     bool
	Conn
}

// NewBitcoinAtom returns a new Bitcoin Atom instance
func NewBitcoinAtom(conf config.BitcoinNetwork, key keystore.BitcoinKey, req swap.Request) (swap.Atom, error) {
	conn, err := NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}
	script, scriptAddr, err := buildInitiateScript(key.AddressString, req, conn.Net())
	if err != nil {
		return nil, err
	}
	return &bitcoinAtom{
		scriptAddr: scriptAddr,
		script:     script,
		key:        key,
		req:        req,
		txVersion:  2,
		fee:        10000,
		verify:     true,
		Conn:       conn,
	}, nil
}

// initiate the atomic swap by funding a HTLC on the Bitcoin blockchain.
func (atom *bitcoinAtom) Initiate() error {
	scriptAddr, err := btcutil.DecodeAddress(atom.scriptAddr, atom.Net())

	initiateScriptP2SHPkScript, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		return NewErrInitiate(NewErrBuildScript(err))
	}

	// creating unsigned transaction and adding transaction outputs
	unsignedTx := wire.NewMsgTx(atom.txVersion)
	unsignedTx.AddTxOut(wire.NewTxOut(atom.req.SendValue.Int64(), initiateScriptP2SHPkScript))

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

	return atom.Conn.PublishTransaction(initiateTxBuffer.Bytes())
}

func (atom *bitcoinAtom) AuditSecret() ([32]byte, error) {
	for {
		txCount, err := atom.Conn.TransactionCount(atom.scriptAddr)
		if err == nil && txCount > 1 {
			break
		}
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
			for i := 0; i < 32; i++ {
				secret[i] = push[i]
			}
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *bitcoinAtom) Refund() error {
	tx, err := atom.Conn.GetFundingTransaction(atom.scriptAddr)
	if err != nil {
		return NewErrRefund(err)
	}

	// decode initiate transaction
	var initiateTx wire.MsgTx
	if err := initiateTx.Deserialize(bytes.NewReader(tx)); err != nil {
		return NewErrRefund(NewErrDecodeTransaction(tx, err))
	}

	// create bitcoin script to pay to the user's personal address
	payToAddrScript, err := txscript.PayToAddrScript(atom.key.Address)
	if err != nil {
		return NewErrRefund(err)
	}

	// finding the relevant output index of the initiate transaction.
	outIndex := -1
	for i, out := range initiateTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, atom.Conn.Net())
		if sc == txscript.ScriptHashTy &&
			addrs[0].(*btcutil.AddressScriptHash).String() == atom.scriptAddr {
			outIndex = i
			break
		}
	}
	if outIndex == -1 {
		return NewErrRefund(ErrMalformedInitiateTx)
	}

	// build transaction
	txHash := initiateTx.TxHash()
	refundTx := wire.NewMsgTx(atom.txVersion)
	refundTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&txHash, uint32(outIndex)), nil, nil))
	refundTx.AddTxOut(wire.NewTxOut(initiateTx.TxOut[outIndex].Value-atom.fee, payToAddrScript))

	// sign transaction
	refundSig, refundPubKey, err := sign(refundTx, outIndex, atom.script, atom.key)
	if err != nil {
		return NewErrRefund(err)
	}

	// build signature script
	refundSigScript, err := newRefundScript(atom.script, refundSig, refundPubKey)
	if err != nil {
		return NewErrRefund(err)
	}
	refundTx.TxIn[outIndex].SignatureScript = refundSigScript

	if atom.verify {
		// verifying the refund script
		e, err := txscript.NewEngine(initiateTx.TxOut[outIndex].PkScript, refundTx, outIndex,
			txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), initiateTx.TxOut[outIndex].Value)
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

	return atom.PublishTransaction(refundTxBuffer.Bytes())
}

func (atom *bitcoinAtom) Audit() error {
	for {
		bal, err := atom.Conn.Balance(atom.scriptAddr)
		if err != nil {
			return NewErrAudit(err)
		}
		if bal >= atom.req.ReceiveValue.Int64() {
			return nil
		}
		if time.Now().Unix() > atom.req.TimeLock {
			break
		}
	}
	return NewErrAudit(ErrTimedOut)
}

// redeem the Atomic Swap by revealing the secret and withdrawing funds from the
// HTLC.
func (atom *bitcoinAtom) Redeem(secret [32]byte) error {
	tx, err := atom.Conn.GetFundingTransaction(atom.scriptAddr)
	if err != nil {
		return NewErrRedeem(err)
	}

	// decode initiate transaction
	var initiateTx wire.MsgTx

	if err := initiateTx.Deserialize(bytes.NewReader(tx)); err != nil {
		return NewErrRedeem(NewErrDecodeTransaction(tx, err))
	}

	// create bitcoin script to pay to the user's personal address
	payToAddrScript, err := txscript.PayToAddrScript(atom.key.Address)
	if err != nil {
		return NewErrRedeem(err)
	}

	// finding the relevant output index of the initiate transaction.
	outIndex := -1
	for i, out := range initiateTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, atom.Conn.Net())
		if sc == txscript.ScriptHashTy &&
			addrs[0].(*btcutil.AddressScriptHash).String() == atom.scriptAddr {
			outIndex = i
			break
		}
	}
	if outIndex == -1 {
		return NewErrRedeem(ErrMalformedInitiateTx)
	}

	// build transaction
	txHash := initiateTx.TxHash()
	redeemTx := wire.NewMsgTx(atom.txVersion)
	redeemTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&txHash, uint32(outIndex)), nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(initiateTx.TxOut[outIndex].Value-atom.fee, payToAddrScript))

	// sign transaction
	redeemSig, redeemPubKey, err := sign(redeemTx, outIndex, atom.script, atom.key)
	if err != nil {
		return NewErrRedeem(err)
	}

	// build signature script
	redeemSigScript, err := newRedeemScript(atom.script, redeemSig, redeemPubKey, secret)
	if err != nil {
		return NewErrRedeem(err)
	}
	redeemTx.TxIn[outIndex].SignatureScript = redeemSigScript

	if atom.verify {
		// verifying the redeem script
		e, err := txscript.NewEngine(initiateTx.TxOut[outIndex].PkScript, redeemTx, outIndex,
			txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), initiateTx.TxOut[outIndex].Value)
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

	return atom.PublishTransaction(redeemTxBuffer.Bytes())
}
