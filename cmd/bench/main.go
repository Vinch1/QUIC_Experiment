package main

import (
	"flag"
	"fmt"
)

func main() {
	var (
		scenario  = flag.String("scenario", "baseline", "scenario name")
		nodes     = flag.Int("nodes", 3, "node count")
		transport = flag.String("transport", "tcp", "transport kind")
	)

	flag.Parse()

	fmt.Printf("benchmark scaffold ready: scenario=%s nodes=%d transport=%s\n", *scenario, *nodes, *transport)
	fmt.Println("next step: implement runner, workload and metrics collection")
}
