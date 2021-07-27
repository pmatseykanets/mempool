BINARY=mempool-cli

default: test

clean:
	@rm -f ./${BINARY}

test:
	go vet ./...
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

build: clean
	@CGO_ENABLED=0 go build -o ${BINARY} -ldflags "-s -w -X 'main.buildVersion=# Built $(shell date -u -R) with $(shell go version) at $(shell git rev-parse HEAD)' -X 'main.version=$(shell git describe --tags --always --dirty --match "v[0-9]*" --abbrev=4 | sed -e 's/^v//')'" ./cmd/mempool-cli

.PHONY: clean, test, build
