.PHONY: build-windows build-macos build-linux

build-windows:
	GOOS=windows GOARCH=amd64 go build

build-macos:
	GOOS=darwin GOARCH=amd64 go build

build-linux:
	GOOS=linux GOARCH=amd64 go build
