package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/swapperd/foundation/swap"
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
//				OP_SIZE
// 				32
//				OP_EQUALVERIFY
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
		b.AddOp(txscript.OP_SIZE)
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
//			1 (True)
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
//			0 (False)
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

func addressToPubKeyHash(addrString string, chainParams *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
	btcAddr, err := btcutil.DecodeAddress(addrString, chainParams)
	if err != nil {
		return nil, fmt.Errorf("address %s is not "+
			"intended for use on %v", addrString, chainParams.Name)
	}
	addr, ok := btcAddr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return nil, errors.New("%s is not p2pkh address")
	}
	return addr, nil
}

func buildInitiateScript(swap swap.Swap, Net *chaincfg.Params) ([]byte, string, error) {
	// decoding bitcoin addresses
	FundingAddr, err := addressToPubKeyHash(swap.FundingAddress, Net)
	if err != nil {
		return nil, "", NewErrDecodeAddress(swap.FundingAddress, err)
	}

	SpendingAddr, err := addressToPubKeyHash(swap.SpendingAddress, Net)
	if err != nil {
		return nil, "", NewErrDecodeAddress(swap.SpendingAddress, err)
	}

	// creating atomic swap initiate script, addressScriptHash and script to
	// deposit bitcoin tokens.
	initiateScript, err := newInitiateScript(
		FundingAddr.Hash160(),
		SpendingAddr.Hash160(),
		swap.TimeLock,
		swap.SecretHash[:],
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
