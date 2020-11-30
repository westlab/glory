
.PHONY: build
build:
	go build -o bin/batch ./cmd/batch
	go build -o bin/fetchID ./cmd/fetchID

.PHONY: deps
deps:
	go mod init
	go mod tidy

.PHONY: test
test:
	go test -v ./
