CURRENT_MILLIS=`date +%s`

release-clean:
	rm -rf dist/
	rm -rf moo4plex.app/Contents

release-deps:
	go install fyne.io/fyne/v2/cmd/fyne@v2.2.4

release-build:
	fyne package --os darwin --release --appBuild $(CURRENT_MILLIS) --icon icon.png
	mkdir -p dist/darwin/amd64
	tar -czvf dist/darwin/amd64/osx-x64.tar.gz moo4plex.app

release: release-clean release-build

run:
	go run main.go