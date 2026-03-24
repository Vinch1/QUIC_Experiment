package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	root "github.com/leo/quic-raft/internal/transport"
)

type clientConn struct {
	conn net.Conn
	mu   sync.Mutex
}

type Transport struct {
	addr    string
	timeout time.Duration

	mu      sync.RWMutex
	handler root.HandlerFunc
	ln      net.Listener
	clients map[string]*clientConn
}

func New(addr string) *Transport {
	return &Transport{
		addr:    addr,
		timeout: 3 * time.Second,
		clients: make(map[string]*clientConn),
	}
}

func (t *Transport) Start(_ context.Context, handler root.HandlerFunc) error {
	ln, err := net.Listen("tcp", t.addr)
	if err != nil {
		return fmt.Errorf("listen tcp %s: %w", t.addr, err)
	}

	t.mu.Lock()
	t.handler = handler
	t.ln = ln
	t.mu.Unlock()

	go t.acceptLoop()
	return nil
}

func (t *Transport) Stop(context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ln != nil {
		_ = t.ln.Close()
		t.ln = nil
	}

	for target, client := range t.clients {
		_ = client.conn.Close()
		delete(t.clients, target)
	}
	return nil
}

func (t *Transport) Send(ctx context.Context, target string, request root.Request) (root.Response, error) {
	client, err := t.getOrCreateClient(target)
	if err != nil {
		return root.Response{}, err
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	if deadline, ok := ctx.Deadline(); ok {
		_ = client.conn.SetDeadline(deadline)
	} else {
		_ = client.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	if err := root.WriteFrame(client.conn, request); err != nil {
		t.dropClient(target)
		return root.Response{}, fmt.Errorf("write request to %s: %w", target, err)
	}

	var response root.Response
	if err := root.ReadFrame(client.conn, &response); err != nil {
		t.dropClient(target)
		return root.Response{}, fmt.Errorf("read response from %s: %w", target, err)
	}

	return response, nil
}

func (t *Transport) Addr() string {
	return t.addr
}

func (t *Transport) Kind() root.Kind {
	return root.KindTCP
}

func (t *Transport) acceptLoop() {
	for {
		t.mu.RLock()
		ln := t.ln
		handler := t.handler
		t.mu.RUnlock()
		if ln == nil {
			return
		}

		conn, err := ln.Accept()
		if err != nil {
			t.mu.RLock()
			stopped := t.ln == nil
			t.mu.RUnlock()
			if stopped {
				return
			}
			continue
		}

		go t.handleConn(conn, handler)
	}
}

func (t *Transport) handleConn(conn net.Conn, handler root.HandlerFunc) {
	defer conn.Close()

	for {
		var request root.Request
		if err := root.ReadFrame(conn, &request); err != nil {
			return
		}

		response, err := handler(context.Background(), request)
		if err != nil {
			response = root.Response{
				OK:      false,
				Message: err.Error(),
			}
		}

		if err := root.WriteFrame(conn, response); err != nil {
			return
		}
	}
}

func (t *Transport) getOrCreateClient(target string) (*clientConn, error) {
	t.mu.RLock()
	existing := t.clients[target]
	t.mu.RUnlock()
	if existing != nil {
		return existing, nil
	}

	conn, err := net.DialTimeout("tcp", target, t.timeout)
	if err != nil {
		return nil, fmt.Errorf("dial tcp %s: %w", target, err)
	}

	client := &clientConn{conn: conn}

	t.mu.Lock()
	defer t.mu.Unlock()
	if current := t.clients[target]; current != nil {
		_ = conn.Close()
		return current, nil
	}
	t.clients[target] = client
	return client, nil
}

func (t *Transport) dropClient(target string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	client := t.clients[target]
	if client == nil {
		return
	}

	_ = client.conn.Close()
	delete(t.clients, target)
}
