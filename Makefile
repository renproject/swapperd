MAIN_VERSION = $(shell cat ./VERSION)
BRANCH = $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT_HASH = $(shell git describe --always --long)
FULL_VERSION = ${MAIN_VERSION}-${BRANCH}-${COMMIT_HASH}

LDFLAGS = -X main.version=${FULL_VERSION}
WIN_LDFLAGS = ${LDFLAGS} -H windowsgui

LOCAL_TARGET = swapper_local.zip

LINUX_TARGET = swapper_linux_amd64.zip
DARWIN_TARGET = swapper_darwin_amd64.zip
WIN_TARGET = swapper_windows_amd64.zip

all: $(DARWIN_TARGET) $(WIN_TARGET) $(LINUX_TARGET)

$(DARWIN_TARGET): darwin

$(LINUX_TARGET): linux

$(WIN_TARGET): windows

darwin: build-unix
	@mkdir -p bin
	@mv installer-darwin-10.6-amd64 bin/installer
	@mv updater-darwin-10.6-amd64 bin/updater
	@mv swapperd-unix-darwin-10.6-amd64 bin/swapperd
	@mv updater-unix-darwin-10.6-amd64 bin/swapperd-updater
	@mv uninstaller-darwin-10.6-amd64 bin/uninstaller
	@zip -r ${DARWIN_TARGET} bin
	@rm -rf bin
	@echo
	@echo "Compiled ${DARWIN_TARGET} (${FULL_VERSION})"

linux: build-unix
	@mkdir -p bin
	@mv installer-linux-amd64 bin/installer
	@mv updater-linux-amd64 bin/updater
	@mv swapperd-unix-linux-amd64 bin/swapperd
	@mv updater-unix-linux-amd64 bin/swapperd-updater
	@mv uninstaller-linux-amd64 bin/uninstaller
	@zip -r ${LINUX_TARGET} bin
	@rm -rf bin
	@echo
	@echo "Compiled ${LINUX_TARGET} (${FULL_VERSION})"

windows: build-win
	@mkdir -p bin
	@mv installer-windows-4.0-amd64.exe bin/installer.exe
	@mv updater-windows-4.0-amd64.exe bin/updater.exe
	@mv swapperd-win-windows-4.0-amd64.exe bin/swapperd.exe
	@mv updater-win-windows-4.0-amd64.exe bin/swapperd-updater.exe
	@mv uninstaller-windows-4.0-amd64.exe bin/uninstaller.exe
	@zip -r ${WIN_TARGET} bin
	@rm -rf bin
	@echo
	@echo "Compiled ${WIN_TARGET} (${FULL_VERSION})"

build:
	@mkdir -p bin
	go build -ldflags="${LDFLAGS}" -o bin/installer ./cmd/installer
	go build -ldflags="${LDFLAGS}" -o bin/updater ./cmd/updater
	go build -ldflags="${LDFLAGS}" -o bin/swapperd ./cmd/swapperd-unix
	go build -ldflags="${LDFLAGS}" -o bin/swapperd-updater ./cmd/updater-unix
	go build -ldflags="${LDFLAGS}" -o bin/uninstaller ./cmd/uninstaller
	@echo
	@zip -r ${LOCAL_TARGET} bin
	@rm -rf bin
	@echo
	@echo "Compiled ${LOCAL_TARGET} (${FULL_VERSION})"

clean:
	rm -rf ${DARWIN_TARGET} ${WIN_TARGET} ${LINUX_TARGET}

version:
	@echo ${FULL_VERSION}

define build_unix
	xgo --targets=darwin/amd64,linux/amd64 -ldflags "${LDFLAGS}" $(1)
endef

define build_win
	xgo --targets=windows/amd64 -ldflags "${WIN_LDFLAGS}" $(1)
endef

build-unix:
	$(call build_unix,./cmd/installer)
	$(call build_unix,./cmd/updater)
	$(call build_unix,./cmd/swapperd-unix)
	$(call build_unix,./cmd/updater-unix)
	$(call build_unix,./cmd/uninstaller)

build-win:
	$(call build_win,./cmd/installer)
	$(call build_win,./cmd/updater)
	$(call build_win,./cmd/swapperd-win)
	$(call build_win,./cmd/updater-win)
	$(call build_win,./cmd/uninstaller)

.PHONY: all build build-unix build-win windows linux darwin version clean

