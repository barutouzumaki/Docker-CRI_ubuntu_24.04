#!/bin/bash

################################################################################
# Kubernetes and Docker Setup Script for Ubuntu 24.04
# This script automates the installation of Docker, Docker CRI, and Kubernetes
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
KUBERNETES_VERSION="1.28"
GO_VERSION="1.21.5"
INIT_CLUSTER=false
POD_NETWORK_CIDR="10.244.0.0/16"

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        print_error "Please run as root (use sudo)"
        exit 1
    fi
}

# Function to check Ubuntu version
check_ubuntu_version() {
    if [ ! -f /etc/os-release ]; then
        print_error "Cannot determine OS version"
        exit 1
    fi

    . /etc/os-release
    if [ "$ID" != "ubuntu" ]; then
        print_warn "This script is designed for Ubuntu. Detected: $ID"
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Function to update system
update_system() {
    print_info "Updating system packages..."
    apt-get update -qq
    apt-get upgrade -y -qq
    print_info "System updated successfully"
}

# Function to install prerequisites
install_prerequisites() {
    print_info "Installing prerequisites..."
    apt-get install -y -qq \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        wget \
        git \
        jq \
        vim \
        net-tools
    print_info "Prerequisites installed successfully"
}

# Function to configure kernel parameters
configure_kernel() {
    print_info "Configuring kernel parameters..."
    
    # Load required kernel modules
    cat <<EOF | tee /etc/modules-load.d/k8s.conf > /dev/null
overlay
br_netfilter
EOF

    modprobe overlay
    modprobe br_netfilter

    # Configure sysctl parameters
    cat <<EOF | tee /etc/sysctl.d/k8s.conf > /dev/null
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

    sysctl --system > /dev/null
    print_info "Kernel parameters configured successfully"
}

# Function to disable swap
disable_swap() {
    print_info "Disabling swap..."
    
    # Temporarily disable swap
    swapoff -a 2>/dev/null || true
    
    # Permanently disable swap
    sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
    
    print_info "Swap disabled successfully"
}

# Function to install Docker
install_docker() {
    print_info "Installing Docker Engine..."
    
    # Check if Docker is already installed
    if command -v docker &> /dev/null; then
        print_warn "Docker is already installed: $(docker --version)"
        read -p "Reinstall Docker? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return
        fi
    fi

    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # Add Docker repository
    . /etc/os-release
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $VERSION_CODENAME stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null

    # Install Docker
    apt-get update -qq
    apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Configure containerd for Kubernetes
    mkdir -p /etc/containerd
    containerd config default | tee /etc/containerd/config.toml > /dev/null
    sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml

    # Start and enable Docker
    systemctl restart containerd
    systemctl enable containerd
    systemctl start docker
    systemctl enable docker

    # Verify Docker installation
    if docker run --rm hello-world &> /dev/null; then
        print_info "Docker installed and verified successfully"
    else
        print_error "Docker installation verification failed"
        exit 1
    fi
}

# Function to install Go
install_go() {
    print_info "Installing Go $GO_VERSION..."
    
    if command -v go &> /dev/null; then
        print_warn "Go is already installed: $(go version)"
        return
    fi

    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    
    print_info "Go installed successfully"
}

# Function to install Docker CRI
install_docker_cri() {
    print_info "Installing Docker CRI..."
    
    # Export Go path
    export PATH=$PATH:/usr/local/go/bin
    
    # Get the script directory
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    cd "$SCRIPT_DIR"
    
    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the Docker-CRI_ubuntu_24.04 directory"
        exit 1
    fi
    
    # Download dependencies
    print_info "Downloading Docker CRI dependencies..."
    go mod download
    go mod tidy
    
    # Build Docker CRI
    print_info "Building Docker CRI binary..."
    mkdir -p bin
    go build -o bin/docker-cri ./main.go
    
    # Install Docker CRI
    cp bin/docker-cri /usr/local/bin/docker-cri
    chmod +x /usr/local/bin/docker-cri
    
    # Create required directories
    mkdir -p /var/lib/docker-cri
    mkdir -p /var/run
    
    # Create systemd service
    print_info "Creating Docker CRI systemd service..."
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

    # Enable and start Docker CRI service
    systemctl daemon-reload
    systemctl enable docker-cri
    systemctl start docker-cri
    
    # Wait for service to start
    sleep 3
    
    if systemctl is-active --quiet docker-cri; then
        print_info "Docker CRI installed and started successfully"
    else
        print_error "Docker CRI service failed to start"
        systemctl status docker-cri
        exit 1
    fi
}

# Function to install Kubernetes components
install_kubernetes() {
    print_info "Installing Kubernetes components (kubeadm, kubelet, kubectl)..."
    
    # Add Kubernetes GPG key
    mkdir -p /etc/apt/keyrings
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v${KUBERNETES_VERSION}/deb/Release.key | \
        gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    
    # Add Kubernetes repository
    echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v${KUBERNETES_VERSION}/deb/ /" | \
        tee /etc/apt/sources.list.d/kubernetes.list > /dev/null
    
    # Install Kubernetes components
    apt-get update -qq
    apt-get install -y -qq kubelet kubeadm kubectl
    apt-mark hold kubelet kubeadm kubectl
    
    # Configure kubelet
    mkdir -p /var/lib/kubelet
    cat <<EOF | tee /var/lib/kubelet/config.yaml > /dev/null
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
EOF

    print_info "Kubernetes components installed successfully"
}

# Function to initialize Kubernetes cluster
init_kubernetes_cluster() {
    print_info "Initializing Kubernetes cluster..."
    
    # Check if cluster is already initialized
    if [ -f /etc/kubernetes/admin.conf ]; then
        print_warn "Kubernetes cluster appears to be already initialized"
        read -p "Reinitialize cluster? This will reset the existing cluster. (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return
        fi
        kubeadm reset -f
    fi
    
    # Initialize cluster
    kubeadm init \
        --pod-network-cidr=${POD_NETWORK_CIDR} \
        --cri-socket=unix:///var/run/docker-cri.sock \
        --ignore-preflight-errors=Swap
    
    # Configure kubectl for root user
    mkdir -p /root/.kube
    cp -i /etc/kubernetes/admin.conf /root/.kube/config
    chown root:root /root/.kube/config
    
    # Configure kubectl for current user (if not root)
    if [ -n "$SUDO_USER" ]; then
        mkdir -p /home/$SUDO_USER/.kube
        cp -i /etc/kubernetes/admin.conf /home/$SUDO_USER/.kube/config
        chown -R $SUDO_USER:$SUDO_USER /home/$SUDO_USER/.kube
    fi
    
    # Install Flannel network plugin
    print_info "Installing Flannel pod network..."
    export KUBECONFIG=/etc/kubernetes/admin.conf
    kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
    
    print_info "Waiting for pods to be ready..."
    sleep 10
    
    # Wait for core pods to be ready
    kubectl wait --for=condition=ready pod --all -n kube-system --timeout=300s || true
    
    print_info "Kubernetes cluster initialized successfully"
    
    # Display join command
    print_info "To join worker nodes, use the following command:"
    kubeadm token create --print-join-command
}

# Function to install crictl
install_crictl() {
    print_info "Installing crictl..."
    
    if command -v crictl &> /dev/null; then
        print_warn "crictl is already installed"
        return
    fi
    
    CRICTL_VERSION="v1.28.0"
    wget -q https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz
    tar zxvf crictl-${CRICTL_VERSION}-linux-amd64.tar.gz -C /usr/local/bin
    rm crictl-${CRICTL_VERSION}-linux-amd64.tar.gz
    
    print_info "crictl installed successfully"
}

# Function to verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    # Check Docker
    if docker --version &> /dev/null; then
        print_info "✓ Docker: $(docker --version | cut -d' ' -f3 | tr -d ',')"
    else
        print_error "✗ Docker verification failed"
    fi
    
    # Check Docker CRI
    if systemctl is-active --quiet docker-cri; then
        print_info "✓ Docker CRI: Service is running"
    else
        print_error "✗ Docker CRI service is not running"
    fi
    
    # Check Kubernetes components
    if kubeadm version &> /dev/null; then
        print_info "✓ kubeadm: $(kubeadm version -o short 2>/dev/null || echo 'installed')"
    else
        print_error "✗ kubeadm verification failed"
    fi
    
    if kubectl version --client &> /dev/null; then
        print_info "✓ kubectl: $(kubectl version --client --short 2>/dev/null | cut -d' ' -f3)"
    else
        print_error "✗ kubectl verification failed"
    fi
    
    if systemctl is-active --quiet kubelet; then
        print_info "✓ kubelet: Service is running"
    else
        print_warn "⚠ kubelet service is not running (normal if cluster not initialized)"
    fi
}

# Main function
main() {
    echo "=========================================="
    echo "Kubernetes and Docker Setup Script"
    echo "Ubuntu 24.04"
    echo "=========================================="
    echo ""
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --init-cluster)
                INIT_CLUSTER=true
                shift
                ;;
            --pod-network-cidr)
                POD_NETWORK_CIDR="$2"
                shift 2
                ;;
            --k8s-version)
                KUBERNETES_VERSION="$2"
                shift 2
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Usage: $0 [--init-cluster] [--pod-network-cidr CIDR] [--k8s-version VERSION]"
                exit 1
                ;;
        esac
    done
    
    # Check prerequisites
    check_root
    check_ubuntu_version
    
    # Installation steps
    update_system
    install_prerequisites
    configure_kernel
    disable_swap
    install_docker
    install_go
    install_docker_cri
    install_kubernetes
    install_crictl
    
    # Verify installation
    verify_installation
    
    # Initialize cluster if requested
    if [ "$INIT_CLUSTER" = true ]; then
        init_kubernetes_cluster
    fi
    
    echo ""
    echo "=========================================="
    echo "Installation Complete!"
    echo "=========================================="
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. Verify Docker:"
    echo "   sudo docker run hello-world"
    echo ""
    echo "2. Verify Docker CRI:"
    echo "   sudo crictl --runtime-endpoint unix:///var/run/docker-cri.sock version"
    echo ""
    if [ "$INIT_CLUSTER" = false ]; then
        echo "3. Initialize Kubernetes cluster:"
        echo "   sudo kubeadm init --pod-network-cidr=${POD_NETWORK_CIDR} --cri-socket=unix:///var/run/docker-cri.sock"
        echo "   mkdir -p \$HOME/.kube"
        echo "   sudo cp -i /etc/kubernetes/admin.conf \$HOME/.kube/config"
        echo "   sudo chown \$(id -u):\$(id -g) \$HOME/.kube/config"
        echo "   kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml"
        echo ""
    fi
    echo "4. Check service status:"
    echo "   sudo systemctl status docker"
    echo "   sudo systemctl status docker-cri"
    echo "   sudo systemctl status kubelet"
    echo ""
    echo "5. View logs:"
    echo "   sudo journalctl -u docker-cri -f"
    echo "   sudo journalctl -u kubelet -f"
    echo ""
}

# Run main function
main "$@"


