GOOS ?= linux
GOARCH ?= arm64 # amd64

ROOT_PACKAGE := github.com/westlab/glory
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/')

GLORY_BIN := ./app/glory
COUNT_CHAR_BIN := ./app/cmd/count_char

.PHONY: all
all: $(GLORY_BIN) $(COUNT_CHAR_BIN)

$(GLORY_BIN): $(SRC)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -a -tags netgo -installsuffix netgo -ldflags \
	'-s -w -extldflags "-static"' -o $@

$(COUNT_CHAR_BIN): $(SRC)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -a -tags netgo -installsuffix netgo -ldflags \
	'-s -w -extldflags "-static"' -o $@ $(ROOT_PACKAGE)/cmd/count_char

.PHONY: build
build: $(GLORY_BIN)

.PHONY: build_count_char
build_count_char: $(COUNT_CHAR_BIN)

.PHONY: clean
clean:
	@rm -f $(GLORY_BIN)
	@rm -f $(COUNT_CHAR_BIN)
