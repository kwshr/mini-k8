// Package proto contains the gRPC service definitions and messages for mini-k8.
package proto

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// DeployRequest is the request message for Worker.Deploy.
type DeployRequest struct {
	Image string
}

// DeployResponse is the response message for Worker.Deploy.
type DeployResponse struct {
	Status      string
	ContainerId string
	Message     string
}

// WorkerServiceServer defines the server interface for WorkerService.
type WorkerServiceServer interface {
	Deploy(context.Context, *DeployRequest) (*DeployResponse, error)
}

// WorkerServiceClient defines the client interface for WorkerService.
type WorkerServiceClient interface {
	Deploy(ctx context.Context, in *DeployRequest, opts ...grpc.CallOption) (*DeployResponse, error)
}

// UnimplementedWorkerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedWorkerServiceServer struct{}

func (UnimplementedWorkerServiceServer) Deploy(context.Context, *DeployRequest) (*DeployResponse, error) {
	return nil, fmt.Errorf("method Deploy not implemented")
}

// RegisterWorkerServiceServer registers a WorkerServiceServer with the gRPC server.
func RegisterWorkerServiceServer(s *grpc.Server, srv WorkerServiceServer) {
	s.RegisterService(&_WorkerService_serviceDesc, srv)
}

// NewWorkerServiceClient creates a new WorkerServiceClient.
func NewWorkerServiceClient(cc grpc.ClientConnInterface) WorkerServiceClient {
	return &workerServiceClient{cc}
}

type workerServiceClient struct {
	cc grpc.ClientConnInterface
}

func (c *workerServiceClient) Deploy(ctx context.Context, in *DeployRequest, opts ...grpc.CallOption) (*DeployResponse, error) {
	out := new(DeployResponse)
	err := c.cc.Invoke(ctx, "/mini_kube.WorkerService/Deploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _WorkerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mini_kube.WorkerService",
	HandlerType: (*WorkerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deploy",
			Handler:    _WorkerService_Deploy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/worker.proto",
}

func _WorkerService_Deploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).Deploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mini_kube.WorkerService/Deploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).Deploy(ctx, req.(*DeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}
