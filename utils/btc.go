package utils

import (
	"os"
	"runtime"

	"github.com/btcsuite/btcd/chaincfg"
)

func GetHome() string {
	system := runtime.GOOS
	switch system{
	case "window":
		return os.Getenv("userprofile")
	case "linux", "darwin":
		return os.Getenv("HOME")
	default:
		panic("unknown Operating System")
	}
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
