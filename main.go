package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/my-experiments/docker-cri/pkg/cri"
	"github.com/my-experiments/docker-cri/pkg/docker"
	"github.com/my-experiments/docker-cri/pkg/server"
)

var (
	socketPath = flag.String("socket", "/var/run/docker-cri.sock", "Path to the CRI socket")
	logLevel   = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	rootDir    = flag.String("root-dir", "/var/lib/docker-cri", "Root directory for storing runtime data")
)

func main() {
	flag.Parse()

	// Set log level
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level: %v", err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Infof("Starting Docker CRI shim...")
	logrus.Infof("Socket path: %s", *socketPath)
	logrus.Infof("Root directory: %s", *rootDir)

	// Create root directory if it doesn't exist
	if err := os.MkdirAll(*rootDir, 0755); err != nil {
		logrus.Fatalf("Failed to create root directory: %v", err)
	}

	// Initialize Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		logrus.Fatalf("Failed to initialize Docker client: %v", err)
	}
	defer dockerClient.Close()

	// Create CRI runtime service
	runtimeService := cri.NewRuntimeService(dockerClient, *rootDir)

	// Create CRI image service
	imageService := cri.NewImageService(dockerClient)

	// Create and start gRPC server
	grpcServer := grpc.NewServer()
	server.RegisterServices(grpcServer, runtimeService, imageService)

	// Remove existing socket if it exists
	if err := os.Remove(*socketPath); err != nil && !os.IsNotExist(err) {
		logrus.Fatalf("Failed to remove existing socket: %v", err)
	}

	// Create socket directory if it doesn't exist
	socketDir := filepath.Dir(*socketPath)
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		logrus.Fatalf("Failed to create socket directory: %v", err)
	}

	// Listen on Unix socket
	listener, err := net.Listen("unix", *socketPath)
	if err != nil {
		logrus.Fatalf("Failed to listen on socket: %v", err)
	}
	defer listener.Close()

	// Set socket permissions
	if err := os.Chmod(*socketPath, 0660); err != nil {
		logrus.Warnf("Failed to set socket permissions: %v", err)
	}

	logrus.Infof("Docker CRI shim listening on %s", *socketPath)

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		logrus.Fatalf("Failed to serve: %v", err)
	}
}

