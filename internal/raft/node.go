package raft

import (
	"context"
	"fmt"
	"sync"

	"github.com/leo/quic-raft/internal/statemachine/kv"
)

type Node struct {
	mu      sync.RWMutex
	config  Config
	state   State
	started bool
	store   *kv.Store
}

func NewNode(cfg Config) (*Node, error) {
	cfg = cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Node{
		config: cfg,
		state:  StateFollower,
		store:  kv.New(),
	}, nil
}

func (n *Node) Start(_ context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.started = true
	if len(n.config.Peers) == 0 {
		n.state = StateLeader
	}
	return nil
}

func (n *Node) Stop(_ context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.started = false
	return nil
}

func (n *Node) Propose(key, value string) error {
	n.mu.RLock()
	started := n.started
	n.mu.RUnlock()

	if !started {
		return fmt.Errorf("node is not running")
	}

	n.store.Put(key, value)
	return nil
}

func (n *Node) Get(key string) (string, bool, error) {
	n.mu.RLock()
	started := n.started
	n.mu.RUnlock()

	if !started {
		return "", false, fmt.Errorf("node is not running")
	}

	value, ok := n.store.Get(key)
	return value, ok, nil
}

func (n *Node) Status() Status {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return Status{
		NodeID:        n.config.NodeID,
		State:         n.state,
		ListenAddr:    n.config.ListenAddr,
		TransportKind: n.config.TransportKind,
		PeerCount:     len(n.config.Peers),
		Started:       n.started,
	}
}
