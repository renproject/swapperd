#!/usr/bin/env bash

mkdir -p bin
cd installer-unix
go build .
xgo --targets=linux/amd64 .
cd ../installer-win
xgo --targets=windows/amd64 -ldflags "-H windowsgui" .

cd ../updater-win/
xgo --targets=windows/amd64 -ldflags "-H windowsgui" .

cd ../swapperd-unix/
go build .
xgo --targets=linux/amd64 .
cd ../swapperd-win/
xgo --targets=windows/amd64 -ldflags "-H windowsgui" .

cd ../uninstaller-unix/
go build .
xgo --targets=linux/amd64 .
cd ../uninstaller-win/
xgo --targets=windows/amd64 -ldflags "-H windowsgui" .

cd ../
mv installer-unix/installer-unix bin/installer
mv swapperd-unix/swapperd-unix bin/swapperd
mv uninstaller-unix/uninstaller-unix bin/uninstaller
zip -r swapper_darwin_amd64.zip bin

mv installer-unix/installer-unix-linux-amd64 bin/installer
mv swapperd-unix/swapperd-unix-linux-amd64 bin/swapperd
mv uninstaller-unix/uninstaller-unix-linux-amd64 bin/uninstaller
zip -r swapper_linux_amd64.zip bin
rm -rf bin

mkdir -p bin
mv installer-win/installer-win-windows-4.0-amd64.exe bin/installer.exe
mv updater-win/updater-win-windows-4.0-amd64.exe bin/updater.exe
mv swapperd-win/swapperd-win-windows-4.0-amd64.exe bin/swapperd.exe 
mv uninstaller-win/uninstaller-win-windows-4.0-amd64.exe bin/uninstaller.exe 
zip -r swapper_windows_amd64.zip bin
rm -rf bin