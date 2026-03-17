.PHONY: build test lint clean install docker-build

BINARY=muah-runner
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/muah-runner/

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

install:
	go install $(LDFLAGS) ./cmd/muah-runner/

docker-build:
	docker build -t muah-runner:$(VERSION) .
