package main

import (
	"flag"
	"log"

	"github.com/kwshr/mini-k8/internal/docker"
	"github.com/kwshr/mini-k8/internal/worker"
)

func main() {
	addr := flag.String("addr", ":9090", "address to listen on")
	flag.Parse()

	runner := docker.NewExecRunner()
	w := worker.New(runner)

	log.Fatal(w.Start(*addr))
}
