#!/usr/bin/env bash

mkdir -p bin
cd installer
go build .
xgo --targets=linux/amd64 .
cd ../swapperd/
go build .
xgo --targets=linux/amd64 .

cd ..
mv installer/installer bin/installer
mv swapperd/swapperd bin/swapperd
zip -r swapper_darwin_amd64.zip bin

mv installer/installer-linux-amd64 bin/installer
mv swapperd/swapperd-linux-amd64 bin/swapperd
zip -r swapper_linux_amd64.zip bin
rm -rf bin