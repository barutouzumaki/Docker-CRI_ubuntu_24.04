# Docker Interview Questions

This document contains commonly asked Docker interview questions organized by topic.

## Table of Contents

- [Basic Concepts](#basic-concepts)
- [Docker Architecture](#docker-architecture)
- [Images & Containers](#images--containers)
- [Dockerfile](#dockerfile)
- [Networking](#networking)
- [Storage & Volumes](#storage--volumes)
- [Docker Compose](#docker-compose)
- [Security](#security)
- [Troubleshooting](#troubleshooting)
- [Advanced Topics](#advanced-topics)

---

## Basic Concepts

### 1. What is Docker and why is it used?

**Answer:**
Docker is a platform for developing, shipping, and running applications using containerization. It's used because:
- **Consistency**: "Works on my machine" problem solved - same environment everywhere
- **Isolation**: Applications run in isolated containers
- **Portability**: Run anywhere (Linux, Windows, cloud, on-premises)
- **Efficiency**: Lightweight compared to VMs, shares host OS kernel
- **Scalability**: Easy to scale applications horizontally
- **Fast deployment**: Containers start in seconds vs minutes for VMs
- **Version control**: Images can be versioned and tracked

### 2. What is the difference between Docker and Virtual Machines?

**Answer:**

| Aspect | Docker (Containers) | Virtual Machines |
|--------|-------------------|-------------------|
| **OS** | Shares host OS kernel | Full guest OS per VM |
| **Size** | MBs (lightweight) | GBs (heavy) |
| **Startup** | Seconds | Minutes |
| **Isolation** | Process-level | Hardware-level |
| **Resource Usage** | Lower overhead | Higher overhead |
| **Portability** | High (OS-specific) | Medium |
| **Security** | Process isolation | Strong isolation |

**Docker** uses the host OS kernel and isolates at the process level.
**VMs** virtualize hardware and run a complete OS.

### 3. What is a Docker Image?

**Answer:**
A Docker image is a read-only template used to create containers. It contains:
- **Application code**: Your application files
- **Dependencies**: Libraries, packages, runtime
- **Configuration**: Environment variables, settings
- **Layers**: Multiple read-only layers stacked together

**Characteristics:**
- Immutable (read-only)
- Built from Dockerfile
- Stored in registries (Docker Hub, private registries)
- Can be versioned with tags
- Layers are cached and shared between images

**Example:**
```bash
# Pull an image
docker pull nginx:latest

# List images
docker images

# Inspect image
docker inspect nginx:latest
```

### 4. What is a Docker Container?

**Answer:**
A Docker container is a runnable instance of an image. It:
- **Runs applications**: Executes the application from the image
- **Has writable layer**: Top layer is writable (changes are ephemeral)
- **Is isolated**: Has its own filesystem, network, and process space
- **Is ephemeral**: Can be created, started, stopped, and deleted

**Container Lifecycle:**
1. **Created**: Container created but not started
2. **Running**: Container is active
3. **Paused**: Execution paused
4. **Stopped**: Container stopped
5. **Removed**: Container deleted

**Example:**
```bash
# Create and run container
docker run -d --name my-nginx nginx:latest

# List running containers
docker ps

# Stop container
docker stop my-nginx

# Remove container
docker rm my-nginx
```

### 5. What is the difference between Docker Image and Container?

**Answer:**

| Aspect | Image | Container |
|--------|-------|-----------|
| **Type** | Template/Blueprint | Instance |
| **State** | Read-only | Read-write (top layer) |
| **Lifecycle** | Built once, reused | Created, started, stopped |
| **Storage** | Registry/filesystem | Runtime instance |
| **Changes** | Immutable | Ephemeral (unless committed) |
| **Count** | One image | Many containers from one image |

**Analogy:** Image is like a class, Container is like an object instance.

**Example:**
```bash
# Image (template)
docker pull nginx:latest

# Multiple containers from same image
docker run -d --name web1 nginx:latest
docker run -d --name web2 nginx:latest
docker run -d --name web3 nginx:latest
```

### 6. What is Docker Hub?

**Answer:**
Docker Hub is Docker's public registry service for storing and sharing Docker images. It provides:
- **Public repositories**: Free public image storage
- **Private repositories**: Paid private image storage
- **Automated builds**: Build images from GitHub/Bitbucket
- **Webhooks**: Trigger actions on image updates
- **Organizations**: Team collaboration features
- **Official images**: Curated official images (nginx, mysql, etc.)

**Usage:**
```bash
# Pull from Docker Hub
docker pull nginx

# Push to Docker Hub
docker tag myapp:latest username/myapp:latest
docker push username/myapp:latest

# Login to Docker Hub
docker login
```

### 7. What is a Dockerfile?

**Answer:**
A Dockerfile is a text file containing instructions for building a Docker image. It defines:
- **Base image**: Starting point (FROM)
- **Commands**: Steps to build the image (RUN, COPY, etc.)
- **Configuration**: Environment variables, exposed ports
- **Entry point**: Default command to run

**Common Instructions:**
- `FROM`: Base image
- `RUN`: Execute commands
- `COPY/ADD`: Copy files
- `WORKDIR`: Set working directory
- `ENV`: Set environment variables
- `EXPOSE`: Expose ports
- `CMD/ENTRYPOINT`: Default command

**Example:**
```dockerfile
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y nginx
COPY index.html /var/www/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 8. Explain Docker layers and how they work.

**Answer:**
Docker images are built using a layered filesystem. Each instruction in a Dockerfile creates a new layer:
- **Layers are read-only**: Once created, layers cannot be changed
- **Layers are cached**: Unchanged layers are reused
- **Copy-on-write**: Containers get a writable layer on top
- **Efficient storage**: Shared layers save disk space

**Benefits:**
- **Faster builds**: Cached layers don't rebuild
- **Smaller images**: Shared base layers
- **Version control**: Each layer is a diff

**Example:**
```dockerfile
FROM ubuntu:24.04          # Layer 1: Base OS
RUN apt-get update         # Layer 2: Package lists
RUN apt-get install nginx  # Layer 3: Nginx installation
COPY app.conf /etc/nginx/   # Layer 4: Configuration
```

**Viewing layers:**
```bash
docker history nginx:latest
docker inspect nginx:latest | grep -A 10 "RootFS"
```

### 9. What is the difference between CMD and ENTRYPOINT?

**Answer:**

| Aspect | CMD | ENTRYPOINT |
|--------|-----|------------|
| **Purpose** | Default command/arguments | Main command (always executed) |
| **Override** | Can be overridden | Cannot be overridden (can append) |
| **Arguments** | Can be replaced | Arguments are appended |
| **Use Case** | Default command | Main application command |

**CMD:**
- Provides default command
- Can be completely overridden: `docker run image <new-command>`
- Used for default arguments

**ENTRYPOINT:**
- Main command always runs
- Arguments are appended: `docker run image <args>` → `entrypoint <args>`
- Used for main application executable

**Example:**
```dockerfile
# CMD - can be overridden
FROM ubuntu
CMD ["echo", "Hello World"]
# docker run image → "Hello World"
# docker run image echo "Hi" → "Hi"

# ENTRYPOINT - always runs
FROM ubuntu
ENTRYPOINT ["echo"]
CMD ["Hello World"]
# docker run image → "Hello World"
# docker run image "Hi" → "Hi"
```

### 10. What is Docker Compose?

**Answer:**
Docker Compose is a tool for defining and running multi-container Docker applications. It:
- **Defines services**: Multiple containers in a YAML file
- **Manages dependencies**: Starts services in correct order
- **Networks**: Creates isolated networks for services
- **Volumes**: Manages shared volumes
- **Single command**: `docker-compose up` starts everything

**Use Cases:**
- Local development environments
- Multi-container applications
- Testing environments
- CI/CD pipelines

**Example docker-compose.yml:**
```yaml
version: '3.8'
services:
  web:
    image: nginx
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html
  db:
    image: postgres
    environment:
      POSTGRES_PASSWORD: password
    volumes:
      - db-data:/var/lib/postgresql/data

volumes:
  db-data:
```

---

## Docker Architecture

### 11. Explain Docker's architecture.

**Answer:**
Docker follows a client-server architecture:

**Components:**
1. **Docker Client (CLI)**: Command-line interface (`docker` command)
2. **Docker Daemon (dockerd)**: Background service managing containers
3. **Docker API**: REST API for Docker operations
4. **Container Runtime**: containerd, runc for running containers

**Architecture Flow:**
```
Docker Client → Docker API → Docker Daemon → containerd → runc → Container
```

**Docker Daemon Responsibilities:**
- Building images
- Running containers
- Managing networks
- Managing volumes
- Managing images

**Communication:**
- Client and daemon can be on same machine or remote
- Communication via REST API over Unix socket or TCP

### 12. What is containerd?

**Answer:**
containerd is an industry-standard container runtime. It:
- **Manages containers**: Lifecycle management (create, start, stop, delete)
- **Manages images**: Image pull, push, storage
- **Manages snapshots**: Filesystem snapshots for layers
- **Low-level runtime**: Used by Docker and Kubernetes
- **OCI-compliant**: Follows Open Container Initiative standards

**Relationship with Docker:**
- Docker uses containerd internally
- containerd uses runc to run containers
- Docker provides higher-level features (networking, volumes)

**Direct Usage:**
```bash
# Using ctr (containerd CLI)
ctr images pull docker.io/library/nginx:latest
ctr containers create docker.io/library/nginx:latest nginx
ctr containers start nginx
```

### 13. What is runc?

**Answer:**
runc is a CLI tool for spawning and running containers according to the OCI specification. It:
- **Creates containers**: Spawns containers from OCI bundles
- **Low-level**: Direct interface to Linux kernel features
- **OCI-compliant**: Implements OCI Runtime Specification
- **Used by containerd**: containerd calls runc to run containers

**Features:**
- Namespace isolation
- cgroup management
- Root filesystem setup
- Process execution

**Container Runtime Stack:**
```
Docker → containerd → runc → Container
```

### 14. What are Docker namespaces?

**Answer:**
Docker uses Linux namespaces to provide isolation. Namespaces isolate:
- **PID namespace**: Process IDs (isolated process tree)
- **Network namespace**: Network interfaces, routing tables
- **Mount namespace**: Filesystem mount points
- **IPC namespace**: Inter-process communication
- **UTS namespace**: Hostname and domain name
- **User namespace**: User and group IDs

**Benefits:**
- Process isolation
- Network isolation
- Filesystem isolation
- Security boundaries

**Viewing namespaces:**
```bash
# Container namespaces
docker inspect <container> | grep -i namespace

# Host namespaces
ls -la /proc/self/ns/
```

### 15. What are cgroups in Docker?

**Answer:**
cgroups (control groups) limit and account for resource usage. Docker uses cgroups to:
- **Limit CPU**: CPU shares, CPU quota
- **Limit Memory**: Memory limits, OOM killer
- **Limit I/O**: Disk I/O throttling
- **Account resources**: Track resource usage

**cgroup v2 (Ubuntu 24.04):**
- Unified hierarchy
- Better resource management
- Improved performance

**Setting limits:**
```bash
# CPU and memory limits
docker run -d --name app \
  --cpus="1.5" \
  --memory="512m" \
  nginx:latest

# View cgroup info
docker stats app
```

### 16. What is the difference between COPY and ADD in Dockerfile?

**Answer:**

| Aspect | COPY | ADD |
|--------|------|-----|
| **Source** | Local files/directories | Local files + URLs + tar extraction |
| **URLs** | Not supported | Downloads from URLs |
| **Tar extraction** | No | Yes (automatic) |
| **Best Practice** | Preferred for files | Use for URLs/tar |
| **Clarity** | More explicit | Less explicit |

**COPY:**
```dockerfile
# Simple file copy
COPY app.py /app/
COPY requirements.txt /app/
```

**ADD:**
```dockerfile
# Can download and extract
ADD https://example.com/file.tar.gz /tmp/
# Automatically extracts tar.gz

# Can download
ADD https://example.com/file.txt /tmp/
```

**Best Practice:** Use COPY for local files, ADD only when you need URL download or tar extraction.

### 17. What is a multi-stage build in Docker?

**Answer:**
Multi-stage builds allow using multiple FROM statements in a Dockerfile to:
- **Reduce image size**: Only include necessary files in final image
- **Separate build and runtime**: Build tools not in final image
- **Optimize layers**: Smaller final image

**Example:**
```dockerfile
# Stage 1: Build
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .
CMD ["./myapp"]
```

**Benefits:**
- Final image only contains compiled binary
- No build tools in production image
- Significantly smaller images

### 18. What is Docker Swarm?

**Answer:**
Docker Swarm is Docker's native clustering and orchestration solution. It provides:
- **Cluster management**: Multiple Docker hosts as single cluster
- **Service orchestration**: Deploy and manage services
- **Load balancing**: Automatic load balancing
- **Scaling**: Scale services up/down
- **Rolling updates**: Update services without downtime

**Swarm Concepts:**
- **Swarm**: Cluster of Docker nodes
- **Node**: Manager or worker node
- **Service**: Definition of tasks to execute
- **Task**: Container running on a node
- **Stack**: Collection of services

**Example:**
```bash
# Initialize swarm
docker swarm init

# Create service
docker service create --replicas 3 --name web nginx:latest

# Scale service
docker service scale web=5
```

### 19. What is the difference between Docker and Docker Compose?

**Answer:**

| Aspect | Docker | Docker Compose |
|--------|--------|----------------|
| **Scope** | Single container | Multiple containers |
| **Configuration** | Command line | YAML file |
| **Use Case** | Simple apps, testing | Multi-container apps |
| **Networking** | Manual setup | Automatic |
| **Orchestration** | Manual | Declarative |

**Docker:**
- Run individual containers
- Command-line based
- Manual networking/volumes
- Good for simple applications

**Docker Compose:**
- Define multi-container apps
- YAML-based configuration
- Automatic networking/volumes
- Good for complex applications

### 20. What is a Docker volume?

**Answer:**
A Docker volume is a mechanism for persisting data generated by containers. It:
- **Persists data**: Data survives container lifecycle
- **Shared storage**: Can be shared between containers
- **Managed by Docker**: Docker manages volume lifecycle
- **Backup/restore**: Can be backed up and restored

**Volume Types:**
- **Named volumes**: Managed by Docker (recommended)
- **Anonymous volumes**: Created automatically
- **Bind mounts**: Mount host directories
- **tmpfs mounts**: In-memory storage

**Example:**
```bash
# Create volume
docker volume create my-volume

# Use in container
docker run -d -v my-volume:/data nginx:latest

# List volumes
docker volume ls

# Inspect volume
docker volume inspect my-volume
```

---

## Images & Containers

### 21. What is the difference between docker run and docker start?

**Answer:**

| Aspect | docker run | docker start |
|--------|-----------|--------------|
| **Action** | Creates and starts new container | Starts existing stopped container |
| **Container** | New container | Existing container |
| **Image** | Requires image | Uses existing container |
| **Options** | Can specify all options | Limited options |
| **ID** | New container ID | Same container ID |

**docker run:**
- Creates new container from image
- Can specify all container options
- Container gets new ID
- Use for new containers

**docker start:**
- Starts previously stopped container
- Container keeps same ID and configuration
- Use to resume stopped containers

**Example:**
```bash
# Create and start new container
docker run -d --name my-nginx nginx:latest

# Stop container
docker stop my-nginx

# Start existing container
docker start my-nginx

# Remove container
docker rm my-nginx
```

### 22. What is the difference between docker stop and docker kill?

**Answer:**

| Aspect | docker stop | docker kill |
|--------|-------------|-------------|
| **Signal** | SIGTERM (graceful) | SIGKILL (forceful) |
| **Grace Period** | 10 seconds default | Immediate |
| **Process Cleanup** | Allows cleanup | No cleanup |
| **Use Case** | Normal shutdown | Force stop |

**docker stop:**
- Sends SIGTERM signal
- Waits for graceful shutdown (default 10s)
- Then sends SIGKILL if still running
- Allows application to cleanup

**docker kill:**
- Sends SIGKILL signal immediately
- Forceful termination
- No cleanup opportunity
- Use when container is unresponsive

**Example:**
```bash
# Graceful stop
docker stop my-container

# Force kill
docker kill my-container

# Custom signal
docker kill --signal=SIGTERM my-container
```

### 23. What is the difference between docker exec and docker attach?

**Answer:**

| Aspect | docker exec | docker attach |
|--------|-------------|---------------|
| **New Process** | Creates new process | Attaches to main process |
| **Multiple Sessions** | Multiple exec sessions | Single attach session |
| **Exit Behavior** | Exiting doesn't stop container | Exiting stops container (if main process) |
| **Use Case** | Run commands, debug | Interact with main process |

**docker exec:**
- Runs new command in running container
- Creates new process
- Multiple exec sessions possible
- Exiting doesn't affect container

**docker attach:**
- Attaches to main container process (PID 1)
- Single session
- Exiting may stop container
- Use for interactive applications

**Example:**
```bash
# Execute command in running container
docker exec -it my-container /bin/bash

# Attach to main process
docker attach my-container

# Detach without stopping (Ctrl+P, Ctrl+Q)
```

### 24. How do you remove unused Docker images, containers, and volumes?

**Answer:**

**Remove unused resources:**
```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove unused networks
docker network prune

# Remove everything unused (containers, networks, images, build cache)
docker system prune

# Remove everything including volumes
docker system prune -a --volumes
```

**Manual removal:**
```bash
# Remove specific image
docker rmi <image-id>

# Remove all images
docker rmi $(docker images -q)

# Remove stopped containers
docker rm $(docker ps -a -q -f status=exited)

# Remove unused volumes
docker volume rm $(docker volume ls -q -f dangling=true)
```

**Best Practice:**
- Regularly clean up unused resources
- Use `docker system prune` periodically
- Be careful with `-a` flag (removes all unused images)

### 25. What is the difference between docker build and docker commit?

**Answer:**

| Aspect | docker build | docker commit |
|--------|--------------|---------------|
| **Method** | From Dockerfile | From running container |
| **Reproducibility** | Reproducible | Not reproducible |
| **Best Practice** | Recommended | Not recommended |
| **Use Case** | Standard builds | Quick saves, debugging |

**docker build:**
- Builds image from Dockerfile
- Reproducible and version controlled
- Best practice for creating images
- Can be automated in CI/CD

**docker commit:**
- Creates image from container changes
- Not reproducible
- Loses build history
- Use only for quick saves or debugging

**Example:**
```bash
# Build from Dockerfile (recommended)
docker build -t myapp:latest .

# Commit from container (not recommended)
docker commit my-container myapp:latest
```

### 26. What is a Docker tag and how is it used?

**Answer:**
A Docker tag is a label attached to images for versioning and identification. It:
- **Identifies versions**: `nginx:1.21`, `nginx:latest`
- **Organizes images**: Different tags for different versions
- **Default tag**: `latest` if not specified
- **Format**: `<repository>:<tag>`

**Tagging:**
```bash
# Tag an image
docker tag nginx:latest nginx:1.21
docker tag nginx:latest myregistry/nginx:v1.21

# List tags
docker images nginx

# Push specific tag
docker push myregistry/nginx:v1.21
```

**Best Practices:**
- Use semantic versioning: `v1.2.3`
- Use `latest` carefully (can be misleading)
- Tag images with build numbers or commit hashes
- Tag for different environments: `dev`, `staging`, `prod`

### 27. How do you copy files from host to container and vice versa?

**Answer:**

**docker cp command:**
```bash
# Copy from host to container
docker cp /path/to/local/file container-name:/path/to/destination

# Copy from container to host
docker cp container-name:/path/to/file /path/to/local/destination

# Copy directory
docker cp /local/dir container-name:/destination/dir
```

**Example:**
```bash
# Copy file to container
docker cp app.conf my-container:/etc/nginx/conf.d/

# Copy file from container
docker cp my-container:/var/log/app.log ./logs/

# Copy directory
docker cp ./config my-container:/app/config
```

**Note:** `docker cp` works on stopped containers too, but filesystem must be accessible.

---

## Dockerfile

### 28. What is .dockerignore and why is it used?

**Answer:**
`.dockerignore` is a file that specifies files/directories to exclude from Docker build context. It:
- **Reduces build context**: Smaller context = faster builds
- **Prevents sensitive data**: Excludes secrets, credentials
- **Improves performance**: Less data to send to Docker daemon
- **Similar to .gitignore**: Same syntax

**Example .dockerignore:**
```
# Git files
.git
.gitignore

# Dependencies (if installing in container)
node_modules
vendor

# Build artifacts
dist
build
*.o

# IDE files
.vscode
.idea

# Secrets
.env
*.pem
secrets/

# Documentation
README.md
docs/
```

**Benefits:**
- Faster builds
- Smaller images
- Security (exclude secrets)
- Cleaner build context

### 29. What is the difference between RUN, CMD, and ENTRYPOINT?

**Answer:**

| Instruction | When Executed | Purpose | Overridable |
|-------------|---------------|---------|-------------|
| **RUN** | During image build | Execute commands in build | N/A |
| **CMD** | Container start | Default command/arguments | Yes |
| **ENTRYPOINT** | Container start | Main command | No (args appended) |

**RUN:**
- Executes during image build
- Creates new layer
- Use for installing packages, setting up environment

**CMD:**
- Provides default command
- Can be overridden: `docker run image <new-command>`
- Use for default arguments

**ENTRYPOINT:**
- Main command always runs
- Arguments appended: `docker run image <args>`
- Use for main executable

**Example:**
```dockerfile
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y curl  # Build time
ENTRYPOINT ["curl"]                            # Always runs
CMD ["-I", "https://example.com"]              # Default args
# docker run image → curl -I https://example.com
# docker run image -v → curl -v
```

### 30. How do you optimize Dockerfile for smaller image size?

**Answer:**

**Best Practices:**

1. **Use multi-stage builds:**
   ```dockerfile
   FROM golang:1.21 AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o app

   FROM alpine:latest
   COPY --from=builder /app/app .
   CMD ["./app"]
   ```

2. **Use minimal base images:**
   ```dockerfile
   # Instead of ubuntu:latest
   FROM alpine:latest
   # Or distroless
   FROM gcr.io/distroless/base
   ```

3. **Combine RUN commands:**
   ```dockerfile
   # Bad
   RUN apt-get update
   RUN apt-get install -y package1
   RUN apt-get install -y package2

   # Good
   RUN apt-get update && \
       apt-get install -y package1 package2 && \
       rm -rf /var/lib/apt/lists/*
   ```

4. **Order instructions (cache optimization):**
   ```dockerfile
   # Copy dependency files first
   COPY requirements.txt .
   RUN pip install -r requirements.txt
   # Copy application code last
   COPY . .
   ```

5. **Remove unnecessary files:**
   ```dockerfile
   RUN apt-get install -y package && \
       rm -rf /var/lib/apt/lists/*
   ```

6. **Use .dockerignore:**
   - Exclude unnecessary files from build context

### 31. What is ARG and ENV in Dockerfile?

**Answer:**

| Aspect | ARG | ENV |
|--------|-----|-----|
| **Scope** | Build-time only | Build-time and runtime |
| **Accessibility** | Not in container | Available in container |
| **Use Case** | Build-time variables | Runtime configuration |
| **Override** | `--build-arg` | `-e` or `--env` |

**ARG:**
- Build-time variable
- Not available in running container
- Use for build configuration
- Example: Version numbers, build flags

**ENV:**
- Available at build-time and runtime
- Accessible in container as environment variable
- Use for runtime configuration
- Example: Database URLs, API keys

**Example:**
```dockerfile
ARG BUILD_VERSION=1.0
ENV APP_VERSION=${BUILD_VERSION}
ENV DB_HOST=localhost

# Build with: docker build --build-arg BUILD_VERSION=2.0 .
# Runtime: DB_HOST available as env var
```

### 32. What is WORKDIR in Dockerfile?

**Answer:**
WORKDIR sets the working directory for subsequent instructions. It:
- **Sets directory**: Changes current directory
- **Creates if missing**: Creates directory if it doesn't exist
- **Affects commands**: All subsequent commands run in this directory
- **Can be used multiple times**: Can change directory multiple times

**Example:**
```dockerfile
FROM ubuntu:24.04
WORKDIR /app
RUN pwd  # Output: /app
COPY . .  # Copies to /app
WORKDIR /app/src
RUN pwd  # Output: /app/src
```

**Benefits:**
- Cleaner paths (no need for absolute paths)
- Consistent working directory
- Better organization

**vs RUN cd:**
```dockerfile
# Bad - doesn't persist
RUN cd /app && do something

# Good - persists
WORKDIR /app
RUN do something
```

---

## Networking

### 33. What are the different Docker network types?

**Answer:**
Docker provides several network drivers:

1. **bridge** (Default):
   - Default network for containers
   - Containers on same bridge can communicate
   - Isolated from host network
   - Use case: Most applications

2. **host**:
   - Removes network isolation
   - Container uses host's network stack
   - No port mapping needed
   - Use case: High performance, network debugging

3. **none**:
   - Disables all networking
   - Container has no network interfaces
   - Use case: Isolated containers

4. **overlay**:
   - Multi-host networking
   - Connects Docker daemons across hosts
   - Use case: Docker Swarm, multi-host

5. **macvlan**:
   - Assigns MAC address to container
   - Appears as physical device on network
   - Use case: Legacy applications, VLANs

**Example:**
```bash
# Create bridge network
docker network create my-network

# Run container on specific network
docker run -d --network my-network nginx

# Use host network
docker run -d --network host nginx
```

### 34. How do you expose container ports in Docker?

**Answer:**

**Port mapping:**
```bash
# Map host port to container port
docker run -p 8080:80 nginx

# Map specific host IP
docker run -p 127.0.0.1:8080:80 nginx

# Map all interfaces
docker run -p 0.0.0.0:8080:80 nginx

# Map random host port
docker run -p 80 nginx

# Map multiple ports
docker run -p 8080:80 -p 8443:443 nginx

# Publish all exposed ports
docker run -P nginx
```

**EXPOSE in Dockerfile:**
```dockerfile
FROM nginx
EXPOSE 80
EXPOSE 443
# Or
EXPOSE 80 443
```

**Note:**
- `EXPOSE` documents ports (doesn't publish)
- `-p` actually publishes ports
- `-P` publishes all exposed ports to random host ports

### 35. What is the difference between EXPOSE and -p in Docker?

**Answer:**

| Aspect | EXPOSE | -p / --publish |
|--------|--------|----------------|
| **When** | Dockerfile instruction | Runtime flag |
| **Purpose** | Documentation | Actual port mapping |
| **Effect** | Documents intent | Publishes port |
| **Required** | No | Yes (to access from host) |

**EXPOSE:**
- Dockerfile instruction
- Documents which ports container uses
- Doesn't actually publish ports
- Used with `-P` flag to publish all

**-p / --publish:**
- Runtime flag
- Actually maps host port to container port
- Required to access container from host
- Format: `host-port:container-port`

**Example:**
```dockerfile
# Dockerfile
EXPOSE 80
```

```bash
# Without -p: port not accessible from host
docker run nginx

# With -p: port accessible
docker run -p 8080:80 nginx
```

---

## Storage & Volumes

### 36. What is the difference between volumes and bind mounts?

**Answer:**

| Aspect | Volumes | Bind Mounts |
|--------|---------|-------------|
| **Location** | Managed by Docker | Host filesystem path |
| **Portability** | Portable | Host-specific |
| **Performance** | Better (Linux) | Depends on host FS |
| **Backup** | Easier | Manual |
| **Use Case** | Production | Development |

**Volumes:**
- Managed by Docker (`/var/lib/docker/volumes/`)
- Portable across hosts
- Better performance on Linux
- Can be backed up/restored easily
- Recommended for production

**Bind Mounts:**
- Direct host filesystem path
- Host-specific (not portable)
- Good for development
- Can mount any host directory

**Example:**
```bash
# Volume (recommended)
docker run -v my-volume:/data nginx

# Bind mount
docker run -v /host/path:/container/path nginx

# Named volume
docker volume create my-volume
docker run -v my-volume:/data nginx
```

### 37. What is a tmpfs mount in Docker?

**Answer:**
tmpfs mount stores data in container's memory (RAM). It:
- **In-memory storage**: Data stored in RAM
- **Fast access**: Very fast read/write
- **Ephemeral**: Data lost when container stops
- **Size limit**: Can specify size limit

**Use Cases:**
- Temporary files
- Cache data
- Sensitive data (cleared on stop)
- High-performance temporary storage

**Example:**
```bash
# tmpfs mount
docker run --tmpfs /tmp nginx

# With size limit
docker run --tmpfs /tmp:rw,noexec,nosuid,size=100m nginx
```

**Characteristics:**
- Faster than disk
- Limited by available RAM
- Data not persisted
- Good for temporary data

### 38. How do you backup and restore Docker volumes?

**Answer:**

**Backup:**
```bash
# Backup volume
docker run --rm \
  -v my-volume:/data \
  -v $(pwd):/backup \
  ubuntu tar czf /backup/backup.tar.gz /data

# Or using volume name
docker run --rm \
  -v my-volume:/data:ro \
  -v $(pwd):/backup \
  ubuntu tar czf /backup/backup.tar.gz /data
```

**Restore:**
```bash
# Restore volume
docker run --rm \
  -v my-volume:/data \
  -v $(pwd):/backup \
  ubuntu tar xzf /backup/backup.tar.gz -C /data

# Or create new volume and restore
docker volume create new-volume
docker run --rm \
  -v new-volume:/data \
  -v $(pwd):/backup \
  ubuntu tar xzf /backup/backup.tar.gz -C /data
```

**Using docker cp:**
```bash
# Backup
docker run -d --name temp -v my-volume:/data ubuntu sleep 3600
docker cp temp:/data ./backup
docker rm -f temp

# Restore
docker run -d --name temp -v my-volume:/data ubuntu sleep 3600
docker cp ./backup/. temp:/data
docker rm -f temp
```

---

## Docker Compose

### 39. What is the difference between docker-compose up and docker-compose start?

**Answer:**

| Aspect | docker-compose up | docker-compose start |
|--------|------------------|---------------------|
| **Action** | Creates and starts | Starts existing |
| **Build** | Builds if needed | Doesn't build |
| **Logs** | Shows logs | Doesn't show logs |
| **Use Case** | First run, development | Restart stopped services |

**docker-compose up:**
- Creates containers if they don't exist
- Builds images if needed
- Starts all services
- Shows logs (use `-d` for detached)
- Use for first run or after changes

**docker-compose start:**
- Starts existing stopped containers
- Doesn't create new containers
- Doesn't build images
- Doesn't show logs
- Use to restart stopped services

**Example:**
```bash
# Create and start (with logs)
docker-compose up

# Create and start (detached)
docker-compose up -d

# Start existing containers
docker-compose start

# Stop containers
docker-compose stop

# Stop and remove
docker-compose down
```

### 40. How do you scale services in Docker Compose?

**Answer:**

**Scale services:**
```bash
# Scale specific service
docker-compose up -d --scale web=3

# Scale multiple services
docker-compose up -d --scale web=3 --scale worker=2
```

**Note:** 
- Service must not have `ports` mapped to single host port
- Use port ranges or remove port mapping
- Load balancing may be needed

**Example docker-compose.yml:**
```yaml
version: '3.8'
services:
  web:
    image: nginx
    # Remove single port mapping for scaling
    # ports:
    #   - "80:80"
    deploy:
      replicas: 3
```

**With load balancer:**
```yaml
version: '3.8'
services:
  nginx:
    image: nginx
    ports:
      - "80:80"
    depends_on:
      - web
  web:
    image: myapp
    # No port mapping
```

**Limitations:**
- Docker Compose scaling is basic
- For production scaling, use Docker Swarm or Kubernetes
- Compose scaling doesn't provide load balancing

---

## Security

### 41. How do you secure Docker containers?

**Answer:**

**Best Practices:**

1. **Use minimal base images:**
   ```dockerfile
   FROM alpine:latest
   # Or distroless
   FROM gcr.io/distroless/base
   ```

2. **Run as non-root:**
   ```dockerfile
   RUN adduser -D appuser
   USER appuser
   ```

3. **Scan images for vulnerabilities:**
   ```bash
   docker scout cves nginx:latest
   trivy image nginx:latest
   ```

4. **Use specific image tags:**
   ```dockerfile
   FROM nginx:1.21.6
   # Not: FROM nginx:latest
   ```

5. **Limit capabilities:**
   ```bash
   docker run --cap-drop=ALL --cap-add=NET_BIND_SERVICE nginx
   ```

6. **Read-only filesystem:**
   ```bash
   docker run --read-only nginx
   ```

7. **Use secrets management:**
   - Docker secrets (Swarm)
   - External tools (Vault, AWS Secrets Manager)

8. **Network isolation:**
   ```bash
   docker network create isolated-network
   docker run --network isolated-network app
   ```

### 42. What is Docker Content Trust?

**Answer:**
Docker Content Trust (DCT) provides the ability to use digital signatures for data sent to and received from remote Docker registries. It:
- **Image signing**: Signs images with cryptographic keys
- **Verification**: Verifies image signatures before pulling
- **Tamper detection**: Detects if images are modified
- **Trust**: Ensures image integrity and publisher

**Enable Content Trust:**
```bash
export DOCKER_CONTENT_TRUST=1
docker pull nginx:latest
```

**Sign images:**
```bash
docker trust sign nginx:1.21
```

**Benefits:**
- Image integrity
- Publisher verification
- Prevents tampering
- Production security

### 43. What is the difference between docker run --privileged and --cap-add?

**Answer:**

| Aspect | --privileged | --cap-add |
|--------|--------------|-----------|
| **Security** | Less secure | More secure |
| **Capabilities** | All capabilities | Specific capabilities |
| **Use Case** | Development, testing | Production |
| **Best Practice** | Avoid | Recommended |

**--privileged:**
- Grants all Linux capabilities
- Disables security features
- Should be avoided in production
- Use only when absolutely necessary

**--cap-add:**
- Grants specific capabilities
- More secure
- Principle of least privilege
- Recommended approach

**Example:**
```bash
# Privileged (not recommended)
docker run --privileged app

# Specific capabilities (recommended)
docker run --cap-add=NET_ADMIN --cap-add=SYS_TIME app
```

### 44. How do you implement secrets management in Docker?

**Answer:**

**Docker Swarm Secrets:**
```bash
# Create secret
echo "my-secret" | docker secret create my_secret -

# Use in service
docker service create \
  --secret my_secret \
  --name my-service \
  nginx
```

**Docker Compose Secrets:**
```yaml
version: '3.8'
services:
  app:
    image: nginx
    secrets:
      - my_secret
secrets:
  my_secret:
    external: true
```

**External Secret Management:**
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Google Secret Manager

**Best Practices:**
- Never commit secrets to version control
- Use external secret management
- Rotate secrets regularly
- Limit secret access

### 45. What is Docker Security Scanning?

**Answer:**
Docker Security Scanning analyzes images for known vulnerabilities. It:
- **Vulnerability detection**: Finds known CVEs
- **Image analysis**: Scans image layers
- **Reports**: Provides detailed vulnerability reports
- **Remediation**: Suggests fixes

**Tools:**
- **Docker Scout**: Docker's native scanning
- **Trivy**: Open-source scanner
- **Clair**: CoreOS scanner
- **Snyk**: Commercial scanner

**Example:**
```bash
# Docker Scout
docker scout cves nginx:latest

# Trivy
trivy image nginx:latest

# Snyk
snyk test --docker nginx:latest
```

**Best Practices:**
- Scan images before deployment
- Integrate into CI/CD pipeline
- Regularly update base images
- Fix high/critical vulnerabilities

---

## Troubleshooting

### 46. How do you debug a container that won't start?

**Answer:**

**Steps:**

1. **Check container logs:**
   ```bash
   docker logs <container-id>
   docker logs --tail 100 <container-id>
   ```

2. **Inspect container:**
   ```bash
   docker inspect <container-id>
   ```

3. **Check exit code:**
   ```bash
   docker inspect <container-id> | grep ExitCode
   ```

4. **Run interactively:**
   ```bash
   docker run -it --entrypoint /bin/sh <image>
   ```

5. **Check Docker daemon logs:**
   ```bash
   sudo journalctl -u docker -f
   ```

6. **Check resources:**
   ```bash
   docker stats
   free -h
   df -h
   ```

**Common Issues:**
- Application errors (check logs)
- Resource constraints
- Port conflicts
- Volume mount issues
- Missing dependencies

### 47. How do you troubleshoot Docker networking issues?

**Answer:**

**Diagnosis:**

1. **Check network connectivity:**
   ```bash
   docker exec <container> ping <target>
   docker exec <container> curl <url>
   ```

2. **Inspect network:**
   ```bash
   docker network inspect <network-name>
   ```

3. **List networks:**
   ```bash
   docker network ls
   ```

4. **Check container network:**
   ```bash
   docker inspect <container> | grep -A 20 "NetworkSettings"
   ```

5. **Test DNS resolution:**
   ```bash
   docker exec <container> nslookup <hostname>
   ```

6. **Check iptables rules:**
   ```bash
   sudo iptables -L -n
   ```

**Common Issues:**
- Network not created
- Containers on different networks
- DNS resolution failures
- Port conflicts
- Firewall rules

### 48. How do you troubleshoot Docker volume issues?

**Answer:**

**Steps:**

1. **List volumes:**
   ```bash
   docker volume ls
   ```

2. **Inspect volume:**
   ```bash
   docker volume inspect <volume-name>
   ```

3. **Check volume mount:**
   ```bash
   docker inspect <container> | grep -A 10 "Mounts"
   ```

4. **Test volume access:**
   ```bash
   docker exec <container> ls -la /mount/point
   docker exec <container> touch /mount/point/test
   ```

5. **Check permissions:**
   ```bash
   docker exec <container> id
   docker exec <container> ls -la /mount/point
   ```

**Common Issues:**
- Volume not created
- Permission denied
- Volume not mounted
- Path doesn't exist
- Disk space issues

### 49. How do you monitor Docker containers?

**Answer:**

**Docker Stats:**
```bash
# Real-time stats
docker stats

# Stats for specific container
docker stats <container-id>

# No-stream mode
docker stats --no-stream
```

**Docker Events:**
```bash
# Monitor events
docker events

# Filter events
docker events --filter 'type=container'
```

**Third-party Tools:**
- **cAdvisor**: Container metrics
- **Prometheus**: Metrics collection
- **Grafana**: Visualization
- **Datadog**: Monitoring platform

**Logs:**
```bash
# Container logs
docker logs -f <container>

# All containers
docker-compose logs -f
```

### 50. How do you optimize Docker image build time?

**Answer:**

**Best Practices:**

1. **Use build cache:**
   ```dockerfile
   # Order matters - copy dependencies first
   COPY requirements.txt .
   RUN pip install -r requirements.txt
   COPY . .
   ```

2. **Multi-stage builds:**
   ```dockerfile
   FROM node:18 AS builder
   WORKDIR /app
   COPY . .
   RUN npm build

   FROM nginx:alpine
   COPY --from=builder /app/dist /usr/share/nginx/html
   ```

3. **Combine RUN commands:**
   ```dockerfile
   # Bad
   RUN apt-get update
   RUN apt-get install -y package1
   RUN apt-get install -y package2

   # Good
   RUN apt-get update && \
       apt-get install -y package1 package2 && \
       rm -rf /var/lib/apt/lists/*
   ```

4. **Use .dockerignore:**
   ```
   node_modules
   .git
   *.log
   ```

5. **Parallel builds:**
   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 .
   ```

6. **Use BuildKit:**
   ```bash
   DOCKER_BUILDKIT=1 docker build .
   ```

---

## Advanced Topics

### 51. What is Docker BuildKit?

**Answer:**
BuildKit is the next-generation build engine for Docker. It provides:
- **Faster builds**: Parallel execution, better caching
- **Advanced features**: Secrets, SSH, cache mounts
- **Multi-platform**: Build for multiple architectures
- **Improved UX**: Better build output

**Enable BuildKit:**
```bash
export DOCKER_BUILDKIT=1
docker build .
```

**Features:**
- Secret mounting
- SSH agent forwarding
- Cache mounts
- Parallel builds
- Better error messages

**Example:**
```dockerfile
# syntax=docker/dockerfile:1.4
FROM alpine
RUN --mount=type=secret,id=mysecret \
    cat /run/secrets/mysecret
```

### 52. What is the difference between docker save and docker export?

**Answer:**

| Aspect | docker save | docker export |
|--------|------------|---------------|
| **Target** | Images | Containers |
| **Format** | tar archive with layers | tar archive (flattened) |
| **Layers** | Preserves layers | Flattens layers |
| **Metadata** | Preserves metadata | Loses metadata |
| **Use Case** | Backup images | Backup container filesystem |

**docker save:**
- Saves images
- Preserves layers and metadata
- Can load back as image
- Use for image backup

**docker export:**
- Exports container filesystem
- Flattens layers
- Loses metadata
- Use for filesystem backup

**Example:**
```bash
# Save image
docker save nginx:latest > nginx.tar
docker load < nginx.tar

# Export container
docker export container-id > container.tar
docker import container.tar new-image:tag
```

### 53. What is Docker Swarm and how does it differ from Kubernetes?

**Answer:**

| Aspect | Docker Swarm | Kubernetes |
|--------|--------------|------------|
| **Complexity** | Simple | Complex |
| **Setup** | Easy | More complex |
| **Features** | Basic orchestration | Advanced features |
| **Ecosystem** | Smaller | Larger |
| **Use Case** | Small-medium clusters | Large-scale production |

**Docker Swarm:**
- Native Docker orchestration
- Simple setup and management
- Good for Docker-focused teams
- Limited features compared to Kubernetes

**Kubernetes:**
- Industry standard
- Rich feature set
- Large ecosystem
- Steeper learning curve

**When to use Swarm:**
- Small to medium deployments
- Docker-focused teams
- Simple orchestration needs
- Quick setup required

### 54. What is the difference between docker-compose and docker stack?

**Answer:**

| Aspect | docker-compose | docker stack |
|--------|----------------|--------------|
| **Scope** | Single host | Swarm cluster |
| **Orchestration** | Compose | Swarm |
| **Networks** | Compose networks | Overlay networks |
| **Secrets** | File-based | Swarm secrets |
| **Use Case** | Development | Production (Swarm) |

**docker-compose:**
- Single host
- Development and testing
- File-based configuration
- Local development

**docker stack:**
- Swarm cluster
- Production deployments
- Swarm-native features
- Multi-host deployments

**Example:**
```bash
# Compose
docker-compose up

# Stack
docker stack deploy -c docker-compose.yml my-stack
```

### 55. What is a Docker registry and how do you set up a private registry?

**Answer:**
A Docker registry is a storage and distribution system for Docker images. Types:
- **Public**: Docker Hub (public images)
- **Private**: Self-hosted or cloud-based

**Set up private registry:**
```bash
# Run registry
docker run -d -p 5000:5000 --name registry registry:2

# Tag and push
docker tag my-image localhost:5000/my-image
docker push localhost:5000/my-image

# Pull
docker pull localhost:5000/my-image
```

**Production registry:**
- **Harbor**: Enterprise registry
- **GitLab Container Registry**: GitLab integrated
- **AWS ECR**: AWS managed
- **Azure ACR**: Azure managed
- **GCR**: Google managed

### 56. What is the difference between docker attach and docker exec?

**Answer:**

| Aspect | docker attach | docker exec |
|--------|---------------|-------------|
| **Process** | Main process (PID 1) | New process |
| **Sessions** | Single | Multiple |
| **Exit** | May stop container | Doesn't affect container |
| **Use Case** | Interactive apps | Debugging, commands |

**docker attach:**
- Attaches to main process
- Single session
- Exiting may stop container
- Use for interactive applications

**docker exec:**
- Creates new process
- Multiple sessions possible
- Exiting doesn't affect container
- Use for debugging, running commands

### 57. How do you implement health checks in Docker?

**Answer:**

**Dockerfile HEALTHCHECK:**
```dockerfile
FROM nginx
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
```

**docker run:**
```bash
docker run --health-cmd="curl -f http://localhost || exit 1" \
  --health-interval=30s \
  --health-timeout=3s \
  --health-retries=3 \
  nginx
```

**docker-compose:**
```yaml
services:
  web:
    image: nginx
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

**Check health:**
```bash
docker ps  # Shows health status
docker inspect <container> | grep -A 10 Health
```

### 58. What is the difference between docker stop and docker pause?

**Answer:**

| Aspect | docker stop | docker pause |
|--------|-------------|--------------|
| **State** | Stops container | Pauses container |
| **Process** | Terminates | Suspends |
| **Resources** | Releases | Keeps allocated |
| **Resume** | docker start | docker unpause |

**docker stop:**
- Terminates container process
- Releases resources
- Container state: Exited
- Resume: `docker start`

**docker pause:**
- Suspends container process
- Keeps resources allocated
- Container state: Paused
- Resume: `docker unpause`

**Use Cases:**
- **stop**: Normal shutdown, resource cleanup
- **pause**: Temporary suspension, debugging

### 59. How do you manage Docker logs?

**Answer:**

**View logs:**
```bash
# Container logs
docker logs <container>

# Follow logs
docker logs -f <container>

# Last N lines
docker logs --tail 100 <container>

# Since timestamp
docker logs --since 2024-01-01T00:00:00 <container>
```

**Log drivers:**
```bash
# JSON file (default)
docker run --log-driver json-file nginx

# Syslog
docker run --log-driver syslog nginx

# Journald
docker run --log-driver journald nginx

# Custom
docker run --log-driver fluentd nginx
```

**Log rotation:**
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

**Centralized logging:**
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Fluentd
- Splunk
- Datadog

### 60. What is the difference between docker build and docker buildx?

**Answer:**

| Aspect | docker build | docker buildx |
|--------|--------------|---------------|
| **Engine** | Legacy builder | BuildKit |
| **Multi-platform** | Limited | Full support |
| **Features** | Basic | Advanced |
| **Performance** | Slower | Faster |

**docker build:**
- Legacy builder
- Single platform
- Basic features
- Still widely used

**docker buildx:**
- BuildKit-based
- Multi-platform builds
- Advanced features (cache, secrets)
- Recommended for new projects

**Example:**
```bash
# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t myapp .

# Create builder
docker buildx create --name mybuilder --use
```

### 61. How do you implement CI/CD with Docker?

**Answer:**

**CI/CD Pipeline:**

1. **Build stage:**
   ```yaml
   - name: Build Docker image
     run: |
       docker build -t myapp:${{ github.sha }} .
   ```

2. **Test stage:**
   ```yaml
   - name: Run tests
     run: |
       docker run myapp:${{ github.sha }} npm test
   ```

3. **Security scan:**
   ```yaml
   - name: Scan image
     run: |
       trivy image myapp:${{ github.sha }}
   ```

4. **Push to registry:**
   ```yaml
   - name: Push image
     run: |
       docker push myapp:${{ github.sha }}
   ```

5. **Deploy:**
   ```yaml
   - name: Deploy
     run: |
       kubectl set image deployment/myapp app=myapp:${{ github.sha }}
   ```

**Best Practices:**
- Use multi-stage builds
- Cache layers
- Scan for vulnerabilities
- Tag with commit SHA
- Use semantic versioning

### 62. What is the difference between docker-compose and docker run?

**Answer:**

| Aspect | docker-compose | docker run |
|--------|----------------|------------|
| **Scope** | Multi-container | Single container |
| **Configuration** | YAML file | Command line |
| **Networking** | Automatic | Manual |
| **Use Case** | Complex apps | Simple apps |

**docker-compose:**
- Manages multiple containers
- YAML-based configuration
- Automatic networking
- Volume management
- Environment variables
- Use for complex applications

**docker run:**
- Single container
- Command-line based
- Manual setup
- Quick testing
- Use for simple applications

### 63. How do you implement blue-green deployment with Docker?

**Answer:**

**Strategy:**
1. Run two identical production environments (blue and green)
2. Route traffic to one environment
3. Deploy new version to idle environment
4. Switch traffic to new environment
5. Keep old environment for rollback

**Implementation:**
```bash
# Blue environment (current)
docker-compose -f docker-compose.blue.yml up -d

# Deploy to green
docker-compose -f docker-compose.green.yml up -d

# Switch traffic (update load balancer)
# Test green environment

# If successful, stop blue
docker-compose -f docker-compose.blue.yml down

# If failed, switch back to blue
```

**With orchestration:**
- Kubernetes: Use multiple deployments
- Docker Swarm: Use service updates
- Load balancer: Route traffic

### 64. What is the difference between docker commit and docker build?

**Answer:**

| Aspect | docker commit | docker build |
|--------|---------------|--------------|
| **Source** | Running container | Dockerfile |
| **Reproducibility** | Not reproducible | Reproducible |
| **Best Practice** | Not recommended | Recommended |
| **Version Control** | No | Yes |
| **Use Case** | Quick saves | Production |

**docker commit:**
- Creates image from container
- Not reproducible
- Loses build history
- Quick but not recommended

**docker build:**
- Builds from Dockerfile
- Reproducible
- Version controlled
- Best practice

### 65. How do you optimize Docker for production?

**Answer:**

**Best Practices:**

1. **Use specific image tags:**
   ```dockerfile
   FROM nginx:1.21.6
   # Not: FROM nginx:latest
   ```

2. **Multi-stage builds:**
   ```dockerfile
   FROM node:18 AS builder
   WORKDIR /app
   COPY . .
   RUN npm build

   FROM nginx:alpine
   COPY --from=builder /app/dist /usr/share/nginx/html
   ```

3. **Resource limits:**
   ```bash
   docker run --memory="512m" --cpus="1.0" app
   ```

4. **Health checks:**
   ```dockerfile
   HEALTHCHECK --interval=30s CMD curl -f http://localhost || exit 1
   ```

5. **Read-only root filesystem:**
   ```bash
   docker run --read-only --tmpfs /tmp app
   ```

6. **Non-root user:**
   ```dockerfile
   RUN adduser -D appuser
   USER appuser
   ```

7. **Security scanning:**
   ```bash
   trivy image myapp:latest
   ```

8. **Log management:**
   ```json
   {
     "log-driver": "json-file",
     "log-opts": {
       "max-size": "10m",
       "max-file": "3"
     }
   }
   ```

---

## Advanced Docker Concepts

### 66. What is the difference between docker build and docker buildx build?

**Answer:**

| Aspect | docker build | docker buildx build |
|--------|--------------|---------------------|
| **Engine** | Legacy builder | BuildKit |
| **Multi-platform** | Limited | Full support |
| **Features** | Basic | Advanced (secrets, SSH) |
| **Performance** | Slower | Faster |

**docker build:**
- Uses legacy builder
- Single platform builds
- Basic features
- Still widely used

**docker buildx build:**
- Uses BuildKit
- Multi-platform builds
- Advanced features (cache mounts, secrets)
- Recommended for new projects

**Example:**
```bash
# Legacy build
docker build -t myapp .

# BuildKit build
DOCKER_BUILDKIT=1 docker build -t myapp .

# Buildx multi-platform
docker buildx build --platform linux/amd64,linux/arm64 -t myapp .
```

### 67. How do you implement multi-stage builds effectively?

**Answer:**

**Best Practices:**
1. **Separate build and runtime:**
   ```dockerfile
   FROM golang:1.21 AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o app

   FROM alpine:latest
   WORKDIR /app
   COPY --from=builder /app/app .
   CMD ["./app"]
   ```

2. **Use minimal base images:**
   - Alpine for runtime
   - Distroless for security
   - Scratch for minimal size

3. **Copy only what's needed:**
   - Don't copy build tools
   - Don't copy source code
   - Only copy artifacts

4. **Name stages:**
   ```dockerfile
   FROM node:18 AS dependencies
   FROM node:18 AS build
   FROM nginx:alpine AS runtime
   ```

### 68. What is the difference between COPY and ADD in Dockerfile?

**Answer:**

| Aspect | COPY | ADD |
|--------|------|-----|
| **Source** | Local files only | Local files + URLs |
| **URLs** | Not supported | Downloads from URLs |
| **Tar extraction** | No | Yes (automatic) |
| **Best Practice** | Preferred | Use for URLs/tar |

**COPY:**
- Simple file copy
- More explicit
- Recommended for local files
- Better caching

**ADD:**
- Can download from URLs
- Can extract tar files
- Less explicit
- Use only when needed

**Example:**
```dockerfile
# COPY (recommended)
COPY app.py /app/
COPY requirements.txt /app/

# ADD (for URLs or tar)
ADD https://example.com/file.tar.gz /tmp/
# Automatically extracts
```

### 69. How do you optimize Docker image layers?

**Answer:**

**Best Practices:**
1. **Combine RUN commands:**
   ```dockerfile
   # Bad - multiple layers
   RUN apt-get update
   RUN apt-get install -y package1
   RUN apt-get install -y package2

   # Good - single layer
   RUN apt-get update && \
       apt-get install -y package1 package2 && \
       rm -rf /var/lib/apt/lists/*
   ```

2. **Order instructions:**
   - Copy dependency files first
   - Install dependencies
   - Copy application code last

3. **Use .dockerignore:**
   - Exclude unnecessary files
   - Reduces build context

4. **Minimize layers:**
   - Combine related operations
   - Use multi-stage builds

### 70. What is Docker BuildKit and its advanced features?

**Answer:**
BuildKit is the next-generation build engine. Advanced features:

**Secret Mounting:**
```dockerfile
# syntax=docker/dockerfile:1.4
FROM alpine
RUN --mount=type=secret,id=mysecret \
    cat /run/secrets/mysecret
```

**SSH Mounting:**
```dockerfile
# syntax=docker/dockerfile:1.4
FROM alpine
RUN --mount=type=ssh \
    ssh-add -l
```

**Cache Mounts:**
```dockerfile
# syntax=docker/dockerfile:1.4
FROM node:18
RUN --mount=type=cache,target=/root/.npm \
    npm install
```

**Benefits:**
- Faster builds
- Better caching
- Parallel execution
- Advanced features

---

## Docker Networking Advanced

### 71. How do you create and manage custom Docker networks?

**Answer:**

**Create Networks:**
```bash
# Bridge network
docker network create my-network

# Overlay network (Swarm)
docker network create --driver overlay my-overlay

# Macvlan network
docker network create --driver macvlan \
  --subnet=192.168.1.0/24 \
  --gateway=192.168.1.1 \
  -o parent=eth0 my-macvlan
```

**Network Types:**
- **bridge**: Default, single host
- **overlay**: Multi-host (Swarm)
- **host**: Uses host network
- **macvlan**: Assigns MAC addresses
- **none**: No networking

**Manage Networks:**
```bash
# List networks
docker network ls

# Inspect network
docker network inspect my-network

# Connect container
docker network connect my-network container

# Disconnect container
docker network disconnect my-network container

# Remove network
docker network rm my-network
```

### 72. What is the difference between bridge and overlay networks?

**Answer:**

| Aspect | Bridge | Overlay |
|--------|--------|---------|
| **Scope** | Single host | Multi-host |
| **Use Case** | Local development | Swarm cluster |
| **Encryption** | No | Optional |
| **Performance** | Better | Slightly slower |

**Bridge:**
- Single Docker host
- Default network driver
- Good for local development
- Better performance

**Overlay:**
- Multiple Docker hosts
- Swarm cluster networking
- Can encrypt traffic
- Service discovery across hosts

**Example:**
```bash
# Bridge (single host)
docker network create --driver bridge my-bridge

# Overlay (Swarm)
docker network create --driver overlay my-overlay
```

### 73. How do you implement service discovery in Docker?

**Answer:**

**Docker DNS:**
- Containers can resolve each other by name
- Works within same network
- Automatic service discovery

**Example:**
```bash
# Create network
docker network create my-network

# Run containers on same network
docker run -d --name web --network my-network nginx
docker run -d --name app --network my-network myapp

# App can resolve 'web' by name
# curl http://web
```

**Docker Compose:**
```yaml
services:
  web:
    image: nginx
  app:
    image: myapp
    # Can access 'web' by service name
```

**External DNS:**
- Use external DNS servers
- Configure in daemon.json
- Use --dns flag

### 74. How do you troubleshoot Docker network connectivity?

**Answer:**

**Diagnosis:**
1. **Check network:**
   ```bash
   docker network inspect <network-name>
   ```

2. **Test connectivity:**
   ```bash
   docker exec <container> ping <target>
   docker exec <container> curl <url>
   ```

3. **Check DNS:**
   ```bash
   docker exec <container> nslookup <hostname>
   ```

4. **Inspect container network:**
   ```bash
   docker inspect <container> | grep -A 20 "NetworkSettings"
   ```

5. **Check iptables:**
   ```bash
   sudo iptables -L -n
   ```

**Common Issues:**
- Containers on different networks
- DNS resolution failures
- Firewall rules
- Network driver issues

### 75. What is Docker's default bridge network?

**Answer:**
The default bridge network is created automatically when Docker is installed. It:
- **Name**: `bridge`
- **Driver**: bridge
- **Isolation**: Isolates containers from host
- **DNS**: No automatic DNS resolution

**Characteristics:**
- Containers can communicate by IP
- No automatic DNS resolution
- Port mapping required for external access
- Isolated from host network

**Custom Bridge:**
- Better than default bridge
- Automatic DNS resolution
- Better isolation
- Recommended for production

**Example:**
```bash
# Default bridge (not recommended)
docker run nginx

# Custom bridge (recommended)
docker network create my-network
docker run --network my-network nginx
```

---

## Docker Storage Advanced

### 76. How do you implement volume backup strategies?

**Answer:**

**Backup Methods:**

1. **Using temporary container:**
   ```bash
   docker run --rm \
     -v my-volume:/data:ro \
     -v $(pwd):/backup \
     ubuntu tar czf /backup/backup.tar.gz /data
   ```

2. **Using docker cp:**
   ```bash
   docker run -d --name temp -v my-volume:/data ubuntu sleep 3600
   docker cp temp:/data ./backup
   docker rm -f temp
   ```

3. **Using bind mount:**
   ```bash
   docker run --rm \
     -v my-volume:/source \
     -v /backup:/backup \
     ubuntu tar czf /backup/backup.tar.gz /source
   ```

**Automated Backup:**
- Cron jobs
- Backup scripts
- Cloud storage integration
- Regular testing

### 77. What is the difference between named volumes and anonymous volumes?

**Answer:**

| Aspect | Named Volume | Anonymous Volume |
|--------|--------------|------------------|
| **Name** | Explicit name | Random name |
| **Management** | Easy to manage | Hard to manage |
| **Backup** | Easier | Harder |
| **Use Case** | Production | Temporary |

**Named Volume:**
- Explicit name: `my-volume`
- Easy to reference
- Better for production
- Easier to backup

**Anonymous Volume:**
- Random name: `a1b2c3d4...`
- Created automatically
- Hard to reference
- Use for temporary data

**Example:**
```bash
# Named volume
docker volume create my-volume
docker run -v my-volume:/data nginx

# Anonymous volume
docker run -v /data nginx
```

### 78. How do you manage Docker disk space?

**Answer:**

**Cleanup Commands:**
```bash
# Remove unused containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove unused networks
docker network prune

# Remove everything unused
docker system prune

# Remove everything including volumes
docker system prune -a --volumes
```

**Monitor Disk Usage:**
```bash
# System disk usage
docker system df

# Detailed usage
docker system df -v
```

**Prevention:**
- Regular cleanup
- Use multi-stage builds
- Remove unused images
- Limit log size
- Use .dockerignore

### 79. What is Docker's storage driver and how to choose one?

**Answer:**

**Storage Drivers:**
- **overlay2**: Recommended for most Linux
- **devicemapper**: Older systems
- **btrfs**: For btrfs filesystems
- **zfs**: For ZFS filesystems
- **aufs**: Legacy

**Check Current Driver:**
```bash
docker info | grep "Storage Driver"
```

**Choose Driver:**
- **overlay2**: Default, best performance
- **devicemapper**: If overlay2 not available
- **btrfs/zfs**: If using those filesystems

**Configure Driver:**
```json
{
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ]
}
```

### 80. How do you implement data persistence in Docker?

**Answer:**

**Methods:**

1. **Volumes (Recommended):**
   ```bash
   docker volume create my-data
   docker run -v my-data:/data nginx
   ```

2. **Bind Mounts:**
   ```bash
   docker run -v /host/path:/container/path nginx
   ```

3. **tmpfs Mounts:**
   ```bash
   docker run --tmpfs /tmp nginx
   ```

**Best Practices:**
- Use volumes for production
- Use bind mounts for development
- Use tmpfs for temporary data
- Backup volumes regularly
- Use appropriate storage drivers

---

## Docker Compose Advanced

### 81. How do you use environment variables in Docker Compose?

**Answer:**

**Methods:**

1. **Environment file (.env):**
   ```bash
   # .env
   DB_PASSWORD=secret
   DB_HOST=localhost
   ```

2. **In docker-compose.yml:**
   ```yaml
   services:
     app:
       image: nginx
       environment:
         - DB_HOST=${DB_HOST}
         - DB_PASSWORD=${DB_PASSWORD}
       env_file:
         - .env
   ```

3. **Command line:**
   ```bash
   DB_HOST=localhost docker-compose up
   ```

**Variable Substitution:**
```yaml
services:
  app:
    image: ${IMAGE_NAME:-nginx}:${IMAGE_TAG:-latest}
    ports:
      - "${PORT:-8080}:80"
```

### 82. How do you implement health checks in Docker Compose?

**Answer:**

**Health Check:**
```yaml
services:
  web:
    image: nginx
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

**Dependencies:**
```yaml
services:
  app:
    image: myapp
    depends_on:
      db:
        condition: service_healthy
  db:
    image: postgres
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
```

**Check Health:**
```bash
docker-compose ps
docker inspect <container> | grep -A 10 Health
```

### 83. How do you scale services in Docker Compose?

**Answer:**

**Scale Services:**
```bash
# Scale specific service
docker-compose up -d --scale web=3

# Scale multiple services
docker-compose up -d --scale web=3 --scale worker=2
```

**Limitations:**
- Service must not have single port mapping
- No built-in load balancing
- Use external load balancer

**With Load Balancer:**
```yaml
services:
  nginx:
    image: nginx
    ports:
      - "80:80"
    depends_on:
      - web
  web:
    image: myapp
    # No port mapping
```

### 84. What is the difference between docker-compose and docker stack?

**Answer:**

| Aspect | docker-compose | docker stack |
|--------|----------------|--------------|
| **Scope** | Single host | Swarm cluster |
| **Orchestration** | Compose | Swarm |
| **Networks** | Compose networks | Overlay networks |
| **Secrets** | File-based | Swarm secrets |
| **Deploy** | docker-compose up | docker stack deploy |

**docker-compose:**
- Single host
- Development and testing
- File-based secrets
- Local development

**docker stack:**
- Swarm cluster
- Production deployments
- Swarm secrets
- Multi-host

**Example:**
```bash
# Compose
docker-compose up

# Stack
docker stack deploy -c docker-compose.yml my-stack
```

### 85. How do you implement blue-green deployment with Docker Compose?

**Answer:**

**Strategy:**
1. Run blue environment (current)
2. Deploy green environment (new)
3. Switch traffic
4. Test green
5. Stop blue if successful

**Implementation:**
```bash
# Blue (current)
docker-compose -f docker-compose.blue.yml up -d

# Green (new)
docker-compose -f docker-compose.green.yml up -d

# Switch traffic (update load balancer)
# Test green environment

# If successful
docker-compose -f docker-compose.blue.yml down

# If failed, switch back
```

**With Load Balancer:**
- Use nginx/HAProxy
- Update upstream configuration
- Health checks
- Gradual traffic shift

---

## Docker Security Advanced

### 86. How do you implement image scanning in CI/CD?

**Answer:**

**Tools:**
- **Trivy**: Open-source scanner
- **Docker Scout**: Docker's scanner
- **Snyk**: Commercial scanner
- **Clair**: CoreOS scanner

**CI/CD Integration:**
```yaml
# GitHub Actions example
- name: Build image
  run: docker build -t myapp:${{ github.sha }} .

- name: Scan image
  run: |
    trivy image myapp:${{ github.sha }}
    # Fail on high/critical vulnerabilities

- name: Push image
  if: success()
  run: docker push myapp:${{ github.sha }}
```

**Best Practices:**
- Scan before push
- Fail on high/critical CVEs
- Regular base image updates
- Document vulnerabilities

### 87. How do you implement secrets management in Docker?

**Answer:**

**Docker Swarm Secrets:**
```bash
# Create secret
echo "my-secret" | docker secret create my_secret -

# Use in service
docker service create \
  --secret my_secret \
  --name my-service \
  nginx
```

**Docker Compose Secrets:**
```yaml
services:
  app:
    image: nginx
    secrets:
      - my_secret
secrets:
  my_secret:
    external: true
```

**External Tools:**
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Google Secret Manager

**Best Practices:**
- Never commit secrets
- Use external management
- Rotate regularly
- Limit access

### 88. What is Docker Content Trust and how to use it?

**Answer:**
Docker Content Trust provides image signing and verification. It:
- **Signs images**: Cryptographic signatures
- **Verifies images**: Checks signatures before pull
- **Prevents tampering**: Detects modifications
- **Trust**: Ensures image integrity

**Enable:**
```bash
export DOCKER_CONTENT_TRUST=1
docker pull nginx:latest
```

**Sign Images:**
```bash
docker trust sign nginx:1.21
```

**Benefits:**
- Image integrity
- Publisher verification
- Tamper detection
- Production security

### 89. How do you secure Docker daemon?

**Answer:**

**Best Practices:**
1. **Use TLS:**
   ```bash
   # Generate certificates
   # Configure daemon with TLS
   ```

2. **Limit access:**
   - Use firewall
   - Restrict network access
   - Use VPN

3. **Regular updates:**
   - Keep Docker updated
   - Security patches
   - Monitor vulnerabilities

4. **Audit logging:**
   - Enable audit logs
   - Monitor access
   - Regular reviews

5. **User permissions:**
   - Use docker group carefully
   - Principle of least privilege
   - Regular access reviews

### 90. What are Docker security best practices?

**Answer:**

**Image Security:**
- Use minimal base images
- Scan for vulnerabilities
- Use specific tags (not latest)
- Multi-stage builds
- Remove unnecessary packages

**Container Security:**
- Run as non-root
- Limit capabilities
- Read-only filesystem
- Resource limits
- Network isolation

**Runtime Security:**
- Use secrets management
- Enable content trust
- Regular updates
- Monitor containers
- Audit logs

**Infrastructure Security:**
- Secure Docker daemon
- Use TLS
- Network policies
- Firewall rules
- Access control

---

## Production & Operations

### 91. How do you implement logging in Docker production?

**Answer:**

**Log Drivers:**
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

**Centralized Logging:**
- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **EFK Stack**: Elasticsearch, Fluentd, Kibana
- **Loki + Grafana**: Lightweight
- **Splunk**: Enterprise

**Example:**
```bash
# JSON file (default)
docker run --log-driver json-file nginx

# Syslog
docker run --log-driver syslog nginx

# Fluentd
docker run --log-driver fluentd nginx
```

**Best Practices:**
- Structured logging (JSON)
- Log rotation
- Centralized collection
- Retention policies

### 92. How do you monitor Docker in production?

**Answer:**

**Docker Stats:**
```bash
docker stats
docker stats --no-stream
```

**Monitoring Tools:**
- **cAdvisor**: Container metrics
- **Prometheus**: Metrics collection
- **Grafana**: Visualization
- **Datadog**: Commercial platform

**Key Metrics:**
- CPU usage
- Memory usage
- Disk I/O
- Network I/O
- Container count

**Example:**
```yaml
# docker-compose.yml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

### 93. How do you implement backup and disaster recovery?

**Answer:**

**Backup Strategy:**
1. **Image backup:**
   ```bash
   docker save myapp:latest > myapp.tar
   ```

2. **Volume backup:**
   ```bash
   docker run --rm \
     -v my-volume:/data:ro \
     -v $(pwd):/backup \
     ubuntu tar czf /backup/backup.tar.gz /data
   ```

3. **Configuration backup:**
   - Backup docker-compose.yml
   - Backup environment files
   - Version control

**Disaster Recovery:**
- Regular backups
- Test restore procedures
- Document recovery process
- Off-site backups
- Recovery time objectives

### 94. How do you optimize Docker for production?

**Answer:**

**Image Optimization:**
- Multi-stage builds
- Minimal base images
- Remove unnecessary packages
- Optimize layers
- Use .dockerignore

**Runtime Optimization:**
- Resource limits
- Health checks
- Log management
- Network optimization
- Storage optimization

**Infrastructure:**
- Right-size hosts
- Use orchestration
- Load balancing
- High availability
- Monitoring

**Example:**
```dockerfile
FROM node:18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY . .
CMD ["node", "app.js"]
```

### 95. How do you implement CI/CD with Docker?

**Answer:**

**Pipeline Stages:**
1. **Build:**
   ```yaml
   - name: Build
     run: docker build -t myapp:${{ github.sha }} .
   ```

2. **Test:**
   ```yaml
   - name: Test
     run: docker run myapp:${{ github.sha }} npm test
   ```

3. **Scan:**
   ```yaml
   - name: Scan
     run: trivy image myapp:${{ github.sha }}
   ```

4. **Push:**
   ```yaml
   - name: Push
     run: docker push myapp:${{ github.sha }}
   ```

5. **Deploy:**
   ```yaml
   - name: Deploy
     run: kubectl set image deployment/myapp app=myapp:${{ github.sha }}
   ```

**Best Practices:**
- Multi-stage builds
- Cache layers
- Parallel builds
- Security scanning
- Tagging strategy

---

## Troubleshooting & Debugging

### 96. How do you debug Docker build failures?

**Answer:**

**Steps:**
1. **Verbose output:**
   ```bash
   docker build --progress=plain --no-cache -t myapp .
   ```

2. **Check Dockerfile:**
   - Syntax errors
   - Missing files
   - Incorrect paths

3. **Test layers:**
   ```bash
   # Build up to specific layer
   docker build --target builder -t myapp:builder .
   ```

4. **Check build context:**
   ```bash
   docker system df
   ```

5. **Interactive debugging:**
   ```bash
   docker run -it --entrypoint /bin/sh myapp:builder
   ```

**Common Issues:**
- Dockerfile syntax
- Missing files
- Network issues
- Disk space
- Base image issues

### 97. How do you troubleshoot container performance issues?

**Answer:**

**Diagnosis:**
1. **Check resources:**
   ```bash
   docker stats
   docker stats <container>
   ```

2. **Check logs:**
   ```bash
   docker logs <container>
   docker logs --tail 100 <container>
   ```

3. **Inspect container:**
   ```bash
   docker inspect <container>
   ```

4. **Check host resources:**
   ```bash
   free -h
   df -h
   top
   ```

5. **Profile application:**
   - Application profiling tools
   - Performance monitoring
   - Resource analysis

**Common Issues:**
- Resource limits too low
- Memory leaks
- CPU throttling
- I/O bottlenecks
- Network issues

### 98. How do you troubleshoot Docker daemon issues?

**Answer:**

**Diagnosis:**
1. **Check daemon status:**
   ```bash
   sudo systemctl status docker
   ```

2. **Check daemon logs:**
   ```bash
   sudo journalctl -u docker -f
   ```

3. **Check daemon info:**
   ```bash
   docker info
   ```

4. **Test connectivity:**
   ```bash
   docker ps
   docker version
   ```

5. **Check configuration:**
   ```bash
   cat /etc/docker/daemon.json
   ```

**Common Issues:**
- Daemon not running
- Configuration errors
- Disk space
- Permission issues
- Network issues

### 99. How do you debug Docker networking issues?

**Answer:**

**Steps:**
1. **Check network:**
   ```bash
   docker network ls
   docker network inspect <network>
   ```

2. **Test connectivity:**
   ```bash
   docker exec <container> ping <target>
   docker exec <container> curl <url>
   ```

3. **Check DNS:**
   ```bash
   docker exec <container> nslookup <hostname>
   ```

4. **Inspect container:**
   ```bash
   docker inspect <container> | grep -A 20 "NetworkSettings"
   ```

5. **Check iptables:**
   ```bash
   sudo iptables -L -n
   ```

**Common Issues:**
- Containers on different networks
- DNS failures
- Firewall rules
- Network driver issues
- Port conflicts

### 100. How do you implement Docker in a multi-host environment?

**Answer:**

**Options:**

1. **Docker Swarm:**
   ```bash
   # Initialize swarm
   docker swarm init

   # Join workers
   docker swarm join --token <token> <manager-ip>:2377

   # Deploy services
   docker stack deploy -c docker-compose.yml my-stack
   ```

2. **Kubernetes:**
   - Use Kubernetes for orchestration
   - Deploy Docker containers
   - Better for large-scale

3. **Docker Compose with remote hosts:**
   - Use DOCKER_HOST environment variable
   - SSH tunneling
   - Not recommended for production

**Best Practices:**
- Use orchestration (Swarm/Kubernetes)
- Overlay networks
- Service discovery
- Load balancing
- High availability

### 101. What is the difference between docker stop and docker kill?

**Answer:**

| Aspect | docker stop | docker kill |
|--------|-------------|-------------|
| **Signal** | SIGTERM | SIGKILL |
| **Grace Period** | 10s default | Immediate |
| **Cleanup** | Allows cleanup | No cleanup |
| **Use Case** | Normal shutdown | Force stop |

**docker stop:**
- Sends SIGTERM
- Waits for graceful shutdown
- Then SIGKILL if needed
- Allows cleanup

**docker kill:**
- Sends SIGKILL immediately
- Forceful termination
- No cleanup
- Use when unresponsive

### 102. How do you implement Docker health checks?

**Answer:**

**Dockerfile:**
```dockerfile
FROM nginx
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
```

**docker run:**
```bash
docker run --health-cmd="curl -f http://localhost || exit 1" \
  --health-interval=30s \
  --health-timeout=3s \
  --health-retries=3 \
  nginx
```

**docker-compose:**
```yaml
services:
  web:
    image: nginx
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

**Check Health:**
```bash
docker ps  # Shows health status
docker inspect <container> | grep -A 10 Health
```

### 103. How do you manage Docker images in production?

**Answer:**

**Best Practices:**
1. **Image registry:**
   - Use private registry
   - Tag appropriately
   - Version control

2. **Image scanning:**
   - Scan before deployment
   - Regular updates
   - Vulnerability management

3. **Image cleanup:**
   - Remove unused images
   - Regular pruning
   - Archive old images

4. **Image policies:**
   - Approved base images
   - Security requirements
   - Update procedures

**Example:**
```bash
# Tag for registry
docker tag myapp:latest registry.example.com/myapp:v1.0.0

# Push to registry
docker push registry.example.com/myapp:v1.0.0

# Pull from registry
docker pull registry.example.com/myapp:v1.0.0
```

### 104. How do you implement Docker resource limits?

**Answer:**

**CPU Limits:**
```bash
# Limit to 1 CPU
docker run --cpus="1.0" nginx

# Limit to 50% of 1 CPU
docker run --cpus="0.5" nginx
```

**Memory Limits:**
```bash
# Limit memory
docker run --memory="512m" nginx

# Memory + swap
docker run --memory="512m" --memory-swap="1g" nginx
```

**docker-compose:**
```yaml
services:
  app:
    image: nginx
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

**Best Practices:**
- Set appropriate limits
- Monitor usage
- Adjust based on needs
- Prevent resource exhaustion

### 105. What are Docker best practices for production?

**Answer:**

**Image Best Practices:**
- Use specific tags
- Multi-stage builds
- Minimal base images
- Scan for vulnerabilities
- Remove unnecessary packages

**Container Best Practices:**
- Run as non-root
- Limit capabilities
- Resource limits
- Health checks
- Read-only filesystem

**Security Best Practices:**
- Secrets management
- Content trust
- Network isolation
- Regular updates
- Security scanning

**Operations Best Practices:**
- Logging strategy
- Monitoring
- Backup procedures
- Disaster recovery
- Documentation

**Example Production Dockerfile:**
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY --chown=nodejs:nodejs . .
USER nodejs
HEALTHCHECK --interval=30s CMD node healthcheck.js
CMD ["node", "app.js"]
```

---

## Real-World Scenarios

### 106. How do you implement blue-green deployment with Docker?

**Answer:**

**Strategy:**
1. Run blue environment (current production)
2. Deploy green environment (new version)
3. Switch traffic to green
4. Test green environment
5. Keep blue for rollback or stop if successful

**Implementation:**
```bash
# Blue environment
docker-compose -f docker-compose.blue.yml up -d

# Deploy green
docker-compose -f docker-compose.green.yml up -d

# Switch traffic (update load balancer config)
# Test green

# If successful, stop blue
docker-compose -f docker-compose.blue.yml down

# If failed, switch back to blue
```

**With Load Balancer:**
```yaml
# nginx load balancer
upstream backend {
    server blue:8080;
    server green:8080 backup;
}
```

### 107. How do you implement canary deployments with Docker?

**Answer:**

**Method 1: Multiple Containers**
```bash
# Run stable version (90% traffic)
docker run -d --name app-stable-1 -p 8080:80 app:v1.0
docker run -d --name app-stable-2 -p 8081:80 app:v1.0
docker run -d --name app-stable-3 -p 8082:80 app:v1.0

# Run canary version (10% traffic)
docker run -d --name app-canary -p 8083:80 app:v1.1

# Configure load balancer to route 10% to canary
```

**Method 2: Docker Swarm**
```yaml
version: '3.8'
services:
  app:
    image: app:v1.0
    deploy:
      replicas: 9
      update_config:
        parallelism: 1
        delay: 10s
      rollback_config:
        parallelism: 1
  app-canary:
    image: app:v1.1
    deploy:
      replicas: 1
```

### 108. How do you handle secrets in Docker Compose?

**Answer:**

**Docker Swarm Secrets:**
```bash
# Create secret
echo "my-secret" | docker secret create db_password -

# Use in stack
version: '3.8'
services:
  app:
    image: myapp
    secrets:
      - db_password
secrets:
  db_password:
    external: true
```

**Environment Variables:**
```yaml
services:
  app:
    image: myapp
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
    env_file:
      - .env
```

**External Secret Management:**
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Docker Secrets (Swarm only)

### 109. How do you implement health checks in production?

**Answer:**

**Dockerfile:**
```dockerfile
FROM nginx
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD curl -f http://localhost/health || exit 1
```

**docker-compose:**
```yaml
services:
  web:
    image: nginx
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

**Docker Swarm:**
```yaml
services:
  web:
    image: nginx
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      update_config:
        failure_action: rollback
```

### 110. How do you implement log aggregation in Docker?

**Answer:**

**Log Drivers:**
```json
{
  "log-driver": "fluentd",
  "log-opts": {
    "fluentd-address": "localhost:24224",
    "tag": "docker.{{.Name}}"
  }
}
```

**ELK Stack:**
```yaml
version: '3.8'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.0.0
  logstash:
    image: docker.elastic.co/logstash/logstash:8.0.0
  kibana:
    image: docker.elastic.co/kibana/kibana:8.0.0
  app:
    image: myapp
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

**Fluentd:**
```bash
docker run --log-driver=fluentd \
  --log-opt fluentd-address=localhost:24224 \
  --log-opt tag=docker.app \
  nginx
```

### 111. How do you implement monitoring for Docker containers?

**Answer:**

**Prometheus + cAdvisor:**
```yaml
version: '3.8'
services:
  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    ports:
      - "8080:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:ro
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

**Docker Stats:**
```bash
# Real-time stats
docker stats

# Export metrics
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

### 112. How do you implement resource limits effectively?

**Answer:**

**CPU Limits:**
```bash
# Limit to 1 CPU
docker run --cpus="1.0" nginx

# Limit to 50% of 1 CPU
docker run --cpus="0.5" nginx

# CPU shares (relative)
docker run --cpu-shares=512 nginx
```

**Memory Limits:**
```bash
# Memory limit
docker run --memory="512m" nginx

# Memory + swap
docker run --memory="512m" --memory-swap="1g" nginx

# OOM kill priority
docker run --oom-kill-disable nginx  # Not recommended
```

**docker-compose:**
```yaml
services:
  app:
    image: nginx
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 113. How do you implement Docker in CI/CD pipelines?

**Answer:**

**GitHub Actions Example:**
```yaml
name: CI/CD Pipeline
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Build image
      run: docker build -t myapp:${{ github.sha }} .
    
    - name: Run tests
      run: docker run myapp:${{ github.sha }} npm test
    
    - name: Scan image
      run: trivy image myapp:${{ github.sha }}
    
    - name: Push to registry
      run: |
        docker login -u ${{ secrets.DOCKER_USERNAME }} -p ${{ secrets.DOCKER_PASSWORD }}
        docker push myapp:${{ github.sha }}
    
    - name: Deploy
      run: |
        kubectl set image deployment/myapp app=myapp:${{ github.sha }}
```

**Best Practices:**
- Multi-stage builds
- Cache layers
- Parallel builds
- Security scanning
- Tag with commit SHA

### 114. How do you implement Docker image versioning?

**Answer:**

**Versioning Strategies:**
1. **Semantic Versioning:**
   ```bash
   docker tag myapp:latest myapp:1.2.3
   docker tag myapp:latest myapp:1.2
   docker tag myapp:latest myapp:1
   ```

2. **Git SHA:**
   ```bash
   docker tag myapp:latest myapp:$(git rev-parse --short HEAD)
   ```

3. **Build Number:**
   ```bash
   docker tag myapp:latest myapp:build-123
   ```

4. **Date-based:**
   ```bash
   docker tag myapp:latest myapp:20240109
   ```

**Best Practices:**
- Use semantic versioning
- Tag with multiple tags
- Keep `latest` updated
- Don't reuse tags
- Document versioning strategy

### 115. How do you implement Docker image caching strategies?

**Answer:**

**Build Cache:**
```dockerfile
# Order matters - copy dependencies first
COPY package.json package-lock.json ./
RUN npm install

# Copy application code last
COPY . .
```

**Cache Mounts (BuildKit):**
```dockerfile
# syntax=docker/dockerfile:1.4
FROM node:18
RUN --mount=type=cache,target=/root/.npm \
    npm install
```

**Registry Caching:**
```bash
# Pull base images
docker pull node:18-alpine

# Use in Dockerfile
FROM node:18-alpine
```

**Best Practices:**
- Order Dockerfile instructions
- Use cache mounts
- Leverage layer caching
- Use specific base image tags

### 116. How do you troubleshoot Docker build performance?

**Answer:**

**Optimization:**
1. **Use BuildKit:**
   ```bash
   DOCKER_BUILDKIT=1 docker build .
   ```

2. **Parallel builds:**
   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 .
   ```

3. **Cache from registry:**
   ```bash
   docker build --cache-from myapp:latest .
   ```

4. **Use .dockerignore:**
   - Exclude unnecessary files
   - Reduces build context

5. **Multi-stage builds:**
   - Separate build and runtime
   - Smaller final images

**Diagnosis:**
```bash
# Build with timing
time docker build .

# Build with progress
docker build --progress=plain .
```

### 117. How do you implement Docker in multi-cloud environments?

**Answer:**

**Strategies:**
1. **Container Registry:**
   - Use cloud-agnostic registry
   - Push to multiple registries
   - Mirror images

2. **Orchestration:**
   - Kubernetes (works on all clouds)
   - Docker Swarm
   - Nomad

3. **Configuration Management:**
   - Environment-specific configs
   - Use environment variables
   - Cloud-agnostic images

**Example:**
```bash
# Build once
docker build -t myapp:latest .

# Push to multiple registries
docker tag myapp:latest gcr.io/project/myapp:latest
docker tag myapp:latest ecr.amazonaws.com/repo/myapp:latest
docker tag myapp:latest registry.azure.io/repo/myapp:latest

docker push gcr.io/project/myapp:latest
docker push ecr.amazonaws.com/repo/myapp:latest
docker push registry.azure.io/repo/myapp:latest
```

### 118. How do you implement Docker for microservices?

**Answer:**

**Architecture:**
```yaml
version: '3.8'
services:
  api-gateway:
    image: nginx
    ports:
      - "80:80"
    depends_on:
      - user-service
      - order-service
  
  user-service:
    image: user-service:latest
    environment:
      - DB_HOST=user-db
  
  order-service:
    image: order-service:latest
    environment:
      - DB_HOST=order-db
  
  user-db:
    image: postgres:14
    volumes:
      - user-db-data:/var/lib/postgresql/data
  
  order-db:
    image: postgres:14
    volumes:
      - order-db-data:/var/lib/postgresql/data

volumes:
  user-db-data:
  order-db-data:
```

**Best Practices:**
- One container per service
- Independent scaling
- Service discovery
- Health checks
- Logging strategy

### 119. How do you implement Docker for batch processing?

**Answer:**

**Batch Job Container:**
```dockerfile
FROM python:3.9
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY batch_script.py .
CMD ["python", "batch_script.py"]
```

**Run as Job:**
```bash
# One-time job
docker run --rm my-batch-job:latest

# Scheduled job (cron)
# Add to crontab
0 2 * * * docker run --rm my-batch-job:latest
```

**Kubernetes Job:**
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: batch-job
spec:
  template:
    spec:
      containers:
      - name: batch
        image: my-batch-job:latest
      restartPolicy: Never
```

### 120. How do you implement Docker for data processing pipelines?

**Answer:**

**Pipeline Architecture:**
```yaml
version: '3.8'
services:
  data-ingest:
    image: ingest-service:latest
    volumes:
      - data:/data
  
  data-process:
    image: process-service:latest
    depends_on:
      - data-ingest
    volumes:
      - data:/data
  
  data-load:
    image: load-service:latest
    depends_on:
      - data-process
    volumes:
      - data:/data

volumes:
  data:
```

**Best Practices:**
- Use volumes for data sharing
- Health checks between stages
- Error handling
- Monitoring
- Retry logic

---

## Advanced Docker Operations

### 121. How do you implement Docker image signing?

**Answer:**

**Docker Content Trust:**
```bash
# Enable content trust
export DOCKER_CONTENT_TRUST=1

# Sign image
docker trust sign myapp:1.0.0

# Verify signature
docker pull myapp:1.0.0
```

**Notary:**
```bash
# Initialize notary
docker trust key generate mykey

# Sign with notary
docker trust signer add mykey myregistry/myapp
docker trust sign myregistry/myapp:1.0.0
```

**Benefits:**
- Image integrity
- Publisher verification
- Prevents tampering
- Production security

### 122. How do you implement Docker image scanning in production?

**Answer:**

**Automated Scanning:**
```yaml
# CI/CD pipeline
- name: Build
  run: docker build -t myapp:${{ github.sha }} .

- name: Scan
  run: |
    trivy image --exit-code 1 --severity HIGH,CRITICAL myapp:${{ github.sha }}
    # Fails build on high/critical vulnerabilities

- name: Push
  if: success()
  run: docker push myapp:${{ github.sha }}
```

**Runtime Scanning:**
- Integrate with registry
- Scan on push
- Block vulnerable images
- Regular scheduled scans

**Tools:**
- Trivy
- Docker Scout
- Snyk
- Clair

### 123. How do you implement Docker for serverless workloads?

**Answer:**

**AWS Lambda:**
```dockerfile
FROM public.ecr.aws/lambda/python:3.9
COPY app.py ${LAMBDA_TASK_ROOT}
CMD [ "app.handler" ]
```

**Azure Container Instances:**
```bash
az container create \
  --resource-group mygroup \
  --name mycontainer \
  --image myapp:latest \
  --cpu 1 \
  --memory 1
```

**Google Cloud Run:**
```bash
gcloud run deploy myapp \
  --image gcr.io/project/myapp:latest \
  --platform managed
```

**Benefits:**
- Pay per use
- Auto-scaling
- No infrastructure management
- Fast deployment

### 124. How do you implement Docker for machine learning workloads?

**Answer:**

**ML Container:**
```dockerfile
FROM tensorflow/tensorflow:latest
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY model.py .
COPY train.py .
CMD ["python", "train.py"]
```

**GPU Support:**
```bash
# Install nvidia-docker
# Run with GPU
docker run --gpus all tensorflow/tensorflow:latest-gpu python train.py
```

**Best Practices:**
- Use GPU-enabled images
- Mount data volumes
- Optimize image size
- Cache model artifacts
- Use multi-stage builds

### 125. How do you implement Docker for IoT devices?

**Answer:**

**ARM Images:**
```dockerfile
FROM arm32v7/python:3.9-slim
WORKDIR /app
COPY . .
CMD ["python", "app.py"]
```

**Build for ARM:**
```bash
# Using buildx
docker buildx build --platform linux/arm/v7 -t myapp:arm .
```

**Deployment:**
- Lightweight images
- ARM-compatible
- Resource constraints
- Edge computing

**Considerations:**
- Limited resources
- Network constraints
- Power consumption
- Update mechanisms

---

## Docker Best Practices & Patterns

### 126. What are Docker anti-patterns to avoid?

**Answer:**

**Common Anti-patterns:**
1. **Running as root:**
   ```dockerfile
   # Bad
   FROM ubuntu
   RUN apt-get install nginx
   
   # Good
   FROM ubuntu
   RUN useradd -m appuser
   USER appuser
   ```

2. **Using latest tag:**
   ```dockerfile
   # Bad
   FROM nginx:latest
   
   # Good
   FROM nginx:1.21.6
   ```

3. **Large images:**
   ```dockerfile
   # Bad - includes build tools
   FROM node:18
   COPY . .
   RUN npm install && npm build
   CMD ["node", "app.js"]
   
   # Good - multi-stage
   FROM node:18 AS builder
   COPY . .
   RUN npm build
   FROM node:18-alpine
   COPY --from=builder /app/dist .
   ```

4. **Sensitive data in images:**
   - Never commit secrets
   - Use secrets management
   - Use .dockerignore

5. **Single process per container:**
   - One process per container
   - Use init systems if needed
   - Don't run multiple services

### 127. How do you implement Docker for development environments?

**Answer:**

**Development Setup:**
```yaml
version: '3.8'
services:
  app:
    build: .
    volumes:
      - .:/app  # Mount source code
      - /app/node_modules  # Exclude node_modules
    environment:
      - NODE_ENV=development
    ports:
      - "3000:3000"
    command: npm run dev  # Hot reload
  
  db:
    image: postgres:14
    environment:
      - POSTGRES_DB=myapp
    volumes:
      - db-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  db-data:
```

**Benefits:**
- Consistent environment
- Easy setup
- Isolated dependencies
- Quick reset

### 128. How do you implement Docker for testing?

**Answer:**

**Test Containers:**
```yaml
version: '3.8'
services:
  app:
    build: .
  
  test:
    build:
      context: .
      dockerfile: Dockerfile.test
    depends_on:
      - app
      - db
    environment:
      - TEST_DB_HOST=db
    command: pytest
  
  db:
    image: postgres:14
    environment:
      - POSTGRES_DB=testdb
```

**CI/CD Integration:**
```yaml
- name: Run tests
  run: |
    docker-compose -f docker-compose.test.yml up --abort-on-container-exit
    docker-compose -f docker-compose.test.yml down
```

### 129. How do you implement Docker for database migrations?

**Answer:**

**Migration Container:**
```dockerfile
FROM postgres:14
WORKDIR /migrations
COPY migrations/ .
COPY migrate.sh .
CMD ["./migrate.sh"]
```

**Run Migrations:**
```bash
# Before app starts
docker run --rm \
  --network my-network \
  -e DB_HOST=postgres \
  migration-container:latest

# Or as init container in Kubernetes
```

**Best Practices:**
- Idempotent migrations
- Version control
- Backup before migration
- Test migrations
- Rollback plan

### 130. How do you implement Docker for API gateways?

**Answer:**

**Nginx API Gateway:**
```dockerfile
FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY upstream.conf /etc/nginx/conf.d/
EXPOSE 80
```

**nginx.conf:**
```nginx
upstream api {
    server api-service:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://api;
    }
}
```

**Kong API Gateway:**
```yaml
services:
  kong:
    image: kong:latest
    environment:
      - KONG_DATABASE=postgres
      - KONG_PG_HOST=kong-db
    ports:
      - "8000:8000"
      - "8443:8443"
```

### 131. How do you implement Docker for message queues?

**Answer:**

**RabbitMQ:**
```yaml
services:
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq
    environment:
      - RABBITMQ_DEFAULT_USER=admin
      - RABBITMQ_DEFAULT_PASS=password

volumes:
  rabbitmq-data:
```

**Redis:**
```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes

volumes:
  redis-data:
```

### 132. How do you implement Docker for caching layers?

**Answer:**

**BuildKit Cache:**
```dockerfile
# syntax=docker/dockerfile:1.4
FROM node:18
RUN --mount=type=cache,target=/root/.npm \
    npm install
```

**Registry Cache:**
```bash
# Build with cache from registry
docker build \
  --cache-from myapp:latest \
  -t myapp:new .
```

**Local Cache:**
- Docker layer caching
- Build context caching
- Volume caching

### 133. How do you implement Docker for load testing?

**Answer:**

**Load Testing Container:**
```dockerfile
FROM alpine:latest
RUN apk add --no-cache curl
COPY loadtest.sh .
CMD ["./loadtest.sh"]
```

**Load Test:**
```bash
# Run load test
docker run --rm \
  --network my-network \
  loadtest:latest \
  --target http://app:8080 \
  --users 100 \
  --duration 5m
```

**Tools:**
- Apache Bench (ab)
- wrk
- k6
- Locust
- JMeter

### 134. How do you implement Docker for reverse proxies?

**Answer:**

**Nginx Reverse Proxy:**
```dockerfile
FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY proxy.conf /etc/nginx/conf.d/
```

**nginx.conf:**
```nginx
upstream backend {
    server app1:8080;
    server app2:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Traefik:**
```yaml
services:
  traefik:
    image: traefik:latest
    command:
      - --api.insecure=true
      - --providers.docker=true
    ports:
      - "80:80"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

### 135. How do you implement Docker for service discovery?

**Answer:**

**Docker DNS:**
- Automatic service discovery
- Containers resolve by name
- Works within same network

**Consul:**
```yaml
services:
  consul:
    image: consul:latest
    ports:
      - "8500:8500"
    command: consul agent -dev -client 0.0.0.0
```

**etcd:**
```yaml
services:
  etcd:
    image: quay.io/coreos/etcd:latest
    ports:
      - "2379:2379"
    environment:
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
```

**Best Practices:**
- Use Docker DNS for simple cases
- Use Consul/etcd for complex scenarios
- Health checks
- Service registration

---

*This is the fifth batch of 30 questions (106-135). More questions will be added in subsequent batches.*


