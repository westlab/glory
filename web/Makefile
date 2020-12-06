
.PHONY: build
build:
	go build -o bin/glory_server ./cmd

.PHONY: deps
deps:
	go mod init
	go mod tidy

.PHONY: test
test:
	go test -v ./
