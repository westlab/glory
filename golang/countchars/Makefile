.PHONY: build
build:
	go build -o bin/countchars ./cmd

.PHONY: deps
deps:
	go mod tidy

.PHONY: test
test:
	go test -v ./
