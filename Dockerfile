# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o docker-cri ./main.go

# Runtime stage
FROM ubuntu:24.04

# Install Docker CLI (for Docker socket access)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install Docker CLI
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && \
    apt-get install -y --no-install-recommends docker-cli && \
    rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /build/docker-cri /usr/local/bin/docker-cri

# Create necessary directories
RUN mkdir -p /var/lib/docker-cri /var/run

# Set permissions
RUN chmod +x /usr/local/bin/docker-cri

# Expose socket directory (mounted as volume)
VOLUME ["/var/run"]

WORKDIR /

ENTRYPOINT ["/usr/local/bin/docker-cri"]

