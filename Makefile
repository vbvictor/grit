.PHONY: all build test clean lint vet fmt run

all: build lint test install

grit:
	go build -o grit ./grit

build:
	go build ./...
	go build -o grit ./grit

install:
	go install ./grit

test:
	go test -v ./...

test-with-coverage:
	go test -v -cover -coverprofile=coverage.out ./...

clean:
	go clean
	rm -f coverage.out

lint:
	golangci-lint run

format:
	gofumpt -w .

coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

deps:
	go mod tidy
	go mod download
