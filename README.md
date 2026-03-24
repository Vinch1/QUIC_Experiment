# quic-raft

Proof-of-concept project for comparing `TCP` vs `QUIC` as the inter-node transport layer beneath a Raft-style replication flow.

## Current Goals

- Build a runnable multi-node experiment skeleton
- Keep `client -> node` and `node -> node` concerns clearly separated
- Make `node -> node` transport switchable between `TCP` and `QUIC`
- Validate majority replication behavior locally before implementing full Raft

## Current Status

This repository now provides:

- Go project scaffolding
- `node` / `client` / `bench` entry points
- Real `node -> node` TCP transport
- Real `node -> node` QUIC transport
- HTTP control plane with `put` / `get` / `status`
- Minimal majority replication path with a static leader

Not yet implemented:

- Dynamic Raft leader election
- Raft term / log index / commit index semantics
- Protobuf-generated message layer
- Snapshot / WAL / membership change
- Benchmark automation and metrics dashboard

## Quick Start

1. Start a 3-node TCP cluster

```bash
go run ./cmd/node --id node-1 --control-addr 127.0.0.1:9001 --raft-addr 127.0.0.1:7001 --transport quic --leader --peers node-2=127.0.0.1:7002,node-3=127.0.0.1:7003
go run ./cmd/node --id node-2 --control-addr 127.0.0.1:9002 --raft-addr 127.0.0.1:7002 --transport quic --leader-id node-1 --peers node-1=127.0.0.1:7001,node-3=127.0.0.1:7003
go run ./cmd/node --id node-3 --control-addr 127.0.0.1:9003 --raft-addr 127.0.0.1:7003 --transport quic --leader-id node-1 --peers node-1=127.0.0.1:7001,node-2=127.0.0.1:7002
```

2. Switch to QUIC by changing only the transport flag

```bash
go run ./cmd/node --id node-1 --control-addr 127.0.0.1:9001 --raft-addr 127.0.0.1:7001 --transport quic --leader --peers node-2=127.0.0.1:7002,node-3=127.0.0.1:7003
```

3. Write through the control plane

```bash
make run-client
```

4. Read from a node

```bash
go run ./cmd/client --addr 127.0.0.1:9002 --command get --key demo
```

5. Inspect status

```bash
go run ./cmd/client --addr 127.0.0.1:9001 --command status
```

## Protocol Split

- `client -> node`: HTTP over TCP
- `node -> node`: configurable `TCP` or `QUIC`

That split keeps the experiment focused on the transport used by replication traffic, instead of mixing client protocol changes into the result.

## Documentation

- Setup: `docs/setup.md`
- Architecture: `docs/architecture.md`
- Experiment Plan: `docs/experiment-plan.md`
