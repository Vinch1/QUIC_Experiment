# QUIC_Experiment

Class project - A proof-of-concept experiment comparing TCP vs QUIC as the transport layer for Raft consensus.

## Current Goals

- Build a minimal runnable Raft experiment skeleton
- Keep the transport layer pluggable for fair TCP/QUIC comparison
- Support local 3-node experiments first, then extend to weak network and fault injection
- Complete environment setup, scaffolding, and experiment plan before implementing protocols and algorithms

## Current Status

This repository provides:

- Go project scaffolding
- `node` / `client` / `bench` entry points
- Minimal interface skeletons for Raft, transport layer, storage layer, and state machine
- HTTP-based minimal control plane with `put` / `get` / `status`
- Local environment check script
- Project plan and setup documentation
- Reserved Docker Compose / Proto / netem directories

Not yet implemented:

- Complete Raft leader election and log replication
- Actual TCP / QUIC network implementation
- Protobuf code generation
- Prometheus metrics collection
- Automated benchmark executor

## Directory Structure

```text
.
├── api/proto
├── cmd
│   ├── bench
│   ├── client
│   └── node
├── deploy
│   ├── docker-compose.yml
│   └── netem
├── docs
├── internal
│   ├── cluster
│   ├── metrics
│   ├── raft
│   ├── statemachine
│   ├── storage
│   └── transport
├── scripts
├── Dockerfile
├── Makefile
└── go.mod
```

## Quick Start

1. Check local dependencies

```bash
make doctor
```

2. Run node scaffold

```bash
make run-node
```

3. Write a key-value pair

```bash
make run-client
```

4. Read a key-value pair

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command get --key demo
```

5. Check node status

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command status
```

## Recommended Development Order

1. Complete `api/proto/raft.proto`
2. Decouple control plane from inter-node transport
3. Implement `internal/transport/tcp` and `internal/transport/quic`
4. Implement minimal leader election and log replication in `internal/raft`
5. Add baseline test scenarios in `cmd/bench`

## Documentation

- Setup: `docs/setup.md`
- Architecture: `docs/architecture.md`
- Experiment Plan: `docs/experiment-plan.md`
