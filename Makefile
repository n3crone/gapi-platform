.PHONY: build test run

build:
	go build -v ./...
    
.PHONY: test test-coverage

test:
	go test -v ./...

test-coverage:
	go test -v -cover ./... -coverprofile=coverage.out else (go test -v -cover ./... -coverprofile=coverage.out)
	go tool cover -html=coverage.out -o coverage.html

run:
	go run cmd/server/main.go

lint:
	golangci-lint run
