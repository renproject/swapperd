#!/bin/sh

cd ../../../drivers/atom-sol

if [ $1 -e "--branch" ]
then
    git checkout $2
fi

npm install

# Setup
sed -i -e 's/"openzeppelin-solidity\/contracts\//".\/openzeppelin-solidity\/contracts\//' contracts/*.sol
sed -i -e 's/"republic-sol\/contracts\//".\/republic-sol\/contracts\//' contracts/*.sol
mkdir ./contracts/openzeppelin-solidity
mkdir ./contracts/republic-sol
cp -r ./node_modules/openzeppelin-solidity/contracts ./contracts/openzeppelin-solidity/contracts
cp -r ./node_modules/republic-sol/contracts ./contracts/republic-sol/contracts

cd contracts/republic-sol
sed -i -e 's/"openzeppelin-solidity\/contracts\//"..\/..\/openzeppelin-solidity\/contracts\//' contracts/*.sol
sed -i -e 's/"openzeppelin-solidity\/contracts\//"..\/..\/..\/openzeppelin-solidity\/contracts\//' contracts/*/*.sol
mkdir ./contracts/openzeppelin-solidity
cp -r ../../node_modules/openzeppelin-solidity/contracts ./contracts/openzeppelin-solidity/contracts
cd ../..

### GENERATE BINDINGS HERE ###
abigen --sol ./contracts/AtomicInfo.sol -pkg eth --out ../../adapters/bindings/eth/info.go
abigen --sol ./contracts/AtomicSwap.sol -pkg eth --out ../../adapters/bindings/eth/atom.go
abigen --sol ./contracts/RenExAtomicSettlement.sol -pkg eth --out ../../adapters/bindings/eth/settlement.go

# Revert setup
sed -i -e 's/".\/openzeppelin-solidity\/contracts\//"openzeppelin-solidity\/contracts\//' contracts/*.sol
sed -i -e 's/"..\/openzeppelin-solidity\/contracts\//"openzeppelin-solidity\/contracts\//' contracts/*/*.sol
sed -i -e 's/".\/republic-sol\/contracts\//"republic-sol\/contracts\//' contracts/*.sol
rm -r ./contracts/openzeppelin-solidity
rm -r ./contracts/republic-sol
