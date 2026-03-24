package httpapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leo/quic-raft/internal/raft"
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

func TestHandlerPutThenGet(t *testing.T) {
	node, err := raft.NewNodeWithTransport(raft.Config{
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
