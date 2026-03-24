package raft

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/leo/quic-raft/internal/cluster"
	"github.com/leo/quic-raft/internal/statemachine/kv"
	roottransport "github.com/leo/quic-raft/internal/transport"
	quictransport "github.com/leo/quic-raft/internal/transport/quic"
	tcptransport "github.com/leo/quic-raft/internal/transport/tcp"
)

const replicationTimeout = 2 * time.Second
const readThroughTimeout = 2 * time.Second

type Node struct {
	mu        sync.RWMutex
	config    Config
	state     State
	started   bool
	store     *kv.Store
	transport roottransport.Transport
	peerAddrs map[string]string
}

func NewNode(cfg Config) (*Node, error) {
	return newNode(cfg, nil)
}

func NewNodeWithTransport(cfg Config, transport roottransport.Transport) (*Node, error) {
	return newNode(cfg, transport)
}

func newNode(cfg Config, transport roottransport.Transport) (*Node, error) {
	cfg = cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	node := &Node{
		config:    cfg,
		state:     StateFollower,
		store:     kv.New(),
		transport: transport,
		peerAddrs: make(map[string]string, len(cfg.Peers)),
	}

	for _, peer := range cfg.Peers {
		node.peerAddrs[peer.ID] = peer.Address
	}

	if cfg.BootstrapLeader || len(cfg.Peers) == 0 {
		node.state = StateLeader
	}

	return node, nil
}

func (n *Node) Start(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.started {
		return nil
	}

	if n.transport == nil {
		transport, err := newTransport(n.config.TransportKind, n.config.RaftAddr)
		if err != nil {
			return err
		}
		n.transport = transport
	}

	if err := n.transport.Start(ctx, n.handleTransportRequest); err != nil {
		return err
	}

	n.started = true
	return nil
}

func (n *Node) Stop(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.started {
		return nil
	}

	n.started = false
	if n.transport != nil {
		return n.transport.Stop(ctx)
	}
	return nil
}

func (n *Node) Propose(ctx context.Context, key, value string) error {
	n.mu.RLock()
	started := n.started
	state := n.state
	peers := append([]cluster.Node(nil), n.config.Peers...)
	n.mu.RUnlock()

	if !started {
		return fmt.Errorf("node is not running")
	}
	if state != StateLeader {
		return fmt.Errorf("node is not leader")
	}

	totalNodes := len(peers) + 1
	majority := totalNodes/2 + 1
	acks := 1

	request := roottransport.Request{
		Type:  roottransport.MessageReplicateSet,
		From:  n.config.NodeID,
		Key:   key,
		Value: value,
	}

	for _, peer := range peers {
		replicateCtx, cancel := context.WithTimeout(ctx, replicationTimeout)
		response, err := n.transport.Send(replicateCtx, peer.Address, request)
		cancel()
		if err != nil {
			continue
		}
		if response.OK {
			acks++
		}
	}

	if acks < majority {
		return fmt.Errorf("replication failed: received %d acks, need %d", acks, majority)
	}

	n.store.Put(key, value)
	return nil
}

func (n *Node) Get(key string) (string, bool, error) {
	n.mu.RLock()
	started := n.started
	state := n.state
	leaderID := n.config.LeaderID
	n.mu.RUnlock()

	if !started {
		return "", false, fmt.Errorf("node is not running")
	}

	value, ok := n.store.Get(key)
	if ok || state == StateLeader || leaderID == "" {
		return value, ok, nil
	}

	leaderAddr, ok := n.peerAddrs[leaderID]
	if !ok {
		return value, false, nil
	}

	response, err := n.fetchFromLeader(leaderAddr, key)
	if err != nil {
		return "", false, nil
	}
	if response.Found {
		n.store.Put(key, response.Value)
		return response.Value, true, nil
	}

	return value, ok, nil
}

func (n *Node) Status() Status {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return Status{
		NodeID:        n.config.NodeID,
		State:         n.state,
		ControlAddr:   n.config.ControlAddr,
		RaftAddr:      n.config.RaftAddr,
		TransportKind: n.config.TransportKind,
		LeaderID:      n.effectiveLeaderID(),
		PeerCount:     len(n.config.Peers),
		Started:       n.started,
	}
}

func (n *Node) IsLeader() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.state == StateLeader
}

func (n *Node) LeaderID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.effectiveLeaderID()
}

func (n *Node) handleTransportRequest(_ context.Context, request roottransport.Request) (roottransport.Response, error) {
	switch request.Type {
	case roottransport.MessagePing:
		return roottransport.Response{
			OK:   true,
			From: n.config.NodeID,
		}, nil
	case roottransport.MessageReplicateSet:
		if request.Key == "" {
			return roottransport.Response{
				OK:      false,
				Message: "key is required",
				From:    n.config.NodeID,
			}, nil
		}

		n.store.Put(request.Key, request.Value)
		return roottransport.Response{
			OK:      true,
			Message: "replicated",
			From:    n.config.NodeID,
		}, nil
	case roottransport.MessageFetchValue:
		if request.Key == "" {
			return roottransport.Response{
				OK:      false,
				Message: "key is required",
				From:    n.config.NodeID,
			}, nil
		}

		value, found := n.store.Get(request.Key)
		return roottransport.Response{
			OK:    true,
			From:  n.config.NodeID,
			Value: value,
			Found: found,
		}, nil
	default:
		return roottransport.Response{
			OK:      false,
			Message: "unsupported transport message",
			From:    n.config.NodeID,
		}, nil
	}
}

func newTransport(kind, addr string) (roottransport.Transport, error) {
	switch kind {
	case string(roottransport.KindTCP):
		return tcptransport.New(addr), nil
	case string(roottransport.KindQUIC):
		return quictransport.New(addr)
	default:
		return nil, fmt.Errorf("unsupported transport kind: %s", kind)
	}
}

func (n *Node) effectiveLeaderID() string {
	if n.state == StateLeader {
		return n.config.NodeID
	}
	if n.config.LeaderID != "" {
		return n.config.LeaderID
	}
	return ""
}

func (n *Node) fetchFromLeader(leaderAddr, key string) (roottransport.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), readThroughTimeout)
	defer cancel()

	return n.transport.Send(ctx, leaderAddr, roottransport.Request{
		Type: roottransport.MessageFetchValue,
		From: n.config.NodeID,
		Key:  key,
	})
}
