package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"mini-k8/internal/api"
)

func main() {

	fmt.Println("========================================")
	fmt.Println("        mini-k8 CLI")
	fmt.Println("========================================")

	// Validate command-line arguments
	if len(os.Args) < 3 {
		fmt.Println("[CLI] Error: Invalid command")
		fmt.Println("[CLI] Usage: mini-k8 deploy <image>")
		return
	}

	command := os.Args[1]
	image := os.Args[2]

	fmt.Printf("[CLI] Command Received : %s\n", command)
	fmt.Printf("[CLI] Docker Image     : %s\n", image)

	if command != "deploy" {
		fmt.Printf("[CLI] Unsupported command: %s\n", command)
		fmt.Println("[CLI] Supported command: deploy")
		return
	}

	// Create deployment request
	req := api.DeployRequest{
		Image: image,
	}

	fmt.Println("[CLI] Creating deployment request...")

	body, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("[CLI] Failed to create JSON request: %v\n", err)
		return
	}

	fmt.Printf("[CLI] Request Body: %s\n", string(body))

	fmt.Println("[CLI] Sending deployment request to Control Plane...")
	fmt.Println("[CLI] Endpoint: http://localhost:8080/deploy")

	resp, err := http.Post(
		"http://localhost:8080/deploy",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		fmt.Printf("[CLI] Failed to contact Control Plane: %v\n", err)
		return
	}

	defer resp.Body.Close()

	fmt.Printf("[CLI] HTTP Response Status: %s\n", resp.Status)

	var result api.DeployResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("[CLI] Failed to decode response: %v\n", err)
		return
	}

	fmt.Println("----------------------------------------")
	fmt.Println("[CLI] Deployment Result")
	fmt.Println("----------------------------------------")
	fmt.Printf("Status : %s\n", result.Status)

	if result.Error != "" {
		fmt.Printf("Error  : %s\n", result.Error)
	} else {
		fmt.Println("Message: Container deployed successfully!")
	}

	fmt.Println("========================================")
}