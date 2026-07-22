package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kwshr/mini-k8/internal/models"
)

const defaultControlPlaneURL = "http://localhost:8080"

func main() {
	if len(os.Args) < 3 || os.Args[1] != "deploy" {
		fmt.Fprintf(os.Stderr, "Usage: mini-kube deploy <image>\n")
		os.Exit(1)
	}

	image := os.Args[2]

	cpURL := os.Getenv("CONTROL_PLANE_URL")
	if cpURL == "" {
		cpURL = defaultControlPlaneURL
	}

	req := models.DeployRequest{Image: image}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	url := cpURL + "/deploy"
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not reach Control Plane at %s — %v\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read response — %v\n", err)
		os.Exit(1)
	}

	var result models.DeployResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid response — %v\n", err)
		os.Exit(1)
	}

	if result.Status == "success" {
		fmt.Printf("✅ Deployment successful!\n")
		fmt.Printf("   Image:        %s\n", image)
		fmt.Printf("   Container ID: %s\n", result.ContainerID)
	} else {
		fmt.Printf("❌ Deployment failed!\n")
		fmt.Printf("   Image: %s\n", image)
		fmt.Printf("   Error: %s\n", result.Error)
		os.Exit(1)
	}
}
