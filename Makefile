.PHONY: all build build-win windows linux darwin

LINUX_TARGET = swapper_linux_amd64.zip
DARWIN_TARGET = swapper_darwin_amd64.zip
WIN_TARGET = swapper_windows_amd64.zip

all: $(DARWIN_TARGET) $(WIN_TARGET) $(LINUX_TARGET)

$(DARWIN_TARGET): darwin

$(LINUX_TARGET): linux

$(WIN_TARGET): windows

darwin: build
	mkdir -p bin
	mv installer-darwin-10.6-amd64 bin/installer
	mv updater-darwin-10.6-amd64 bin/updater
	mv swapperd-unix-darwin-10.6-amd64 bin/swapperd
	mv updater-unix-darwin-10.6-amd64 bin/swapperd-updater
	mv uninstaller-darwin-10.6-amd64 bin/uninstaller
	zip -r $(DARWIN_TARGET) bin
	rm -rf bin

linux: build
	mkdir -p bin
	mv installer-linux-amd64 bin/installer
	mv updater-linux-amd64 bin/updater
	mv swapperd-unix-linux-amd64 bin/swapperd
	mv updater-unix-linux-amd64 bin/swapperd-updater
	mv uninstaller-linux-amd64 bin/uninstaller
	zip -r $(LINUX_TARGET) bin
	rm -rf bin

windows: build-win
	mkdir -p bin
	mv installer-windows-4.0-amd64.exe bin/installer.exe
	mv updater-windows-4.0-amd64.exe bin/updater.exe
	mv swapperd-win-windows-4.0-amd64.exe bin/swapperd.exe
	mv updater-win-windows-4.0-amd64.exe bin/swapperd-updater.exe
	mv uninstaller-windows-4.0-amd64.exe bin/uninstaller.exe
	zip -r $(WIN_TARGET) bin
	rm -rf bin

build:
	xgo --targets=darwin/amd64,linux/amd64 ./cmd/installer
	xgo --targets=darwin/amd64,linux/amd64 ./cmd/updater
	xgo --targets=darwin/amd64,linux/amd64 ./cmd/swapperd-unix
	xgo --targets=darwin/amd64,linux/amd64 ./cmd/updater-unix
	xgo --targets=darwin/amd64,linux/amd64 ./cmd/uninstaller

build-win:
	xgo --targets=windows/amd64 -ldflags "-H windowsgui" ./cmd/installer
	xgo --targets=windows/amd64 -ldflags "-H windowsgui" ./cmd/updater
	xgo --targets=windows/amd64 -ldflags "-H windowsgui" ./cmd/swapperd-win
	xgo --targets=windows/amd64 -ldflags "-H windowsgui" ./cmd/updater-win
	xgo --targets=windows/amd64 -ldflags "-H windowsgui" ./cmd/uninstaller

clean:
	rm -rf $(DARWIN_TARGET) $(WIN_TARGET) $(LINUX_TARGET)
