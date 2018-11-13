#!/usr/bin/env bash

mkdir -p bin
cd installer
go build .
xgo --targets=linux/amd64,windows/amd64 .
cd ../swapperd-unix/
go build .
xgo --targets=linux/amd64 .
cd ../swapperd-win/
xgo --targets=windows/amd64 .

cd ..
mv installer/installer bin/installer
mv swapperd-unix/swapperd-unix bin/swapperd
zip -r swapper_darwin_amd64.zip bin

mv installer/installer-linux-amd64 bin/installer
mv swapperd-unix/swapperd-unix-linux-amd64 bin/swapperd
zip -r swapper_linux_amd64.zip bin
rm -rf bin

mkdir -p bin
mv installer/installer-windows-4.0-amd64.exe bin/installer.exe
mv swapperd-win/swapperd-win-windows-4.0-amd64.exe bin/swapper.exe 
zip -r swapper_windows_amd64.zip bin
rm -rf bin