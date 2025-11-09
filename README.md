# Docker CRI (Container Runtime Interface)

A Kubernetes CRI (Container Runtime Interface) implementation for Docker, compatible with Ubuntu 24.04.

## Overview

This project provides a gRPC-based CRI shim that translates Kubernetes CRI API calls to Docker API calls, allowing Kubernetes to use Docker as its container runtime.

## Features

- **RuntimeService**: Full implementation of container lifecycle management
  - Pod sandbox creation and management
  - Container creation, start, stop, and removal
  - Container status and statistics
  - Container listing and filtering

- **ImageService**: Complete image management
  - Image pulling with authentication
  - Image listing and status
  - Image removal
  - Filesystem information

## Prerequisites

### Ubuntu 24.04

1. **Docker Engine** (20.10 or later)
   ```bash
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   
   # Add your user to docker group (optional)
   sudo usermod -aG docker $USER
   ```

2. **Go** (1.21 or later)
   ```bash
   # Install Go
   wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

3. **Kubernetes** (optional, for testing)
   ```bash
   # Install kubectl
   curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
   sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
   ```

## Installation

### Method 1: Build from Source

```bash
# Clone or navigate to the project directory
cd docker-cri

# Download dependencies
make deps

# Build the binary
make build

# Install system-wide
sudo make install
```

### Method 2: Using Docker

```bash
# Build the Docker image
docker build -t docker-cri:latest .

# Run the container
docker run -d \
  --name docker-cri \
  --privileged \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/run:/var/run \
  -v /var/lib/docker-cri:/var/lib/docker-cri \
  docker-cri:latest
```

## Configuration

### Command Line Options

```bash
docker-cri [options]

Options:
  -socket string
        Path to the CRI socket (default "/var/run/docker-cri.sock")
  -log-level string
        Log level: debug, info, warn, error (default "info")
  -root-dir string
        Root directory for storing runtime data (default "/var/lib/docker-cri")
```

### Example Usage

```bash
# Run with default settings
docker-cri

# Run with custom socket path and log level
docker-cri -socket /var/run/custom-cri.sock -log-level debug

# Run as a systemd service
sudo systemctl start docker-cri
```

## Systemd Service

Create a systemd service file for automatic startup:

```bash
sudo tee /etc/systemd/system/docker-cri.service > /dev/null <<EOF
[Unit]
Description=Docker CRI Shim
After=docker.service
Requires=docker.service

[Service]
Type=simple
ExecStart=/usr/local/bin/docker-cri -socket /var/run/docker-cri.sock -log-level info
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable docker-cri
sudo systemctl start docker-cri
```

## Kubernetes Configuration

To use this CRI with Kubernetes, configure kubelet:

```bash
# Edit kubelet configuration
sudo nano /var/lib/kubelet/config.yaml
```

Add or modify:
```yaml
containerRuntimeEndpoint: "unix:///var/run/docker-cri.sock"
```

Or use command-line flag:
```bash
--container-runtime-endpoint=unix:///var/run/docker-cri.sock
```

## Testing

### Test the CRI Service

```bash
# Check if the socket is created
ls -la /var/run/docker-cri.sock

# Test with crictl (if installed)
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock version
```

### Build and Run Tests

```bash
# Run unit tests
make test

# Run with verbose output
go test -v ./...
```

## Architecture

```
┌─────────────────┐
│   Kubernetes    │
│   (kubelet)     │
└────────┬────────┘
         │ gRPC (CRI API)
         │
┌────────▼────────┐
│   docker-cri    │
│   (CRI Shim)    │
└────────┬────────┘
         │ Docker API
         │
┌────────▼────────┐
│  Docker Engine  │
└─────────────────┘
```

## API Coverage

### RuntimeService
- ✅ Version
- ✅ RunPodSandbox
- ✅ StopPodSandbox
- ✅ RemovePodSandbox
- ✅ PodSandboxStatus
- ✅ ListPodSandbox
- ✅ CreateContainer
- ✅ StartContainer
- ✅ StopContainer
- ✅ RemoveContainer
- ✅ ContainerStatus
- ✅ ListContainers
- ✅ ContainerStats
- ✅ ListContainerStats
- ✅ UpdateContainerResources
- ✅ ReopenContainerLog
- ✅ Status
- ⚠️ ExecSync (not implemented)
- ⚠️ Exec (not implemented)
- ⚠️ Attach (not implemented)
- ⚠️ PortForward (not implemented)

### ImageService
- ✅ ListImages
- ✅ ImageStatus
- ✅ PullImage
- ✅ RemoveImage
- ✅ ImageFsInfo

## Troubleshooting

### Socket Permission Issues

```bash
# Fix socket permissions
sudo chmod 666 /var/run/docker-cri.sock
sudo chown root:docker /var/run/docker-cri.sock
```

### Docker Connection Issues

```bash
# Verify Docker is running
sudo systemctl status docker

# Test Docker connection
docker ps
```

### Logs

```bash
# View systemd logs
sudo journalctl -u docker-cri -f

# View with custom log level
docker-cri -log-level debug
```

## Development

### Project Structure

```
docker-cri/
├── main.go              # Entry point
├── go.mod               # Go module definition
├── Makefile             # Build automation
├── Dockerfile           # Container image
├── README.md            # This file
└── pkg/
    ├── cri/
    │   ├── runtime.go   # RuntimeService implementation
    │   └── image.go     # ImageService implementation
    ├── docker/
    │   └── client.go    # Docker client wrapper
    ├── server/
    │   └── server.go    # gRPC server setup
    └── utils/
        └── utils.go     # Utility functions
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Build: `make build`
6. Submit a pull request

## License

This project is provided as-is for experimental use.

## References

- [Kubernetes CRI Specification](https://kubernetes.io/docs/concepts/architecture/cri/)
- [Docker API Documentation](https://docs.docker.com/engine/api/)
- [CRI API Protobuf Definitions](https://github.com/kubernetes/cri-api)

## Support

For issues and questions, please open an issue in the repository.

