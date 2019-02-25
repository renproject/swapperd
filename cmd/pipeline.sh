#!/usr/bin/env bash

mkdir -p bin
cd installer
go build .
xgo --targets=linux/amd64,windows/amd64 .

cd ../updater/
go build .
xgo --targets=linux/amd64,windows/amd64 .

cd ../swapperd-unix/
go build .
xgo --targets=linux/amd64 .
cd ../swapperd-win/
xgo --targets=windows/amd64 .

cd ../updater-unix/
go build .
xgo --targets=linux/amd64 .
cd ../updater-win/
xgo --targets=windows/amd64 .

cd ../uninstaller/
go build .
xgo --targets=linux/amd64,windows/amd64 .

cd ../
mv installer/installer bin/installer
mv updater/updater bin/updater
mv swapperd-unix/swapperd-unix bin/swapperd
mv updater-unix/updater-unix bin/swapperd-updater
mv uninstaller/uninstaller bin/uninstaller
zip -r swapper_darwin_amd64.zip bin

mv installer/installer-linux-amd64 bin/installer
mv updater/updater-linux-amd64 bin/updater
mv swapperd-unix/swapperd-unix-linux-amd64 bin/swapperd
mv updater-unix/updater-unix-linux-amd64 bin/swapperd-updater
mv uninstaller/uninstaller-linux-amd64 bin/uninstaller
zip -r swapper_linux_amd64.zip bin
rm -rf bin

mkdir -p bin
mv installer/installer-windows-4.0-amd64.exe bin/installer.exe
mv updater/updater-windows-4.0-amd64.exe bin/updater.exe
mv swapperd-win/swapperd-win-windows-4.0-amd64.exe bin/swapperd.exe 
mv updater-win/updater-win-windows-4.0-amd64.exe bin/swapperd-updater.exe 
mv uninstaller/uninstaller-windows-4.0-amd64.exe bin/uninstaller.exe 
zip -r swapper_windows_amd64.zip bin
rm -rf bin