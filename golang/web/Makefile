GOOS ?= linux
GOARCH ?= arm64

ROOT_PACKAGE := github.com/westlab/glory/web
SRC := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GLORY_BIN := ./app/glory

$(GLORY_BIN): $(SRC)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -tags netgo -installsuffix netgo -ldflags '-s -w -extldflags "-static"' -o $@

.PHONY: build
build: $(GLORY_BIN))

.PHONY: clean
clean:
	@rm -f $(GLORY_BIN)
