package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"
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

	// Balance of the given address on Bitcoin blockchain.
	Balance(address string) (int64, error)

	// GetUnspentOutputs of the given address on Bitcoin blockchain.
	GetUnspentOutputs(address string) (UnspentOutputs, error)
}

type bitcoinAtom struct {
	scriptAddr string
	script     []byte
	key        keystore.BitcoinKey
	req        swapDomain.Request
	txVersion  int32
	fee        int64
	verify     bool
	Conn
}

// NewBitcoinAtom returns a new Bitcoin Atom instance
func NewBitcoinAtom(conf config.BitcoinNetwork, key keystore.BitcoinKey, req swapDomain.Request) (swap.Atom, error) {
	conn, err := NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}
	script, scriptAddr, err := buildInitiateScript(key.AddressString, req, conn.Net())
	if err != nil {
		return nil, err
	}
	fmt.Println("Script Address: ", scriptAddr)
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
	fmt.Println("Initiating on bitcoin blockchain")
	scriptAddr, err := btcutil.DecodeAddress(atom.scriptAddr, atom.Net())
	if err != nil {
		return err
	}

	// FIXME: Change to greater than or equal to
	if bal, err := atom.Balance(atom.scriptAddr); bal == atom.req.SendValue.Int64() || err != nil {
		if err != nil {
			return err
		}
		return swap.ErrSwapAlreadyInitiated
	}

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
	fmt.Println("Initiated on bitcoin blockchain")
	return atom.Conn.PublishTransaction(initiateTxBuffer.Bytes())
}

func (atom *bitcoinAtom) AuditSecret() ([32]byte, error) {
	for {
		fmt.Println("Auditing secret on bitcoin blockchain")
		bal, err := atom.Conn.Balance(atom.scriptAddr)
		if err == nil && bal < atom.req.SendValue.Int64() {
			break
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return [32]byte{}, NewErrAuditSecret(ErrTimedOut)
		}
		time.Sleep(1 * time.Minute)
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
			fmt.Println("... success")
			return secret, nil
		}
	}
	return [32]byte{}, NewErrAuditSecret(ErrMalformedRedeemTx)
}

// refund the Atomic Swap after expiry and withdraw funds from the HTLC.
func (atom *bitcoinAtom) Refund() error {
	// TODO: Use all the outputs
	outs, err := atom.Conn.GetUnspentOutputs(atom.scriptAddr)
	if err != nil {
		return NewErrRefund(err)
	}
	output := outs.Outputs[0]

	if bal, err := atom.Balance(atom.scriptAddr); bal == 0 || err != nil {
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
	refundTx.AddTxOut(wire.NewTxOut(atom.req.SendValue.Int64()-atom.fee, payToAddrScript))

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
			txscript.NewTxSigHashes(refundTx), atom.req.SendValue.Int64())
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
		fmt.Println("Auditing on bitcoin blockchain")
		bal, _ := atom.Conn.Balance(atom.scriptAddr)
		if bal >= atom.req.ReceiveValue.Int64() {
			fmt.Println(".... done")
			return nil
		}
		if bal != 0 {
			fmt.Printf("Expected receive value: %v actual value: %v\n", atom.req.ReceiveValue.Int64(), bal)
		}
		if time.Now().Unix() > atom.req.TimeLock {
			break
		}
		time.Sleep(1 * time.Minute)
	}
	return NewErrAudit(ErrTimedOut)
}

// redeem the Atomic Swap by revealing the secret and withdrawing funds from the
// HTLC.
func (atom *bitcoinAtom) Redeem(secret [32]byte) error {
	fmt.Println("Redeeming on bitcoin blockchain")
	outs, err := atom.Conn.GetUnspentOutputs(atom.scriptAddr)
	if err != nil {
		return NewErrRedeem(err)
	}
	output := outs.Outputs[0]

	if bal, err := atom.Balance(atom.scriptAddr); bal == 0 || err != nil {
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
	redeemTx.AddTxOut(wire.NewTxOut(atom.req.ReceiveValue.Int64()-atom.fee, payToAddrScript))

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
			txscript.NewTxSigHashes(redeemTx), atom.req.ReceiveValue.Int64())
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
	fmt.Println("Redeem Transaction: ", hex.EncodeToString(redeemTxBuffer.Bytes()))

	fmt.Println("Redeemed on bitcoin blockchain")
	return atom.PublishTransaction(redeemTxBuffer.Bytes())
}

// func (atom *bitcoinAtom) waitForInitiation(val *big.Int) error {
// 	for {
// 		fmt.Println("Auditing on bitcoin blockchain")
// 		bal, err := atom.Conn.Balance(atom.scriptAddr)
// 		if bal >= atom.req.ReceiveValue.Int64() {
// 			return nil
// 		}
// 		if time.Now().Unix() > atom.req.TimeLock {
// 			break
// 		}
// 		time.Sleep(1 * time.Minute)
// 	}
// 	return NewErrAudit(ErrTimedOut)
// }

// func (atom *bitcoinAtom) waitForRedemption() error {

// }
