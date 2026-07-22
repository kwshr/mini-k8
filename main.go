package main

import (
	"fmt"
	"os"

	"github.com/mini-k8/mini-k8/pkg/cli"
	"github.com/mini-k8/mini-k8/pkg/controlplane"
	"github.com/mini-k8/mini-k8/pkg/worker"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "deploy":
		if len(os.Args) < 3 {
			fmt.Println("Error: image name required")
			printUsage()
			os.Exit(1)
		}

		image := os.Args[2]
		controlPlaneURL := "http://127.0.0.1:8080"

		// Check for --control-plane flag
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--control-plane" && i+1 < len(os.Args) {
				controlPlaneURL = os.Args[i+1]
			}
		}

		if err := cli.Deploy(image, controlPlaneURL); err != nil {
			fmt.Fprintf(os.Stderr, "Deployment failed: %v\n", err)
			os.Exit(1)
		}

	case "control-plane":
		port := "8080"
		workerAddress := "127.0.0.1:50051"

		// Parse flags
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--port" && i+1 < len(os.Args) {
				port = os.Args[i+1]
			}
			if os.Args[i] == "--worker" && i+1 < len(os.Args) {
				workerAddress = os.Args[i+1]
			}
		}

		if err := controlplane.Start(port, workerAddress); err != nil {
			fmt.Fprintf(os.Stderr, "Control plane failed: %v\n", err)
			os.Exit(1)
		}

	case "worker":
		port := "50051"

		// Parse flags
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--port" && i+1 < len(os.Args) {
				port = os.Args[i+1]
			}
		}

		if err := worker.Start(port); err != nil {
			fmt.Fprintf(os.Stderr, "Worker failed: %v\n", err)
			os.Exit(1)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  mini-kube deploy <image> [--control-plane http://127.0.0.1:8080]")
	fmt.Println("  mini-kube control-plane [--port 8080] [--worker 127.0.0.1:50051]")
	fmt.Println("  mini-kube worker [--port 50051]")
}
