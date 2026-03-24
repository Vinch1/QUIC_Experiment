# Project Plan

## Current Status

The project has moved from an empty scaffold to a runnable multi-node replication prototype with switchable inter-node transport.

## What Has Been Done

### 1. Project scaffolding

- Initialized the Go module and repository layout
- Added `cmd/node`, `cmd/client`, and `cmd/bench`
- Added core directories for `raft`, `transport`, `storage`, `statemachine`, and docs

### 2. Control plane

- Implemented a simple HTTP control API for:
  - `PUT /kv`
  - `GET /kv`
  - `GET /status`
- Implemented a CLI client that can:
  - write a key
  - read a key
  - inspect node status

### 3. Node-to-node transport

- Defined a shared inter-node transport interface
- Implemented a real TCP transport
- Implemented a real QUIC transport using `quic-go`
- Added framed request/response messaging for node-to-node RPCs

### 4. Minimal distributed replication flow

- Split node addresses into:
  - `control-addr` for client access
  - `raft-addr` for inter-node traffic
- Implemented a static leader model
- Implemented majority-ack write replication
- Implemented follower read-through to leader on local miss

### 5. Reliability and usability improvements

- Fixed startup readiness so the control API is actually listening before startup is reported as successful
- Added client retry logic for short startup races
- Added `make stop-cluster` to clean stale node processes and free ports
- Improved bind error messages for easier diagnosis

### 6. Validation

- `make fmt` passes
- `make test` passes
- `make build` passes
- Verified 3-node QUIC mode:
  - write through leader succeeds
  - read from follower succeeds
- Verified the same code structure supports TCP / QUIC switching through configuration

## What This Means

At this point, the project proves:

- client-to-node communication works
- node-to-node communication works over both TCP and QUIC
- multi-node replication works in a minimal majority-based model
- QUIC can carry the current replication RPC flow successfully

This is already enough for basic transport-level experimentation, but it is not yet a full Raft implementation.

## What Is Not Done Yet

- real Raft leader election
- Raft terms and voting
- log entries with index / term
- `AppendEntries` semantics
- conflict resolution and log repair
- commit index / applied index
- persistent WAL / recovery
- snapshot and membership changes
- benchmark automation and metrics dashboard

## Next Step

The next step should be:

## Phase 1: Minimal Correct Raft Core

### Goal

Replace the current static-leader replication demo with a real minimal Raft core.

### Priority tasks

1. Add Raft message types
   - `RequestVote`
   - `RequestVoteResponse`
   - `AppendEntries`
   - `AppendEntriesResponse`

2. Add persistent Raft state model
   - current term
   - voted-for
   - in-memory log entries

3. Implement follower / candidate / leader transitions
   - election timeout
   - heartbeat loop
   - vote counting

4. Replace direct `replicate_set` write with log replication
   - leader appends entry to local log
   - leader sends `AppendEntries`
   - followers validate and append
   - leader commits on majority
   - state machine applies only committed entries

5. Replace follower read-through hack with Raft-aware reads
   - leader reads from committed state
   - follower redirects to leader, or later supports lease/read-index

### Acceptance criteria

- 3 nodes can elect a leader automatically
- client writes succeed after election completes
- leader failure triggers a new election
- committed writes remain readable after leader change
- TCP and QUIC can still be swapped without changing Raft core logic

## Recommended Immediate Work Order

1. Define transport message schema for `RequestVote` and `AppendEntries`
2. Add in-memory log structure and term/vote state
3. Implement election timeout and heartbeat loop
4. Implement leader commit flow
5. Update client write/read semantics to follow the new Raft model

