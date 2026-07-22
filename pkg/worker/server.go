package worker

import (
	"context"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/mini-k8/mini-k8/internal/proto"
)

// Start launches the gRPC worker server.
func Start(port string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	server := grpc.NewServer()
	pb.RegisterWorkerServiceServer(server, &workerServer{})

	fmt.Printf("Worker listening on 0.0.0.0:%s\n", port)

	return server.Serve(listener)
}

type workerServer struct {
	pb.UnimplementedWorkerServiceServer
}

func (s *workerServer) Deploy(ctx context.Context, req *pb.DeployRequest) (*pb.DeployResponse, error) {
	image := strings.TrimSpace(req.Image)
	if image == "" {
		return &pb.DeployResponse{
			Status:  "failure",
			Message: "image is required",
		}, fmt.Errorf("image is required")
	}

	containerID, err := RunContainer(image)
	if err != nil {
		return &pb.DeployResponse{
			Status:  "failure",
			Message: err.Error(),
		}, err
	}

	return &pb.DeployResponse{
		Status:      "success",
		ContainerId: containerID,
		Message:     fmt.Sprintf("container started from %s", image),
	}, nil
}
