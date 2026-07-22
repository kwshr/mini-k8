package main

import (
	"flag"
	"log"

	"github.com/kwshr/mini-k8/internal/controlplane"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	workerURL := flag.String("worker-url", "http://localhost:9090", "worker base URL")
	flag.Parse()

	cp := controlplane.New(*workerURL)

	log.Fatal(cp.Start(*addr))
}
