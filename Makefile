.PHONY: all build test clean lint vet fmt run

all: lint test build

build:
	go build ./grit

install:
	go install ./grit

test:
	go test -v -cover ./...

clean:
	go clean
	rm -f coverage.out

lint:
	golangci-lint run

vet:
	go vet ./...

format:
	gofumpt -w .

run:
	go run ./cmd/main.go

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

deps:
	go mod download
	go mod tidy
