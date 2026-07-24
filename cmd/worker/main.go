package main

import (
	"log"
	"net/http"

	"mini-k8/internal/worker"
)

func main() {

	http.HandleFunc("/run", worker.RunHandler)

	log.Println("Worker running on :9001")

	log.Fatal(http.ListenAndServe(":9001", nil))
}