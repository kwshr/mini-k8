package main

import (
	"log"
	"net/http"

	"mini-k8/internal/controlplane"
)

func main() {

	http.HandleFunc("/deploy", controlplane.DeployHandler)

	log.Println("Control Plane running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}