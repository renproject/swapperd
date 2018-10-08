#!/usr/bin/env bash

mkdir -p bin
cd installer
go build .
xgo --targets=linux/amd64,windows/amd64 .
cd ../swapper/
go build .
xgo --targets=linux/amd64,windows/amd64 .
cd ../installer-win/
xgo --targets=windows/amd64 .

cd ..
mv installer/installer bin/installer
mv swapper/swapper bin/swapper
zip -r swapper_darwin_amd64.zip bin

mv installer/installer-linux-amd64 bin/installer
mv swapper/swapper-linux-amd64 bin/swapper
zip -r swapper_linux_amd64.zip bin

rm -rf bin
mkdir -p bin

mv installer/installer-windows-4.0-amd64.exe bin/installer.exe
mv swapper/swapper-windows-4.0-amd64.exe bin/swapper.exe
mv installer-win/installer-win-windows-4.0-amd64.exe bin/service.exe
zip -r swapper_windows_amd64.zip bin

rm -rf bin