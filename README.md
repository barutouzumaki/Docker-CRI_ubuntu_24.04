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

### Why Docker CRI Setup Matters for Ubuntu 24.04

Ubuntu 24.04 LTS (Noble Numbat) introduces several improvements that make it an excellent platform for containerized workloads:

- **Enhanced Kernel Support**: Linux kernel 6.8+ with improved cgroup v2 support
- **Better Systemd Integration**: Native systemd cgroup driver support for better resource management
- **Improved Security**: Enhanced AppArmor and seccomp profiles
- **Performance Optimizations**: Better I/O scheduling and memory management

### Docker CRI vs containerd vs CRI-O

When setting up Kubernetes on Ubuntu 24.04, you have three main container runtime options:

#### Docker CRI (Docker Engine + CRI Shim)

**Advantages:**
- ✅ Familiar Docker tooling and ecosystem
- ✅ Full Docker CLI access (`docker ps`, `docker logs`, etc.)
- ✅ Easy debugging with standard Docker commands
- ✅ Rich ecosystem of Docker images and tools
- ✅ Good for teams already familiar with Docker
- ✅ Supports Docker Compose for local development

**Disadvantages:**
- ❌ Additional layer (CRI shim) adds slight overhead
- ❌ Not the most lightweight option
- ❌ Docker Engine includes features not needed by Kubernetes
- ❌ More complex architecture (dockerd → containerd → runc)

**Best For:**
- Teams transitioning from Docker to Kubernetes
- Development environments
- When you need Docker CLI access alongside Kubernetes

#### containerd

**Advantages:**
- ✅ Industry-standard, production-ready runtime
- ✅ Lightweight and efficient (no Docker daemon overhead)
- ✅ Direct CRI implementation (no shim needed)
- ✅ Better performance and resource usage
- ✅ Used by Docker Engine internally
- ✅ Excellent for production Kubernetes clusters
- ✅ Active development and wide adoption

**Disadvantages:**
- ❌ No Docker CLI (uses `crictl` instead)
- ❌ Less familiar to Docker users
- ❌ Requires learning new tools

**Best For:**
- Production Kubernetes clusters
- When you want optimal performance
- Cloud-native deployments
- Most recommended by Kubernetes community

#### CRI-O

**Advantages:**
- ✅ Lightweight and minimal (built specifically for Kubernetes)
- ✅ OCI-compliant runtime
- ✅ Good security defaults
- ✅ Active development by Red Hat
- ✅ No Docker dependencies

**Disadvantages:**
- ❌ Less mature ecosystem compared to containerd
- ❌ Primarily focused on Kubernetes (less flexible)
- ❌ Smaller community than containerd
- ❌ Requires learning CRI-O specific tools

**Best For:**
- Red Hat/CentOS/RHEL environments
- Minimal Kubernetes-only deployments
- Security-focused environments

#### Comparison Table

| Feature | Docker CRI | containerd | CRI-O |
|---------|-----------|------------|-------|
| **Performance** | Good | Excellent | Excellent |
| **Resource Usage** | Higher | Lower | Lower |
| **Docker CLI Support** | ✅ Yes | ❌ No | ❌ No |
| **Production Ready** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Complexity** | Higher | Lower | Lower |
| **Learning Curve** | Lower (if familiar with Docker) | Medium | Medium |
| **Community Support** | Large | Very Large | Growing |
| **Recommended For** | Development, Migration | Production | Production (RHEL) |

#### Recommendation for Ubuntu 24.04

For **production environments**, **containerd** is the recommended choice due to:
- Better performance and resource efficiency
- Direct CRI implementation
- Industry standard and widely adopted
- Excellent Ubuntu 24.04 support

For **development or migration scenarios**, **Docker CRI** provides:
- Familiar Docker tooling
- Easy debugging
- Smooth transition path

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

#### Problem: Docker daemon not running

**Symptoms:**
- `docker: Cannot connect to the Docker daemon`
- `Error response from daemon: dial unix /var/run/docker.sock: connect: permission denied`

**Diagnosis:**
```bash
# Check Docker service status
sudo systemctl status docker

# Check Docker logs
sudo journalctl -u docker -n 50 --no-pager

# Check if Docker socket exists
ls -la /var/run/docker.sock
```

**Solution:**
```bash
# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Verify Docker is running
sudo docker ps
```

#### Problem: Permission denied when running docker commands

**Symptoms:**
- `permission denied while trying to connect to the Docker daemon socket`

**Solution:**
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Apply changes (logout/login or use newgrp)
newgrp docker

# Verify
docker ps
```

#### Problem: Docker containers not starting

**Diagnosis:**
```bash
# Check Docker daemon logs
sudo journalctl -u docker -f

# Check container logs
sudo docker logs <container-id>

# Inspect container
sudo docker inspect <container-id>

# Check Docker info
sudo docker info
```

**Common Causes:**
- Insufficient disk space: `df -h`
- Memory issues: `free -h`
- Image pull failures: `sudo docker pull <image>`
- Port conflicts: `sudo netstat -tulpn | grep <port>`

#### Problem: Docker build failures

**Diagnosis:**
```bash
# Build with verbose output
sudo docker build --progress=plain --no-cache -t <image> .

# Check build context
sudo docker system df
```

**Common Issues:**
- Dockerfile syntax errors
- Network issues during build
- Insufficient disk space
- Base image pull failures

### Docker CRI Issues

#### Problem: Docker CRI service not starting

**Symptoms:**
- `Failed to start docker-cri.service`
- Socket file not created

**Diagnosis:**
```bash
# Check service status
sudo systemctl status docker-cri

# View detailed logs
sudo journalctl -u docker-cri -n 100 --no-pager

# Check if Docker is running (dependency)
sudo systemctl status docker

# Check socket file
ls -la /var/run/docker-cri.sock
```

**Solution:**
```bash
# Check Docker CRI binary exists
which docker-cri
/usr/local/bin/docker-cri --help

# Restart Docker first
sudo systemctl restart docker

# Restart Docker CRI
sudo systemctl restart docker-cri

# Check logs in real-time
sudo journalctl -u docker-cri -f
```

#### Problem: Socket file not found

**Symptoms:**
- `crictl: connect: no such file or directory: /var/run/docker-cri.sock`

**Diagnosis:**
```bash
# Check if socket exists
ls -la /var/run/docker-cri.sock

# Check service is running
sudo systemctl is-active docker-cri

# Check permissions
sudo ls -la /var/run/ | grep docker-cri
```

**Solution:**
```bash
# Restart service
sudo systemctl restart docker-cri

# Wait a few seconds
sleep 5

# Verify socket creation
ls -la /var/run/docker-cri.sock

# Test CRI connection
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock version
```

#### Problem: CRI connection timeout

**Symptoms:**
- `crictl: context deadline exceeded`

**Diagnosis:**
```bash
# Check service health
sudo systemctl status docker-cri

# Test socket connectivity
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock info

# Check for port conflicts
sudo netstat -tulpn | grep docker
```

### Kubernetes Issues

#### Problem: kubelet not running

**Symptoms:**
- Node shows as `NotReady`
- Pods not starting

**Diagnosis:**
```bash
# Check kubelet status
sudo systemctl status kubelet

# View kubelet logs
sudo journalctl -u kubelet -n 100 --no-pager

# Check kubelet configuration
cat /var/lib/kubelet/config.yaml

# Check for errors
sudo journalctl -u kubelet | grep -i error
```

**Common Errors and Solutions:**

1. **Cgroup driver mismatch:**
   ```bash
   # Check current cgroup driver
   docker info | grep -i cgroup
   
   # Update kubelet config
   sudo sed -i 's/cgroupDriver:.*/cgroupDriver: systemd/' /var/lib/kubelet/config.yaml
   sudo systemctl restart kubelet
   ```

2. **CRI socket not found:**
   ```bash
   # Verify CRI socket
   ls -la /var/run/docker-cri.sock
   
   # Update kubelet args (if needed)
   sudo vim /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
   # Add: --container-runtime-endpoint=unix:///var/run/docker-cri.sock
   sudo systemctl daemon-reload
   sudo systemctl restart kubelet
   ```

3. **Swap not disabled:**
   ```bash
   # Check swap
   free -h
   swapon --show
   
   # Disable swap
   sudo swapoff -a
   sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
   sudo systemctl restart kubelet
   ```

#### Problem: Pods stuck in Pending state

**Symptoms:**
- `kubectl get pods` shows `Pending`
- Pods not scheduling

**Diagnosis:**
```bash
# Describe pod for events
kubectl describe pod <pod-name> -n <namespace>

# Check node status
kubectl get nodes
kubectl describe node <node-name>

# Check node resources
kubectl top nodes

# Check pod network
kubectl get pods -n kube-system
```

**Common Causes:**

1. **Insufficient resources:**
   ```bash
   # Check node capacity
   kubectl describe node | grep -A 5 "Allocated resources"
   
   # Check pod resource requests
   kubectl describe pod <pod-name> | grep -A 5 "Requests:"
   ```

2. **Node not ready:**
   ```bash
   # Check node conditions
   kubectl get nodes -o wide
   
   # Check kubelet on node
   ssh <node> sudo systemctl status kubelet
   ```

3. **Pod network not installed:**
   ```bash
   # Check CNI pods
   kubectl get pods -n kube-flannel
   
   # Install Flannel if missing
   kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
   ```

4. **Taints preventing scheduling:**
   ```bash
   # Check node taints
   kubectl describe node | grep Taints
   
   # Remove taint if needed (for master node)
   kubectl taint nodes --all node-role.kubernetes.io/control-plane-
   ```

#### Problem: Pods stuck in ContainerCreating

**Symptoms:**
- Pods show `ContainerCreating` for extended time
- Containers not starting

**Diagnosis:**
```bash
# Describe pod
kubectl describe pod <pod-name> -n <namespace>

# Check events
kubectl get events --sort-by='.lastTimestamp' -n <namespace>

# Check container runtime
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock ps -a

# Check image pull issues
kubectl describe pod <pod-name> | grep -i image
```

**Common Causes:**

1. **Image pull errors:**
   ```bash
   # Check image pull secrets
   kubectl get secrets
   
   # Manually pull image
   sudo docker pull <image-name>
   ```

2. **Volume mount issues:**
   ```bash
   # Check persistent volumes
   kubectl get pv
   kubectl get pvc
   
   # Check volume mounts in pod
   kubectl describe pod <pod-name> | grep -A 10 "Mounts:"
   ```

3. **Resource constraints:**
   ```bash
   # Check node resources
   kubectl top nodes
   
   # Check pod resource limits
   kubectl describe pod <pod-name> | grep -A 5 "Limits:"
   ```

#### Problem: Pods in CrashLoopBackOff

**Symptoms:**
- Pods restarting repeatedly
- `CrashLoopBackOff` status

**Diagnosis:**
```bash
# Check pod logs
kubectl logs <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace> --previous

# Describe pod for events
kubectl describe pod <pod-name> -n <namespace>

# Check container exit codes
kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.status.containerStatuses[0].lastState.terminated.exitCode}'
```

**Common Causes:**
- Application errors (check logs)
- Configuration errors
- Resource limits exceeded
- Health check failures

#### Problem: Service not accessible

**Symptoms:**
- Cannot connect to service
- Service endpoints empty

**Diagnosis:**
```bash
# Check service
kubectl get svc <service-name> -n <namespace>
kubectl describe svc <service-name> -n <namespace>

# Check endpoints
kubectl get endpoints <service-name> -n <namespace>

# Check pod labels match service selector
kubectl get pods -n <namespace> --show-labels
kubectl get svc <service-name> -n <namespace> -o yaml | grep selector
```

**Solution:**
```bash
# Verify pod labels match service selector
# Update labels or service selector as needed

# Check kube-proxy
kubectl get pods -n kube-system | grep kube-proxy
kubectl logs -n kube-system <kube-proxy-pod>
```

#### Problem: Network connectivity issues

**Diagnosis:**
```bash
# Check CNI pods
kubectl get pods -n kube-flannel

# Check CNI logs
kubectl logs -n kube-flannel <cni-pod>

# Check network policies
kubectl get networkpolicies --all-namespaces

# Test pod-to-pod connectivity
kubectl run test-pod --image=busybox --rm -it -- ping <target-pod-ip>
```

#### Problem: etcd issues

**Symptoms:**
- API server errors
- Cluster state inconsistent

**Diagnosis:**
```bash
# Check etcd pods (if using etcd in cluster)
kubectl get pods -n kube-system | grep etcd

# Check etcd health (on control plane)
sudo ETCDCTL_API=3 etcdctl --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key \
  endpoint health
```

### Log Locations and Commands

#### System Logs
```bash
# All system logs
sudo journalctl -xe

# Docker logs
sudo journalctl -u docker -f
sudo journalctl -u docker --since "1 hour ago"

# Docker CRI logs
sudo journalctl -u docker-cri -f
sudo journalctl -u docker-cri --since "1 hour ago"

# Kubelet logs
sudo journalctl -u kubelet -f
sudo journalctl -u kubelet --since "1 hour ago"

# All Kubernetes component logs
sudo journalctl -u kubelet -u docker -u docker-cri --since "1 hour ago"
```

#### Kubernetes Logs
```bash
# Pod logs
kubectl logs <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace> --previous
kubectl logs <pod-name> -n <namespace> --tail=100 -f

# All pods in namespace
kubectl logs -l app=<label> -n <namespace>

# Container logs (multi-container pod)
kubectl logs <pod-name> -c <container-name> -n <namespace>

# System component logs
kubectl logs -n kube-system <component-pod>
```

#### Container Runtime Logs
```bash
# Docker logs
sudo docker logs <container-id>
sudo docker logs --tail 100 -f <container-id>

# CRI logs
sudo crictl logs <container-id>
sudo crictl logs --tail 100 <container-id>
```

### Common Error Messages and Solutions

#### Error: "failed to get sandbox runtime: no runtime for \"docker\" is configured"
**Solution:**
```bash
# Verify CRI socket configuration
sudo systemctl status docker-cri
ls -la /var/run/docker-cri.sock

# Update kubelet configuration
sudo vim /var/lib/kubelet/config.yaml
# Ensure containerRuntimeEndpoint is set correctly
```

#### Error: "cgroupfs is not supported with systemd cgroup driver"
**Solution:**
```bash
# Update containerd config
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
sudo systemctl restart containerd

# Update kubelet config
echo "cgroupDriver: systemd" | sudo tee -a /var/lib/kubelet/config.yaml
sudo systemctl restart kubelet
```

#### Error: "The connection to the server localhost:8080 was refused"
**Solution:**
```bash
# Set kubeconfig
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Or export KUBECONFIG
export KUBECONFIG=/etc/kubernetes/admin.conf
```

#### Error: "failed to pull image: rpc error"
**Solution:**
```bash
# Check image registry access
sudo docker pull <image-name>

# Check registry credentials
kubectl get secrets -n <namespace>

# Check network connectivity
ping registry-url
```

### Reset Kubernetes Cluster

**Complete Reset:**
```bash
# Drain and delete node (if in cluster)
kubectl drain <node-name> --ignore-daemonsets --delete-emptydir-data
kubectl delete node <node-name>

# Reset kubeadm
sudo kubeadm reset -f

# Clean up
sudo rm -rf /etc/cni/net.d
sudo rm -rf /var/lib/etcd
sudo rm -rf /var/lib/kubelet
sudo rm -rf /etc/kubernetes
sudo rm -rf ~/.kube

# Restart services
sudo systemctl restart docker
sudo systemctl restart docker-cri
sudo systemctl restart kubelet
```

### Diagnostic Commands

```bash
# System resources
free -h                    # Memory
df -h                      # Disk space
top                        # CPU and processes
iostat -x 1 5              # I/O statistics

# Network
ip addr show               # Network interfaces
ip route show              # Routing table
netstat -tulpn             # Listening ports
ss -tulpn                  # Socket statistics

# Docker
sudo docker ps             # Running containers
sudo docker images         # Images
sudo docker system df      # Disk usage
sudo docker info           # Docker information

# Kubernetes
kubectl cluster-info       # Cluster information
kubectl get nodes -o wide # Node status
kubectl get pods --all-namespaces # All pods
kubectl get events --all-namespaces --sort-by='.lastTimestamp' # Recent events

# CRI
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock ps
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock images
sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock info
```

## Performance Tuning

### Kubelet Performance Tuning

#### Resource Limits and Requests

```bash
# Configure kubelet resource limits
sudo vim /var/lib/kubelet/config.yaml
```

Add or modify:
```yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
# Maximum number of pods per node
maxPods: 110
# System reserved resources
systemReserved:
  cpu: "500m"
  memory: "1Gi"
  ephemeral-storage: "1Gi"
# Kube reserved resources
kubeReserved:
  cpu: "500m"
  memory: "1Gi"
  ephemeral-storage: "1Gi"
# Eviction thresholds
evictionHard:
  memory.available: "200Mi"
  nodefs.available: "10%"
  imagefs.available: "15%"
# Pod PIDs limit
podPidsLimit: 4096
```

#### Kubelet Garbage Collection

```yaml
# Configure image garbage collection
imageGCHighThresholdPercent: 85
imageGCLowThresholdPercent: 80
# Configure container garbage collection
containerGCMaxPerPodContainer: 5
containerGCMaxContainers: 50
```

#### Kubelet API Server Connection

```yaml
# Optimize API server connection
kubeAPIQPS: 50
kubeAPIBurst: 100
# Connection timeout
httpCheckTimeout: 20s
```

**Apply changes:**
```bash
sudo systemctl daemon-reload
sudo systemctl restart kubelet
```

### CRI Performance Tuning

#### Docker Daemon Tuning

```bash
# Edit Docker daemon configuration
sudo vim /etc/docker/daemon.json
```

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ],
  "default-ulimits": {
    "nofile": {
      "Name": "nofile",
      "Hard": 64000,
      "Soft": 64000
    }
  },
  "max-concurrent-downloads": 10,
  "max-concurrent-uploads": 5,
  "default-address-pools": [
    {
      "base": "172.17.0.0/16",
      "size": 24
    }
  ]
}
```

**Apply changes:**
```bash
sudo systemctl restart docker
```

#### containerd Tuning (if using containerd directly)

```bash
sudo vim /etc/containerd/config.toml
```

```toml
[plugins."io.containerd.grpc.v1.cri"]
  # Max concurrent downloads
  max_concurrent_downloads = 10
  # Snapshotter
  snapshotter = "overlayfs"
  # Stream server settings
  stream_server_address = "127.0.0.1"
  stream_server_port = "0"
```

### Networking Performance Tuning

#### kube-proxy Tuning

For better performance, use IPVS mode instead of iptables:

```bash
# Edit kube-proxy ConfigMap
kubectl edit configmap kube-proxy -n kube-system
```

Change mode to IPVS:
```yaml
mode: "ipvs"
ipvs:
  scheduler: "rr"  # round-robin, or "lc" for least-connection
  minSyncPeriod: 5s
  syncPeriod: 30s
```

**Restart kube-proxy:**
```bash
kubectl delete pods -n kube-system -l k8s-app=kube-proxy
```

#### Network Plugin Tuning (Flannel)

```bash
# Edit Flannel ConfigMap
kubectl edit configmap kube-flannel-cfg -n kube-flannel
```

```yaml
net-conf.json: |
  {
    "Network": "10.244.0.0/16",
    "Backend": {
      "Type": "vxlan",
      "VNI": 1,
      "Port": 8472
    }
  }
```

#### System Network Tuning

```bash
# Add to /etc/sysctl.conf
sudo tee -a /etc/sysctl.conf <<EOF
# Increase connection tracking table size
net.netfilter.nf_conntrack_max = 1048576

# TCP tuning
net.core.somaxconn = 32768
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_tw_reuse = 1
net.ipv4.ip_local_port_range = 1024 65535

# Increase buffer sizes
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216

# Disable TCP slow start after idle
net.ipv4.tcp_slow_start_after_idle = 0
EOF

sudo sysctl -p
```

### Storage Performance Tuning

#### Docker Storage Driver Optimization

For Ubuntu 24.04, overlay2 is recommended:

```bash
# Verify storage driver
docker info | grep "Storage Driver"

# If using overlay2, optimize
sudo vim /etc/docker/daemon.json
```

```json
{
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true",
    "overlay2.size=20G"
  ]
}
```

#### Disk I/O Tuning

```bash
# Set I/O scheduler (for SSDs)
echo deadline | sudo tee /sys/block/sda/queue/scheduler

# Increase read-ahead buffer
sudo blockdev --setra 8192 /dev/sda

# Make permanent
sudo tee -a /etc/rc.local <<EOF
echo deadline > /sys/block/sda/queue/scheduler
blockdev --setra 8192 /dev/sda
EOF
```

#### Volume Performance

For persistent volumes, consider:
- Using local SSDs for high I/O workloads
- Configuring appropriate storage classes
- Using ReadWriteMany volumes for shared access
- Implementing volume snapshots for backups

### System-Level Performance Tuning

#### CPU Governor

```bash
# Set performance governor
sudo cpupower frequency-set -g performance

# Make permanent
sudo tee /etc/default/cpupower <<EOF
GOVERNOR=performance
EOF
```

#### Memory Management

```bash
# Disable swap (already done, but verify)
sudo swapoff -a

# Tune vm.swappiness
echo "vm.swappiness=1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Tune overcommit memory
echo "vm.overcommit_memory=1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

#### File Descriptor Limits

```bash
# Increase file descriptor limits
sudo tee -a /etc/security/limits.conf <<EOF
* soft nofile 65536
* hard nofile 65536
root soft nofile 65536
root hard nofile 65536
EOF
```

### Monitoring Performance

```bash
# Install metrics-server (if not installed)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Monitor node resources
kubectl top nodes

# Monitor pod resources
kubectl top pods --all-namespaces

# Check kubelet metrics
curl http://localhost:10255/metrics

# Check Docker stats
sudo docker stats --no-stream
```

## Security Guidance

### AppArmor Profiles

AppArmor provides mandatory access control (MAC) for applications.

#### Install AppArmor

```bash
# Check if AppArmor is installed
sudo aa-status

# Install AppArmor utilities
sudo apt-get install -y apparmor-utils apparmor-profiles
```

#### Create AppArmor Profile for Containers

```bash
# Generate profile for a container
sudo aa-genprof docker-default

# Or create custom profile
sudo vim /etc/apparmor.d/docker-default
```

Example profile:
```
#include <tunables/global>

profile docker-default flags=(attach_disconnected,mediate_deleted) {
  #include <abstractions/base>
  
  network,
  capability,
  file,
  umount,
  
  deny @{PROC}/* w,
  deny /sys/[^f]** w,
  deny /sys/f[^s]** w,
  deny /sys/fs/[^c]** w,
  deny /sys/fs/c[^g]** w,
  deny /sys/fs/cg[^r]** w,
  deny /sys/firmware/** rwkl,
  deny /sys/kernel/security/** rwkl,
}
```

#### Apply AppArmor Profile to Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: apparmor-pod
  annotations:
    container.apparmor.security.beta.kubernetes.io/<container-name>: localhost/docker-default
spec:
  containers:
  - name: <container-name>
    image: nginx
```

#### Load and Enable Profile

```bash
# Load profile
sudo apparmor_parser -r /etc/apparmor.d/docker-default

# Check status
sudo aa-status | grep docker-default
```

### Seccomp Profiles

Seccomp (Secure Computing Mode) restricts system calls available to containers.

#### Default Seccomp Profile

Kubernetes applies a default seccomp profile. Check if enabled:

```bash
# Verify seccomp support
grep -i seccomp /proc/self/status
```

#### Create Custom Seccomp Profile

```bash
# Create seccomp profile directory
sudo mkdir -p /var/lib/kubelet/seccomp/profiles

# Create custom profile
sudo vim /var/lib/kubelet/seccomp/profiles/deny-all.json
```

```json
{
  "defaultAction": "SCMP_ACT_ERRNO",
  "architectures": [
    "SCMP_ARCH_X86_64",
    "SCMP_ARCH_X86",
    "SCMP_ARCH_X32"
  ],
  "syscalls": [
    {
      "names": [
        "accept",
        "accept4",
        "access",
        "adjtimex",
        "alarm",
        "bind",
        "brk",
        "capget",
        "capset",
        "chdir",
        "chmod",
        "chown",
        "chroot",
        "clock_getres",
        "clock_gettime",
        "clock_nanosleep",
        "close",
        "connect",
        "copy_file_range",
        "creat",
        "dup",
        "dup2",
        "dup3",
        "epoll_create",
        "epoll_create1",
        "epoll_ctl",
        "epoll_pwait",
        "epoll_wait",
        "eventfd",
        "eventfd2",
        "execve",
        "execveat",
        "exit",
        "exit_group",
        "faccessat",
        "fadvise64",
        "fallocate",
        "fanotify_mark",
        "fchdir",
        "fchmod",
        "fchmodat",
        "fchown",
        "fchownat",
        "fcntl",
        "fdatasync",
        "fgetxattr",
        "finit_module",
        "flistxattr",
        "flock",
        "fork",
        "fremovexattr",
        "fsetxattr",
        "fstat",
        "fstatfs",
        "fsync",
        "ftruncate",
        "futex",
        "getcwd",
        "getdents",
        "getdents64",
        "getegid",
        "geteuid",
        "getgid",
        "getgroups",
        "getpeername",
        "getpgid",
        "getpgrp",
        "getpid",
        "getppid",
        "getpriority",
        "getrandom",
        "getresgid",
        "getresuid",
        "getrlimit",
        "getrusage",
        "getsid",
        "getsockname",
        "getsockopt",
        "get_thread_area",
        "gettid",
        "gettimeofday",
        "getuid",
        "getxattr",
        "inotify_add_watch",
        "inotify_init",
        "inotify_init1",
        "inotify_rm_watch",
        "io_cancel",
        "io_destroy",
        "io_getevents",
        "io_setup",
        "io_submit",
        "ioctl",
        "ioprio_get",
        "ioprio_set",
        "ipc",
        "keyctl",
        "kill",
        "lgetxattr",
        "link",
        "linkat",
        "listen",
        "listxattr",
        "llistxattr",
        "_llseek",
        "lremovexattr",
        "lseek",
        "lsetxattr",
        "lstat",
        "madvise",
        "memfd_create",
        "mincore",
        "mkdir",
        "mkdirat",
        "mknod",
        "mknodat",
        "mlock",
        "mlock2",
        "mlockall",
        "mmap",
        "mmap2",
        "mprotect",
        "mq_getsetattr",
        "mq_notify",
        "mq_open",
        "mq_timedreceive",
        "mq_timedsend",
        "mq_unlink",
        "mremap",
        "msgctl",
        "msgget",
        "msgrcv",
        "msgsnd",
        "msync",
        "munlock",
        "munlockall",
        "munmap",
        "nanosleep",
        "newfstatat",
        "_newselect",
        "open",
        "openat",
        "pause",
        "personality",
        "pipe",
        "pipe2",
        "poll",
        "ppoll",
        "prctl",
        "pread64",
        "preadv",
        "prlimit64",
        "pselect6",
        "ptrace",
        "pwrite64",
        "pwritev",
        "read",
        "readahead",
        "readlink",
        "readlinkat",
        "readv",
        "recv",
        "recvfrom",
        "recvmmsg",
        "recvmsg",
        "remap_file_pages",
        "removexattr",
        "rename",
        "renameat",
        "renameat2",
        "restart_syscall",
        "rmdir",
        "rt_sigaction",
        "rt_sigpending",
        "rt_sigprocmask",
        "rt_sigqueueinfo",
        "rt_sigreturn",
        "rt_sigsuspend",
        "rt_sigtimedwait",
        "rt_tgsigqueueinfo",
        "sched_getaffinity",
        "sched_getattr",
        "sched_getparam",
        "sched_get_priority_max",
        "sched_get_priority_min",
        "sched_getscheduler",
        "sched_setaffinity",
        "sched_setattr",
        "sched_setparam",
        "sched_setscheduler",
        "sched_yield",
        "seccomp",
        "select",
        "semctl",
        "semget",
        "semop",
        "semtimedop",
        "send",
        "sendfile",
        "sendfile64",
        "sendmmsg",
        "sendmsg",
        "sendto",
        "setfsgid",
        "setfsuid",
        "setgid",
        "setgroups",
        "setitimer",
        "setpgid",
        "setpriority",
        "setregid",
        "setresgid",
        "setresuid",
        "setreuid",
        "setrlimit",
        "setsid",
        "setsockopt",
        "set_thread_area",
        "set_tid_address",
        "setuid",
        "setxattr",
        "shmat",
        "shmctl",
        "shmdt",
        "shmget",
        "shutdown",
        "sigaltstack",
        "signalfd",
        "signalfd4",
        "sigreturn",
        "socket",
        "socketcall",
        "socketpair",
        "splice",
        "stat",
        "statfs",
        "statfs64",
        "symlink",
        "symlinkat",
        "sync",
        "sync_file_range",
        "syncfs",
        "sysinfo",
        "syslog",
        "tee",
        "tgkill",
        "time",
        "timer_create",
        "timer_delete",
        "timerfd_create",
        "timerfd_gettime",
        "timerfd_settime",
        "timer_getoverrun",
        "timer_gettime",
        "timer_settime",
        "times",
        "tkill",
        "truncate",
        "umask",
        "umount",
        "umount2",
        "uname",
        "unlink",
        "unlinkat",
        "utime",
        "utimensat",
        "utimes",
        "vfork",
        "vmsplice",
        "wait4",
        "waitid",
        "write",
        "writev"
      ],
      "action": "SCMP_ACT_ALLOW"
    }
  ]
}
```

#### Apply Seccomp Profile to Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: seccomp-pod
spec:
  securityContext:
    seccompProfile:
      type: Localhost
      localhostProfile: profiles/deny-all.json
  containers:
  - name: test-container
    image: nginx
```

### Rootless Containers

Rootless containers run without root privileges, improving security.

#### Docker Rootless Setup

```bash
# Install rootless Docker
curl -fsSL https://get.docker.com/rootless | sh

# Set environment variables
export PATH=$HOME/bin:$PATH
export DOCKER_HOST=unix://$HOME/.docker/run/docker.sock

# Start rootless Docker
dockerd-rootless.sh --experimental --storage-driver vfs
```

#### Kubernetes Rootless (Experimental)

Kubernetes rootless support is experimental. For production, use:
- Pod Security Standards
- Security Contexts
- Non-root containers

#### Run Containers as Non-Root

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: non-root-pod
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 2000
  containers:
  - name: app
    image: nginx
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
        add:
        - NET_BIND_SERVICE
```

### Pod Security Standards

Kubernetes Pod Security Standards provide three policies:

#### Privileged (Most Permissive)
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: privileged-ns
  labels:
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/audit: privileged
    pod-security.kubernetes.io/warn: privileged
```

#### Baseline (Recommended)
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: baseline-ns
  labels:
    pod-security.kubernetes.io/enforce: baseline
    pod-security.kubernetes.io/audit: baseline
    pod-security.kubernetes.io/warn: baseline
```

#### Restricted (Most Secure)
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: restricted-ns
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### Network Policies

Implement network policies to control pod-to-pod communication:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: production
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

### RBAC Configuration

Implement Role-Based Access Control:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
subjects:
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

### Image Security

#### Scan Images for Vulnerabilities

```bash
# Using Trivy
trivy image nginx:latest

# Using Docker Scout
docker scout cves nginx:latest
```

#### Use Trusted Base Images

- Prefer official images from Docker Hub
- Use minimal base images (Alpine, Distroless)
- Regularly update base images
- Pin specific image versions

#### Image Pull Secrets

```bash
# Create secret
kubectl create secret docker-registry regcred \
  --docker-server=<registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  --docker-email=<email>

# Use in pod
apiVersion: v1
kind: Pod
metadata:
  name: private-reg-pod
spec:
  imagePullSecrets:
  - name: regcred
  containers:
  - name: app
    image: private-registry/app:latest
```

### Security Best Practices Summary

1. **Always run containers as non-root** when possible
2. **Use Pod Security Standards** to enforce security policies
3. **Implement Network Policies** to restrict network traffic
4. **Enable RBAC** and follow principle of least privilege
5. **Scan images** for vulnerabilities before deployment
6. **Use AppArmor/Seccomp** profiles to restrict capabilities
7. **Keep components updated** with security patches
8. **Use secrets management** (Kubernetes Secrets or external tools)
9. **Enable audit logging** for security monitoring
10. **Regularly review** and update security configurations

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kubeadm Documentation](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)
- [Container Runtime Interface (CRI)](https://kubernetes.io/docs/concepts/architecture/cri/)
- [Kubernetes Security Best Practices](https://kubernetes.io/docs/concepts/security/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)

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
