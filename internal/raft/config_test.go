package raft

import "testing"

func TestConfigNormalizeAndValidate(t *testing.T) {
	cfg := Config{
		NodeID:        "node-1",
		ListenAddr:    "127.0.0.1:9001",
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
		ListenAddr:    "127.0.0.1:9001",
		TransportKind: "udp",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for unsupported transport")
	}
}

func TestNodeProposalLifecycle(t *testing.T) {
	node, err := NewNode(Config{
		NodeID:        "node-1",
		ListenAddr:    "127.0.0.1:9001",
		TransportKind: "tcp",
	})
	if err != nil {
		t.Fatalf("new node: %v", err)
	}

	if err := node.Start(t.Context()); err != nil {
		t.Fatalf("start node: %v", err)
	}

	if err := node.Propose("demo", "hello"); err != nil {
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
