package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	root "github.com/leo/quic-raft/internal/transport"
	libquic "github.com/quic-go/quic-go"
)

type peerConn struct {
	conn *libquic.Conn
}

type Transport struct {
	addr       string
	timeout    time.Duration
	serverTLS  *tls.Config
	clientTLS  *tls.Config
	quicConfig *libquic.Config

	mu       sync.RWMutex
	handler  root.HandlerFunc
	listener *libquic.Listener
	clients  map[string]*peerConn
}

func New(addr string) (*Transport, error) {
	serverTLS, clientTLS, err := generateTLSConfigs()
	if err != nil {
		return nil, err
	}

	return &Transport{
		addr:       addr,
		timeout:    3 * time.Second,
		serverTLS:  serverTLS,
		clientTLS:  clientTLS,
		quicConfig: &libquic.Config{KeepAlivePeriod: 10 * time.Second},
		clients:    make(map[string]*peerConn),
	}, nil
}

func (t *Transport) Start(ctx context.Context, handler root.HandlerFunc) error {
	listener, err := libquic.ListenAddr(t.addr, t.serverTLS, t.quicConfig)
	if err != nil {
		return fmt.Errorf("listen quic %s: %w", t.addr, err)
	}

	t.mu.Lock()
	t.handler = handler
	t.listener = listener
	t.mu.Unlock()

	go t.acceptLoop(ctx)
	return nil
}

func (t *Transport) Stop(context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.listener != nil {
		_ = t.listener.Close()
		t.listener = nil
	}

	for target, client := range t.clients {
		_ = client.conn.CloseWithError(0, "shutdown")
		delete(t.clients, target)
	}
	return nil
}

func (t *Transport) Send(ctx context.Context, target string, request root.Request) (root.Response, error) {
	conn, err := t.getOrCreateConn(ctx, target)
	if err != nil {
		return root.Response{}, err
	}

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		t.dropConn(target)
		return root.Response{}, fmt.Errorf("open quic stream to %s: %w", target, err)
	}
	defer stream.Close()

	if err := root.WriteFrame(stream, request); err != nil {
		t.dropConn(target)
		return root.Response{}, fmt.Errorf("write request to %s: %w", target, err)
	}

	var response root.Response
	if err := root.ReadFrame(stream, &response); err != nil {
		t.dropConn(target)
		return root.Response{}, fmt.Errorf("read response from %s: %w", target, err)
	}

	return response, nil
}

func (t *Transport) Addr() string {
	return t.addr
}

func (t *Transport) Kind() root.Kind {
	return root.KindQUIC
}

func (t *Transport) acceptLoop(ctx context.Context) {
	for {
		t.mu.RLock()
		listener := t.listener
		handler := t.handler
		t.mu.RUnlock()
		if listener == nil {
			return
		}

		conn, err := listener.Accept(ctx)
		if err != nil {
			t.mu.RLock()
			stopped := t.listener == nil
			t.mu.RUnlock()
			if stopped {
				return
			}
			continue
		}

		go t.handleConn(conn, handler)
	}
}

func (t *Transport) handleConn(conn *libquic.Conn, handler root.HandlerFunc) {
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go func() {
			defer stream.Close()

			var request root.Request
			if err := root.ReadFrame(stream, &request); err != nil {
				return
			}

			response, err := handler(context.Background(), request)
			if err != nil {
				response = root.Response{
					OK:      false,
					Message: err.Error(),
				}
			}

			_ = root.WriteFrame(stream, response)
		}()
	}
}

func (t *Transport) getOrCreateConn(ctx context.Context, target string) (*libquic.Conn, error) {
	t.mu.RLock()
	existing := t.clients[target]
	t.mu.RUnlock()
	if existing != nil {
		return existing.conn, nil
	}

	dialCtx := ctx
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		dialCtx, cancel = context.WithTimeout(ctx, t.timeout)
		defer cancel()
	}

	conn, err := libquic.DialAddr(dialCtx, target, t.clientTLS, t.quicConfig)
	if err != nil {
		return nil, fmt.Errorf("dial quic %s: %w", target, err)
	}

	client := &peerConn{conn: conn}

	t.mu.Lock()
	defer t.mu.Unlock()
	if current := t.clients[target]; current != nil {
		_ = conn.CloseWithError(0, "duplicate connection")
		return current.conn, nil
	}
	t.clients[target] = client
	return client.conn, nil
}

func (t *Transport) dropConn(target string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	client := t.clients[target]
	if client == nil {
		return
	}

	_ = client.conn.CloseWithError(0, "reset")
	delete(t.clients, target)
}
