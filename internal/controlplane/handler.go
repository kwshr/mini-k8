package controlplane

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"mini-k8/internal/api"
)

func DeployHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("========================================")
	log.Println("[CONTROL PLANE] Deployment request received")

	var req api.DeployRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[CONTROL PLANE] Invalid request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[CONTROL PLANE] Image: %s\n", req.Image)

	body, _ := json.Marshal(req)

	log.Println("[CONTROL PLANE] Forwarding request to Worker")

	resp, err := http.Post(
		"http://localhost:9001/run",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {

		log.Printf("[CONTROL PLANE] Worker unreachable: %v\n", err)

		http.Error(w, err.Error(), 500)

		return
	}

	defer resp.Body.Close()

	log.Println("[CONTROL PLANE] Response received from Worker")

	var response api.DeployResponse

	json.NewDecoder(resp.Body).Decode(&response)

	log.Printf("[CONTROL PLANE] Worker Status: %s\n", response.Status)

	if response.Error != "" {
		log.Printf("[CONTROL PLANE] Worker Error: %s\n", response.Error)
	}

	log.Println("[CONTROL PLANE] Sending response back to CLI")

	json.NewEncoder(w).Encode(response)
}