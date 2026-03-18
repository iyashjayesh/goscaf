build:
	go build -o bin/goscaf ./main.go

install:
	go install ./...

run:
	go run ./main.go

test:
	GOTOOLCHAIN=auto go test -race -cover ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

smoke-test: build
	./bin/goscaf init smoke-test-project --defaults
	ls smoke-test-project/
	rm -rf smoke-test-project/
	@echo "Smoke test passed!"

fmt:
	go fmt ./...

clean:
	rm -rf bin/
	rm -rf coverage.out
	rm -rf smoke-test-project/

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: build install run test lint tidy smoke-test fmt clean install-tools
