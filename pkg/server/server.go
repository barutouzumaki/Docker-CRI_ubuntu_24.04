package server

import (
	"k8s.io/cri-api/pkg/apis/runtime/v1"
	"google.golang.org/grpc"

	"github.com/my-experiments/docker-cri/pkg/cri"
)

// RegisterServices registers CRI services with the gRPC server
func RegisterServices(grpcServer *grpc.Server, runtimeService *cri.RuntimeService, imageService *cri.ImageService) {
	v1.RegisterRuntimeServiceServer(grpcServer, runtimeService)
	v1.RegisterImageServiceServer(grpcServer, imageService)
}

