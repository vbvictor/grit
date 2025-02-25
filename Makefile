.PHONY: all build test clean lint vet fmt run

all: build lint test

grit:
	go build -o grit ./grit

build:
	go build ./...

install:
	go install ./grit

test:
	go test -v -cover ./...

clean:
	go clean
	rm -f coverage.out

lint:
	golangci-lint run

format:
	gofumpt -w .

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

deps:
	go mod download
	go mod tidy
