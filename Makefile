build:
	go build -o bin/gostart ./main.go

install:
	go install ./...

run:
	go run ./main.go

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

smoke-test: build
	./bin/gostart init smoke-test-project --defaults
	ls smoke-test-project/
	rm -rf smoke-test-project/
	@echo "Smoke test passed!"

.PHONY: build install run test lint tidy smoke-test
