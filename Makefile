darwin:
	xgo --targets=darwin/amd64 ./cmd/installer
	xgo --targets=darwin/amd64 ./cmd/updater
	xgo --targets=darwin/amd64 ./cmd/swapperd-unix
	xgo --targets=darwin/amd64 ./cmd/updater-unix
	xgo --targets=darwin/amd64 ./cmd/uninstaller
	mkdir -p bin
	mv installer-darwin-10.6-amd64 bin/installer
	mv updater-darwin-10.6-amd64 bin/updater
	mv swapperd-unix-darwin-10.6-amd64 bin/swapperd
	mv updater-unix-darwin-10.6-amd64 bin/swapperd-updater
	mv uninstaller-darwin-10.6-amd64 bin/uninstaller
	zip -r swapper_darwin_amd64.zip bin
	rm -rf bin

