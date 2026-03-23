package httpapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leo/quic-raft/internal/raft"
)

func TestHandlerPutThenGet(t *testing.T) {
	node, err := raft.NewNode(raft.Config{
		NodeID:        "node-1",
		ListenAddr:    "127.0.0.1:9001",
		TransportKind: "tcp",
	})
	if err != nil {
		t.Fatalf("new node: %v", err)
	}

	if err := node.Start(context.Background()); err != nil {
		t.Fatalf("start node: %v", err)
	}

	handler := NewHandler(node)

	putReq := httptest.NewRequest(http.MethodPost, "/kv", bytes.NewBufferString(`{"command":"put","key":"demo","value":"hello"}`))
	putRec := httptest.NewRecorder()
	handler.ServeHTTP(putRec, putReq)

	if putRec.Code != http.StatusOK {
		t.Fatalf("unexpected put status: %d", putRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/kv?key=demo", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("unexpected get status: %d", getRec.Code)
	}
}
