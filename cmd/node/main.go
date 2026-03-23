package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leo/quic-raft/internal/control/httpapi"
	"github.com/leo/quic-raft/internal/raft"
)

func main() {
	var (
		id        = flag.String("id", "node-1", "node id")
		listen    = flag.String("listen", "127.0.0.1:9001", "listen address")
		transport = flag.String("transport", "tcp", "transport kind: tcp|quic")
	)

	flag.Parse()

	cfg := raft.Config{
		NodeID:        *id,
		ListenAddr:    *listen,
		TransportKind: *transport,
		Peers:         nil,
	}

	node, err := raft.NewNode(cfg)
	if err != nil {
		log.Fatalf("create node: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := node.Start(ctx); err != nil {
		log.Fatalf("start node: %v", err)
	}

	server := &http.Server{
		Addr:    *listen,
		Handler: httpapi.NewHandler(node),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	fmt.Printf("node started: %+v\n", node.Status())
	fmt.Printf("control api listening on http://%s\n", *listen)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3_000_000_000)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown server: %v", err)
	}

	if err := node.Stop(context.Background()); err != nil {
		log.Fatalf("stop node: %v", err)
	}
}
