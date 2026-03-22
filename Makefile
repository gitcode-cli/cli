# gitcode-cli Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build clean test install

all: build

build:
	go build $(LDFLAGS) -o bin/gc ./cmd/gc

run:
	go run ./cmd/gc

test:
	go test -v -race ./...

clean:
	rm -rf bin/

install: build
	cp bin/gc /usr/local/bin/gc

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

.PHONY: help
help:
	@echo "gitcode-cli Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build     Build the binary"
	@echo "  make run       Run the application"
	@echo "  make test      Run tests"
	@echo "  make clean     Clean build artifacts"
	@echo "  make install   Install to /usr/local/bin"
	@echo "  make fmt       Format code"
	@echo "  make lint      Run linter"