package cluster

import "testing"

func TestParsePeers(t *testing.T) {
	peers, err := ParsePeers("node-2=127.0.0.1:7002,node-3=127.0.0.1:7003")
	if err != nil {
		t.Fatalf("parse peers: %v", err)
	}
	if len(peers) != 2 {
		t.Fatalf("unexpected peer count: %d", len(peers))
	}
	if peers[0].ID != "node-2" || peers[1].Address != "127.0.0.1:7003" {
		t.Fatalf("unexpected peers: %+v", peers)
	}
}
