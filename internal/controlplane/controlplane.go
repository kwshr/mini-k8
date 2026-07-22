package controlplane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kwshr/mini-k8/internal/models"
)

// ControlPlane is the Façade that sits between the user and the worker(s).
// It validates deployment requests and forwards them to a worker.
type ControlPlane struct {
	workerURL string // e.g. "http://localhost:9090"
}

// New creates a ControlPlane that forwards requests to the given worker URL.
func New(workerURL string) *ControlPlane {
	return &ControlPlane{workerURL: workerURL}
}

// HandleDeploy is the HTTP handler for POST /deploy.
func (cp *ControlPlane) HandleDeploy(w http.ResponseWriter, r *http.Request) {
	// Only accept POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request
	var req models.DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.DeployResponse{
			Status: "error",
			Error:  "invalid request body: " + err.Error(),
		})
		return
	}

	// Validate — image must not be empty
	if req.Image == "" {
		writeJSON(w, http.StatusBadRequest, models.DeployResponse{
			Status: "error",
			Error:  "image name is required",
		})
		return
	}

	log.Printf("[ControlPlane] Received deploy request for image: %s", req.Image)

	// Forward to worker
	resp, err := cp.forwardToWorker(req)
	if err != nil {
		log.Printf("[ControlPlane] Failed to forward to worker: %v", err)
		writeJSON(w, http.StatusBadGateway, models.DeployResponse{
			Status: "error",
			Error:  "worker communication failed: " + err.Error(),
		})
		return
	}

	log.Printf("[ControlPlane] Deployment result: status=%s container_id=%s", resp.Status, resp.ContainerID)
	statusCode := http.StatusOK
	if resp.Status != "success" {
		statusCode = http.StatusInternalServerError
	}
	writeJSON(w, statusCode, resp)
}

// forwardToWorker sends the deploy request to the worker's /run endpoint.
func (cp *ControlPlane) forwardToWorker(req models.DeployRequest) (*models.DeployResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := cp.workerURL + "/run"
	httpResp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("POST %s: %w", url, err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var resp models.DeployResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &resp, nil
}

// Start begins listening for requests on the given address.
func (cp *ControlPlane) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/deploy", cp.HandleDeploy)

	log.Printf("[ControlPlane] Listening on %s (worker: %s)", addr, cp.workerURL)
	return http.ListenAndServe(addr, mux)
}

// writeJSON is a helper that writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
