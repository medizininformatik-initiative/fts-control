.PHONY: build test coverage vuln

build:
	go build -o ftsctl

test:
	go test -v ./...

coverage:
	go test -v -coverprofile=coverage.txt ./...

vuln:
	go tool govulncheck ./...
