package utils

import (
	"github.com/btcsuite/btcd/chaincfg"
)

func GetChainParams(network string) *chaincfg.Params {
	switch network {
	case "regtest":
		return &chaincfg.RegressionNetParams
	case "testnet":
		return &chaincfg.TestNet3Params
	case "mainnet":
		return &chaincfg.MainNetParams
	default:
		panic("unimplemented")
	}
}
