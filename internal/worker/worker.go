package worker

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kwshr/mini-k8/internal/docker"
	"github.com/kwshr/mini-k8/internal/models"
)

// Worker listens for deployment requests from the Control Plane
// and executes them via a docker.Runner.
type Worker struct {
	runner docker.Runner
}

// New creates a Worker with the given docker.Runner.
func New(runner docker.Runner) *Worker {
	return &Worker{runner: runner}
}

// HandleRun is the HTTP handler for POST /run.
// It decodes the deploy request, runs the container, and returns the result.
func (w *Worker) HandleRun(rw http.ResponseWriter, r *http.Request) {
	// Only accept POST
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req models.DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(rw, http.StatusBadRequest, models.DeployResponse{
			Status: "error",
			Error:  "invalid request body: " + err.Error(),
		})
		return
	}

	log.Printf("[Worker] Received run request for image: %s", req.Image)

	// Execute docker run
	containerID, err := w.runner.Run(req.Image)
	if err != nil {
		log.Printf("[Worker] Docker run failed: %v", err)
		writeJSON(rw, http.StatusInternalServerError, models.DeployResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	log.Printf("[Worker] Container started: %s", containerID)
	writeJSON(rw, http.StatusOK, models.DeployResponse{
		Status:      "success",
		ContainerID: containerID,
	})
}

// Start begins listening for requests on the given address.
func (w *Worker) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/run", w.HandleRun)

	log.Printf("[Worker] Listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

// writeJSON is a helper that writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
