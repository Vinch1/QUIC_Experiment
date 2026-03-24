package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/leo/quic-raft/internal/cluster"
	"github.com/leo/quic-raft/internal/control/httpapi"
	"github.com/leo/quic-raft/internal/raft"
)

func main() {
	var (
		id          = flag.String("id", "node-1", "node id")
		controlAddr = flag.String("control-addr", "127.0.0.1:9001", "control api address")
		raftAddr    = flag.String("raft-addr", "127.0.0.1:7001", "raft transport address")
		transport   = flag.String("transport", "tcp", "transport kind: tcp|quic")
		peersRaw    = flag.String("peers", "", "comma separated peers: node-2=127.0.0.1:7002,node-3=127.0.0.1:7003")
		leader      = flag.Bool("leader", false, "bootstrap this node as leader")
		leaderID    = flag.String("leader-id", "", "known leader node id")
	)

	flag.Parse()

	peers, err := cluster.ParsePeers(*peersRaw)
	if err != nil {
		log.Fatalf("parse peers: %v", err)
	}

	cfg := raft.Config{
		NodeID:          *id,
		ControlAddr:     *controlAddr,
		RaftAddr:        *raftAddr,
		TransportKind:   *transport,
		Peers:           peers,
		BootstrapLeader: *leader,
		LeaderID:        *leaderID,
	}

	node, err := raft.NewNode(cfg)
	if err != nil {
		log.Fatalf("create node: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := node.Start(ctx); err != nil {
		log.Fatalf("start node: %v", formatBindError(err, *raftAddr, *controlAddr))
	}

	controlListener, err := net.Listen("tcp", *controlAddr)
	if err != nil {
		log.Fatalf("listen control api: %v", formatBindError(err, *raftAddr, *controlAddr))
	}

	server := &http.Server{
		Handler: httpapi.NewHandler(node),
	}

	go func() {
		if err := server.Serve(controlListener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	fmt.Printf("node started: %+v\n", node.Status())
	fmt.Printf("control api listening on http://%s\n", *controlAddr)
	fmt.Printf("raft transport listening on %s://%s\n", *transport, *raftAddr)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown server: %v", err)
	}

	if err := node.Stop(context.Background()); err != nil {
		log.Fatalf("stop node: %v", err)
	}
}

func formatBindError(err error, raftAddr, controlAddr string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, syscall.EADDRINUSE) || strings.Contains(err.Error(), "address already in use") {
		return fmt.Errorf("%w\nhint: another node is already using %s or %s; run `make stop-cluster` or inspect with `lsof -nP -iUDP:%s -iTCP:%s`", err, raftAddr, controlAddr, strings.Split(raftAddr, ":")[1], strings.Split(controlAddr, ":")[1])
	}

	return err
}
