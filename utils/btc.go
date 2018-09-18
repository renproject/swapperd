package utils

import (
	"os"

	"github.com/btcsuite/btcd/chaincfg"
)

func GetHome() string {
	winHome := os.Getenv("userprofile")
	unixHome := os.Getenv("HOME")
	if winHome != "" {
		return winHome
	}

	if unixHome != "" {
		return unixHome
	}
	panic("unknown Operating System")
}

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
