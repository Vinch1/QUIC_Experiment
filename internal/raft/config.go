package raft

import (
	"fmt"
	"time"

	"github.com/leo/quic-raft/internal/cluster"
)

const (
	DefaultElectionTimeout  = 1500 * time.Millisecond
	DefaultHeartbeatTimeout = 300 * time.Millisecond
)

type Config struct {
	NodeID            string
	ControlAddr       string
	RaftAddr          string
	TransportKind     string
	Peers             []cluster.Node
	BootstrapLeader   bool
	LeaderID          string
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
	if c.ControlAddr == "" {
		return fmt.Errorf("control address is required")
	}
	if c.RaftAddr == "" {
		return fmt.Errorf("raft address is required")
	}
	switch c.TransportKind {
	case "tcp", "quic":
	default:
		return fmt.Errorf("unsupported transport kind: %s", c.TransportKind)
	}

	for _, peer := range c.Peers {
		if peer.ID == c.NodeID {
			return fmt.Errorf("peer list must not include self node id %s", c.NodeID)
		}
	}

	return nil
}
