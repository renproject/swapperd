#!/usr/bin/env bash

mkdir -p bin
cd installer
go build .
xgo --targets=linux/amd64 .
cd ../swapper/
go build .
xgo --targets=linux/amd64 .

cd ..
mv installer/installer bin/installer
mv swapper/swapper bin/swapper
zip -r swapper_darwin_amd64.zip bin

mv installer/installer-linux-amd64 bin/installer
mv swapper/swapper-linux-amd64 bin/swapper
zip -r swapper_linux_amd64.zip bin

rm -rf bin