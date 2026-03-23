APP_NAME := quic-raft
BIN_DIR := bin
GOCACHE := $(CURDIR)/.cache/go-build
GOMODCACHE := $(CURDIR)/.cache/gomod
GOENV := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

.PHONY: doctor fmt test build run-node run-client run-bench clean

doctor:
	./scripts/check_env.sh

fmt:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go fmt ./...

test:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go test ./...

build:
	mkdir -p $(BIN_DIR)
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go build -o $(BIN_DIR)/node ./cmd/node
	$(GOENV) go build -o $(BIN_DIR)/client ./cmd/client
	$(GOENV) go build -o $(BIN_DIR)/bench ./cmd/bench

run-node:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go run ./cmd/node --id node-1 --listen 127.0.0.1:9001 --transport quic

run-client:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go run ./cmd/client --addr 127.0.0.1:9001 --command put --key demo --value hello

run-bench:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go run ./cmd/bench --scenario baseline --nodes 3 --transport tcp

clean:
	rm -rf $(BIN_DIR)

