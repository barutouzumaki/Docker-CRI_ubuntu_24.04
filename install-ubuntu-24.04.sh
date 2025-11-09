#!/bin/bash

# Installation script for Docker CRI on Ubuntu 24.04

set -e

echo "=========================================="
echo "Docker CRI Installation for Ubuntu 24.04"
echo "=========================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (use sudo)"
    exit 1
fi

# Check Ubuntu version
if [ ! -f /etc/os-release ]; then
    echo "Cannot determine OS version"
    exit 1
fi

. /etc/os-release
if [ "$ID" != "ubuntu" ] || [ "$VERSION_ID" != "24.04" ]; then
    echo "Warning: This script is designed for Ubuntu 24.04"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Install prerequisites
echo "Installing prerequisites..."
apt-get update
apt-get install -y curl wget ca-certificates

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
    echo "Docker installed successfully"
else
    echo "Docker is already installed"
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go 1.21..."
    wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
    rm go1.21.5.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "Go installed successfully"
else
    echo "Go is already installed: $(go version)"
fi

# Build docker-cri
echo "Building docker-cri..."
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Please run this script from the docker-cri directory"
    exit 1
fi

# Download dependencies
echo "Downloading dependencies..."
export PATH=$PATH:/usr/local/go/bin
go mod download
go mod tidy

# Build
echo "Building binary..."
mkdir -p bin
go build -o bin/docker-cri ./main.go

# Install
echo "Installing docker-cri..."
cp bin/docker-cri /usr/local/bin/docker-cri
chmod +x /usr/local/bin/docker-cri

# Create directories
echo "Creating directories..."
mkdir -p /var/lib/docker-cri
mkdir -p /var/run

# Create systemd service
echo "Creating systemd service..."
cat > /etc/systemd/system/docker-cri.service <<EOF
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

# Reload systemd and enable service
echo "Enabling docker-cri service..."
systemctl daemon-reload
systemctl enable docker-cri

echo ""
echo "=========================================="
echo "Installation complete!"
echo "=========================================="
echo ""
echo "To start the service:"
echo "  sudo systemctl start docker-cri"
echo ""
echo "To check status:"
echo "  sudo systemctl status docker-cri"
echo ""
echo "To view logs:"
echo "  sudo journalctl -u docker-cri -f"
echo ""
echo "To test the CRI:"
echo "  sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock version"
echo ""

