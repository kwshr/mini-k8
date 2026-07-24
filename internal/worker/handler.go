package worker

import (
	"encoding/json"
	"log"
	"net/http"

	"mini-k8/internal/api"
	"mini-k8/internal/docker"
)

func RunHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("========================================")
	log.Println("[WORKER] Received deployment request")

	var req api.DeployRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[WORKER] Failed to decode request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[WORKER] Image Requested: %s\n", req.Image)

	log.Println("[WORKER] Starting Docker container...")

	err = docker.Run(req.Image)
	if err != nil {

		log.Printf("[WORKER] Docker failed: %v\n", err)

		json.NewEncoder(w).Encode(api.DeployResponse{
			Status: "failed",
			Error:  err.Error(),
		})

		return
	}

	log.Println("[WORKER] Docker container started successfully")
	log.Println("[WORKER] Sending success response to Control Plane")

	json.NewEncoder(w).Encode(api.DeployResponse{
		Status: "success",
	})
}