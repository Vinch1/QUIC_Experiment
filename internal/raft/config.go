package raft

import (
	"fmt"
	"time"
)

const (
	DefaultElectionTimeout  = 1500 * time.Millisecond
	DefaultHeartbeatTimeout = 300 * time.Millisecond
)

type Config struct {
	NodeID            string
	ListenAddr        string
	TransportKind     string
	Peers             []string
	ElectionTimeout   time.Duration
	HeartbeatInterval time.Duration
}

func (c Config) Normalize() Config {
	if c.ElectionTimeout == 0 {
		c.ElectionTimeout = DefaultElectionTimeout
	}
	if c.HeartbeatInterval == 0 {
		c.HeartbeatInterval = DefaultHeartbeatTimeout
	}
	return c
}

func (c Config) Validate() error {
	if c.NodeID == "" {
		return fmt.Errorf("node id is required")
	}
	if c.ListenAddr == "" {
		return fmt.Errorf("listen address is required")
	}
	switch c.TransportKind {
	case "tcp", "quic":
	default:
		return fmt.Errorf("unsupported transport kind: %s", c.TransportKind)
	}
	return nil
}
