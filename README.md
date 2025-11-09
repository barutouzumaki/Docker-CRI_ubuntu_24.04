# Docker and Kubernetes Setup Guide for Ubuntu 24.04

This guide provides comprehensive instructions for setting up Docker, Docker CRI (Container Runtime Interface), and Kubernetes on Ubuntu 24.04.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Overview](#overview)
- [Manual Installation Steps](#manual-installation-steps)
- [Automated Installation](#automated-installation)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)
- [Additional Resources](#additional-resources)

## Prerequisites

Before starting, ensure you have:

- Ubuntu 24.04 LTS installed
- Root or sudo access
- At least 2 CPU cores
- Minimum 2GB RAM (4GB+ recommended)
- Internet connectivity
- Swap disabled (recommended for Kubernetes)

## Overview

This setup includes:

1. **Docker Engine**: Container runtime for running containers
2. **Docker CRI**: Container Runtime Interface shim that allows Kubernetes to use Docker
3. **Kubernetes Components**:
   - `kubeadm`: Tool for bootstrapping Kubernetes clusters
   - `kubelet`: Node agent that runs on each machine
   - `kubectl`: Command-line tool for interacting with Kubernetes clusters

## Manual Installation Steps

### Step 1: System Preparation

#### Update System Packages
```bash
sudo apt-get update
sudo apt-get upgrade -y
```

#### Install Prerequisites
```bash
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    wget \
    git
```

#### Disable Swap (Required for Kubernetes)
```bash
# Temporarily disable swap
sudo swapoff -a

# Permanently disable swap
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
```

#### Configure Kernel Parameters
```bash
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sudo sysctl --system
```

### Step 2: Install Docker Engine

#### Add Docker's Official GPG Key
```bash
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
```

#### Add Docker Repository
```bash
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

#### Install Docker Engine
```bash
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

#### Configure Docker for Kubernetes
```bash
# Configure containerd
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml

# Enable systemd cgroup driver
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml

# Restart containerd
sudo systemctl restart containerd
sudo systemctl enable containerd
```

#### Start and Enable Docker
```bash
sudo systemctl start docker
sudo systemctl enable docker
```

#### Verify Docker Installation
```bash
sudo docker run hello-world
```

### Step 3: Install Docker CRI

#### Install Go (if not already installed)
```bash
wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
rm go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

#### Build and Install Docker CRI
```bash
# Navigate to the docker-cri directory
cd Docker-CRI_ubuntu_24.04

# Download dependencies
go mod download
go mod tidy

# Build the binary
mkdir -p bin
go build -o bin/docker-cri ./main.go

# Install docker-cri
sudo cp bin/docker-cri /usr/local/bin/docker-cri
sudo chmod +x /usr/local/bin/docker-cri

# Create required directories
sudo mkdir -p /var/lib/docker-cri
sudo mkdir -p /var/run
```

#### Create Systemd Service for Docker CRI
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

### Step 4: Install Kubernetes Components

#### Add Kubernetes GPG Key
```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
```

#### Add Kubernetes Repository
```bash
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
```

#### Install Kubernetes Components
```bash
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

#### Configure kubelet to use Docker CRI
```bash
sudo mkdir -p /var/lib/kubelet
cat <<EOF | sudo tee /var/lib/kubelet/config.yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
EOF
```

### Step 5: Initialize Kubernetes Cluster (Master Node)

#### Initialize the Cluster
```bash
sudo kubeadm init \
  --pod-network-cidr=10.244.0.0/16 \
  --cri-socket=unix:///var/run/docker-cri.sock
```

**Important**: Save the `kubeadm join` command that is displayed after initialization. You'll need it to join worker nodes.

#### Configure kubectl
```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

#### Install Pod Network Addon (Flannel)
```bash
kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
```

#### Verify Cluster Status
```bash
kubectl get nodes
kubectl get pods --all-namespaces
```

### Step 6: Join Worker Nodes (Optional)

On each worker node, run the `kubeadm join` command that was displayed during master initialization:

```bash
sudo kubeadm join <MASTER_IP>:6443 --token <TOKEN> \
  --discovery-token-ca-cert-hash sha256:<HASH> \
  --cri-socket=unix:///var/run/docker-cri.sock
```

## Automated Installation

For automated installation, use the provided initialization script:

```bash
# Make the script executable
chmod +x setup-k8s-docker.sh

# Run the script (requires sudo)
sudo ./setup-k8s-docker.sh
```

The script will:
1. Update system packages
2. Install prerequisites
3. Configure kernel parameters
4. Install Docker Engine
5. Install Docker CRI
6. Install Kubernetes components
7. Optionally initialize a Kubernetes cluster

## Verification

### Verify Docker
```bash
sudo docker version
sudo docker info
sudo docker run hello-world
```

### Verify Docker CRI
```bash
# Install crictl if not already installed
wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.28.0/crictl-v1.28.0-linux-amd64.tar.gz
sudo tar zxvf crictl-v1.28.0-linux-amd64.tar.gz -C /usr/local/bin
rm crictl-v1.28.0-linux-amd64.tar.gz

# Test Docker CRI
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock version
```

### Verify Kubernetes
```bash
kubectl version --client
kubectl get nodes
kubectl cluster-info
```

### Check Service Status
```bash
sudo systemctl status docker
sudo systemctl status containerd
sudo systemctl status docker-cri
sudo systemctl status kubelet
```

## Troubleshooting

### Docker Issues

**Problem**: Docker daemon not running
```bash
sudo systemctl start docker
sudo systemctl status docker
```

**Problem**: Permission denied when running docker commands
```bash
sudo usermod -aG docker $USER
# Log out and log back in
```

### Docker CRI Issues

**Problem**: Docker CRI service not starting
```bash
sudo systemctl status docker-cri
sudo journalctl -u docker-cri -f
```

**Problem**: Socket file not found
```bash
# Check if socket exists
ls -la /var/run/docker-cri.sock

# Restart service
sudo systemctl restart docker-cri
```

### Kubernetes Issues

**Problem**: kubelet not running
```bash
sudo systemctl status kubelet
sudo journalctl -u kubelet -f
```

**Problem**: Pods stuck in Pending state
```bash
# Check node status
kubectl describe node

# Check if pod network is installed
kubectl get pods -n kube-flannel
```

**Problem**: Reset Kubernetes cluster
```bash
sudo kubeadm reset
sudo rm -rf /etc/cni/net.d
sudo rm -rf /var/lib/etcd
```

### Common Commands

```bash
# View all logs
sudo journalctl -xe

# Check system resources
free -h
df -h

# Check network connectivity
ping 8.8.8.8
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kubeadm Documentation](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)
- [Container Runtime Interface (CRI)](https://kubernetes.io/docs/concepts/architecture/cri/)

## Security Considerations

1. **Firewall**: Configure firewall rules to allow Kubernetes ports (6443, 10250, etc.)
2. **RBAC**: Configure Role-Based Access Control for Kubernetes
3. **Network Policies**: Implement network policies for pod-to-pod communication
4. **Secrets Management**: Use Kubernetes secrets or external secret management tools
5. **Regular Updates**: Keep all components updated with security patches

## Maintenance

### Update Components

```bash
# Update Docker
sudo apt-get update
sudo apt-get upgrade docker-ce docker-ce-cli containerd.io

# Update Kubernetes
sudo apt-get update
sudo apt-get upgrade kubelet kubeadm kubectl
```

### Backup Configuration

```bash
# Backup Kubernetes configuration
sudo cp -r /etc/kubernetes ~/k8s-backup
sudo cp -r ~/.kube ~/kube-backup
```

## Support

For issues or questions:
1. Check the troubleshooting section
2. Review service logs using `journalctl`
3. Consult official documentation
4. Check GitHub issues for known problems
