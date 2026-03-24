APP_NAME := quic-raft
BIN_DIR := bin
GOCACHE := $(CURDIR)/.cache/go-build
GOMODCACHE := $(CURDIR)/.cache/gomod
GOENV := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

.PHONY: doctor fmt test build run-node run-client run-bench stop-cluster clean

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
	$(GOENV) go run ./cmd/node --id node-1 --control-addr 127.0.0.1:9001 --raft-addr 127.0.0.1:7001 --transport tcp --leader

run-client:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go run ./cmd/client --addr 127.0.0.1:9001 --command put --key demo --value hello

run-bench:
	mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GOENV) go run ./cmd/bench --scenario baseline --nodes 3 --transport tcp

stop-cluster:
	pkill -f '/Users/leo/Code/Local/quic_experiment/bin/node' || true
	pkill -f 'go run ./cmd/node' || true
	for port in 7001 7002 7003 9001 9002 9003; do \
		pids=$$(lsof -ti tcp:$$port -sTCP:LISTEN 2>/dev/null || true); \
		if [ -n "$$pids" ]; then kill $$pids || true; fi; \
		pids=$$(lsof -ti udp:$$port 2>/dev/null || true); \
		if [ -n "$$pids" ]; then kill $$pids || true; fi; \
	done

clean:
	rm -rf $(BIN_DIR)
