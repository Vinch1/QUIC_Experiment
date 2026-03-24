package cluster

import (
	"fmt"
	"strings"
)

type Node struct {
	ID      string
	Address string
}

type Topology struct {
	Nodes []Node
}

func (t Topology) Size() int {
	return len(t.Nodes)
}

func ParsePeers(raw string) ([]Node, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	items := strings.Split(raw, ",")
	peers := make([]Node, 0, len(items))
	seen := make(map[string]struct{}, len(items))

	for _, item := range items {
		item = strings.TrimSpace(item)
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid peer entry %q, expected node-id=host:port", item)
		}

		id := strings.TrimSpace(parts[0])
		addr := strings.TrimSpace(parts[1])
		if id == "" || addr == "" {
			return nil, fmt.Errorf("invalid peer entry %q, expected node-id=host:port", item)
		}
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("duplicate peer id %q", id)
		}
		seen[id] = struct{}{}

		peers = append(peers, Node{
			ID:      id,
			Address: addr,
		})
	}

	return peers, nil
}
