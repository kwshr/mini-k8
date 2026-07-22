package controlplane

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mini-k8/mini-k8/internal/proto"
)

// Start launches the REST control plane server.
func Start(port, workerAddress string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /deploy", func(w http.ResponseWriter, r *http.Request) {
		handleDeploy(w, r, workerAddress)
	})

	addr := fmt.Sprintf("127.0.0.1:%s", port)
	fmt.Printf("Control plane listening on http://%s\n", addr)
	fmt.Printf("Forwarding deployments to %s\n", workerAddress)

	return http.ListenAndServe(addr, mux)
}

func handleDeploy(w http.ResponseWriter, r *http.Request, workerAddress string) {
	w.Header().Set("Content-Type", "application/json")

	// Read and parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "failure",
			"message": "failed to read request body",
		})
		return
	}

	var req map[string]string
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "failure",
			"message": "invalid JSON in request body",
		})
		return
	}

	image := strings.TrimSpace(req["image"])
	if image == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "failure",
			"message": "image is required",
		})
		return
	}

	// Connect to worker via gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, workerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "failure",
			"message": fmt.Sprintf("failed to connect to worker: %v", err),
		})
		return
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)
	resp, err := client.Deploy(ctx, &pb.DeployRequest{Image: image})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "failure",
			"message": fmt.Sprintf("worker deployment failed: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":      resp.Status,
		"image":       image,
		"containerId": resp.ContainerId,
		"message":     resp.Message,
	})
}
