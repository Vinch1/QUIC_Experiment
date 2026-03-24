package raft

import (
	"context"
	"testing"

	"github.com/leo/quic-raft/internal/cluster"
	roottransport "github.com/leo/quic-raft/internal/transport"
)

type fakeTransport struct{}

func (fakeTransport) Start(context.Context, roottransport.HandlerFunc) error {
	return nil
}

func (fakeTransport) Stop(context.Context) error {
	return nil
}

func (fakeTransport) Send(context.Context, string, roottransport.Request) (roottransport.Response, error) {
	return roottransport.Response{OK: true}, nil
}

func (fakeTransport) Addr() string {
	return "in-memory"
}

func (fakeTransport) Kind() roottransport.Kind {
	return roottransport.KindTCP
}

func TestConfigNormalizeAndValidate(t *testing.T) {
	cfg := Config{
		NodeID:        "node-1",
		ControlAddr:   "127.0.0.1:9001",
		RaftAddr:      "127.0.0.1:7001",
		TransportKind: "tcp",
	}

	cfg = cfg.Normalize()

	if cfg.ElectionTimeout == 0 {
		t.Fatal("expected election timeout to be set")
	}
	if cfg.HeartbeatInterval == 0 {
		t.Fatal("expected heartbeat interval to be set")
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("validate config: %v", err)
	}
}

func TestConfigValidateRejectsUnknownTransport(t *testing.T) {
	cfg := Config{
		NodeID:        "node-1",
		ControlAddr:   "127.0.0.1:9001",
		RaftAddr:      "127.0.0.1:7001",
		TransportKind: "udp",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for unsupported transport")
	}
}

func TestNodeProposalLifecycle(t *testing.T) {
	node, err := NewNodeWithTransport(Config{
		NodeID:          "node-1",
		ControlAddr:     "127.0.0.1:9001",
		RaftAddr:        "127.0.0.1:7001",
		TransportKind:   "tcp",
		BootstrapLeader: true,
	}, fakeTransport{})
	if err != nil {
		t.Fatalf("new node: %v", err)
	}

	if err := node.Start(t.Context()); err != nil {
		t.Fatalf("start node: %v", err)
	}

	if err := node.Propose(t.Context(), "demo", "hello"); err != nil {
		t.Fatalf("propose: %v", err)
	}

	value, ok, err := node.Get("demo")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "hello" {
		t.Fatalf("unexpected value: %s", value)
	}
}

type fakeReadThroughTransport struct{}

func (fakeReadThroughTransport) Start(context.Context, roottransport.HandlerFunc) error {
	return nil
}

func (fakeReadThroughTransport) Stop(context.Context) error {
	return nil
}

func (fakeReadThroughTransport) Send(_ context.Context, _ string, request roottransport.Request) (roottransport.Response, error) {
	if request.Type == roottransport.MessageFetchValue && request.Key == "demo" {
		return roottransport.Response{
			OK:    true,
			Found: true,
			Value: "hello",
		}, nil
	}
	return roottransport.Response{OK: true}, nil
}

func (fakeReadThroughTransport) Addr() string {
	return "in-memory"
}

func (fakeReadThroughTransport) Kind() roottransport.Kind {
	return roottransport.KindQUIC
}

func TestFollowerGetReadsThroughLeader(t *testing.T) {
	node, err := NewNodeWithTransport(Config{
		NodeID:        "node-2",
		ControlAddr:   "127.0.0.1:9002",
		RaftAddr:      "127.0.0.1:7002",
		TransportKind: "quic",
		LeaderID:      "node-1",
		Peers: []cluster.Node{
			{ID: "node-1", Address: "127.0.0.1:7001"},
		},
	}, fakeReadThroughTransport{})
	if err != nil {
		t.Fatalf("new node: %v", err)
	}

	if err := node.Start(t.Context()); err != nil {
		t.Fatalf("start node: %v", err)
	}

	value, ok, err := node.Get("demo")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !ok {
		t.Fatal("expected key to be read from leader")
	}
	if value != "hello" {
		t.Fatalf("unexpected value: %s", value)
	}
}
