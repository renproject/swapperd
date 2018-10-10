package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/swapperd/adapter/keystore"
	"github.com/republicprotocol/swapperd/foundation"
	"golang.org/x/crypto/ripemd160"
)

// AtomicSwapRefundScriptSize is the size of the Bitcoin Atomic Swap's
// RefundScript
const AtomicSwapRefundScriptSize = 1 + 73 + 1 + 33 + 1

// AtomicSwapRedeemScriptSize is the size of the Bitcoin Atomic Swap's
// RedeemScript
const AtomicSwapRedeemScriptSize = 1 + 73 + 1 + 33 + 1 + 32 + 1

// newInitiateScript creates a Bitcoin Atomic Swap initiate script.
//
//			OP_IF
//				OP_SHA256
//				<secret_hash>
//				OP_EQUALVERIFY
//				OP_DUP
//				OP_HASH160
//				<foreign_address>
//			OP_ELSE
//				<lock_time>
//				OP_CHECKLOCKTIMEVERIFY
//				OP_DROP
//				OP_DUP
//				OP_HASH160
//				<personal_address>
//			OP_ENDIF
//			OP_EQUALVERIFY
//			OP_CHECKSIG
//
func newInitiateScript(pkhMe, pkhThem *[ripemd160.Size]byte, locktime int64, secretHash []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()

	b.AddOp(txscript.OP_IF)
	{
		b.AddOp(txscript.OP_SIZE) // TODO: Update comments
		b.AddData([]byte{32})
		b.AddOp(txscript.OP_EQUALVERIFY)
		b.AddOp(txscript.OP_SHA256)
		b.AddData(secretHash)
		b.AddOp(txscript.OP_EQUALVERIFY)
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhThem[:])
	}
	b.AddOp(txscript.OP_ELSE)
	{
		b.AddInt64(locktime)
		b.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
		b.AddOp(txscript.OP_DROP)
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhMe[:])
	}
	b.AddOp(txscript.OP_ENDIF)
	b.AddOp(txscript.OP_EQUALVERIFY)
	b.AddOp(txscript.OP_CHECKSIG)

	return b.Script()
}

// newRedeemScript creates a Redeem Script for the Bitcoin Atomic Swap.
//
//			<Signature>
//			<PublicKey>
//			<Secret>
//			<True>(Int 1)
//			<InitiateScript>
//
func newRedeemScript(initiateScript, sig, pubkey []byte, secret [32]byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddData(secret[:])
	b.AddInt64(1)
	b.AddData(initiateScript)
	return b.Script()
}

// newRefundScript creates a Bitcoin Refund Atomic Swap.
//
//			<Signature>
//			<PublicKey>
//			<False>(Int 0)
//			<InitiateScript>
//
func newRefundScript(initiateScript, sig, pubkey []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddInt64(0)
	b.AddData(initiateScript)
	return b.Script()
}

// helper functions
func sign(tx *wire.MsgTx, idx int, pkScript []byte, key keystore.BitcoinKey) (sig, pubkey []byte, err error) {
	sig, err = txscript.RawTxInSignature(tx, idx, pkScript, txscript.SigHashAll, key.PrivateKey)
	if err != nil {
		return nil, nil, err
	}
	return sig, key.PublicKey, nil
}

func addressToPubKeyHash(addr string, chainParams *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
	btcAddr, err := btcutil.DecodeAddress(addr, chainParams)
	if err != nil {
		return nil, fmt.Errorf("address %s is not "+
			"intended for use on %v", addr, chainParams.Name)
	}
	Addr, ok := btcAddr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return nil, fmt.Errorf("address %s is not Pay to Public Key Hash")
	}
	return Addr, nil
}

func buildInitiateScript(personalAddress string, req foundation.Swap, Net *chaincfg.Params) ([]byte, string, error) {
	var PayerAddress, SpenderAddress string
	var locktime int64

	if (req.IsFirst && req.SendToken == foundation.TokenBTC) || (!req.IsFirst && req.ReceiveToken == foundation.TokenBTC) {
		locktime = req.TimeLock
	} else {
		locktime = req.TimeLock - 24*60*60
	}

	if req.SendToken == foundation.TokenBTC {
		PayerAddress = personalAddress
		SpenderAddress = req.SendToAddress
	} else {
		PayerAddress = req.ReceiveFromAddress
		SpenderAddress = personalAddress
	}

	// decoding bitcoin addresses
	PayerAddr, err := addressToPubKeyHash(PayerAddress, Net)
	if err != nil {
		return nil, "", NewErrDecodeAddress(PayerAddress, err)
	}

	SpenderAddr, err := addressToPubKeyHash(SpenderAddress, Net)
	if err != nil {
		return nil, "", NewErrDecodeAddress(SpenderAddress, err)
	}

	// creating atomic swap initiate script, addressScriptHash and script to
	// deposit bitcoin tokens.
	initiateScript, err := newInitiateScript(
		PayerAddr.Hash160(),
		SpenderAddr.Hash160(),
		locktime,
		req.SecretHash[:],
	)
	if err != nil {
		return nil, "", NewErrBuildScript(err)
	}
	initiateScriptP2SH, err := btcutil.NewAddressScriptHash(initiateScript, Net)
	if err != nil {
		return nil, "", NewErrBuildScript(err)
	}

	return initiateScript, initiateScriptP2SH.EncodeAddress(), nil
}

func verifyTransaction(scriptPubKey []byte, tx *wire.MsgTx, idx int, receiveValue int64) error {
	e, err := txscript.NewEngine(scriptPubKey, tx, idx,
		txscript.StandardVerifyFlags, txscript.NewSigCache(10),
		txscript.NewTxSigHashes(tx), receiveValue)
	if err != nil {
		return err
	}
	if err := e.Execute(); err != nil {
		return NewErrScriptExec(err)
	}
	return nil
}
