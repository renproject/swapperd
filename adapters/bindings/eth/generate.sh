#!/bin/sh

cd ../../../drivers/atom-sol

if [ $1 = "--branch" ]
then
    git checkout $2
fi

### GENERATE BINDINGS HERE
abigen --sol ./contracts/AtomInfo.sol -pkg eth --out ../../adapters/bindings/eth/info.go
abigen --sol ./contracts/AtomNetwork.sol -pkg eth --out ../../adapters/bindings/eth/network.go
abigen --sol ./contracts/AtomSwap.sol -pkg eth --out ../../adapters/bindings/eth/atom.go
abigen --sol ./contracts/AtomWallet.sol -pkg eth --out ../../adapters/bindings/eth/wallet.go