release-clean:
	rm -rf dist

release-build:
	env GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/moo4plex cmd/main.go 
	tar -C dist/darwin/amd64/ -czvf dist/darwin/amd64/osx-x64.tar.gz moo4plex

release: release-clean release-build

run:
	go run cmd/main.go