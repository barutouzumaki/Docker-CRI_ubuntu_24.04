# Kubernetes Interview Questions

This document contains commonly asked Kubernetes interview questions organized by topic.

## Table of Contents

- [Basic Concepts](#basic-concepts)
- [Architecture & Components](#architecture--components)
- [Pods & Containers](#pods--containers)
- [Services & Networking](#services--networking)
- [Deployments & Scaling](#deployments--scaling)
- [Storage](#storage)
- [Security](#security)
- [Troubleshooting](#troubleshooting)
- [Advanced Topics](#advanced-topics)

---

## Basic Concepts

### 1. What is Kubernetes and why is it used?

**Answer:**
Kubernetes is an open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications. It's used because:
- **Automated scaling**: Automatically scales applications based on demand
- **Self-healing**: Restarts failed containers and replaces them
- **Service discovery**: Automatically discovers and load balances services
- **Rolling updates**: Updates applications without downtime
- **Resource management**: Efficiently manages compute resources
- **Portability**: Works across different cloud providers and on-premises

### 2. What is a Pod in Kubernetes?

**Answer:**
A Pod is the smallest deployable unit in Kubernetes. It represents a single instance of a running process in the cluster. A Pod can contain:
- One or more containers that share storage and network
- Shared IP address and port space
- Shared volumes for data persistence
- Containers in a Pod are always co-located and co-scheduled

**Example:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
```

### 3. What is the difference between a Pod and a Container?

**Answer:**
- **Container**: A lightweight, standalone executable package that includes everything needed to run an application
- **Pod**: A Kubernetes abstraction that wraps one or more containers with shared storage/network and a specification for how to run them

Key differences:
- Pods are Kubernetes-specific; containers are platform-agnostic
- Pods can contain multiple containers that work together
- Pods have Kubernetes metadata and lifecycle management
- Containers in a Pod share the same network namespace and can communicate via localhost

### 4. What is a Namespace in Kubernetes?

**Answer:**
A Namespace is a virtual cluster within a physical Kubernetes cluster. It provides:
- **Resource isolation**: Separate resources into logical groups
- **Access control**: Apply RBAC policies per namespace
- **Resource quotas**: Limit resource consumption per namespace
- **Scoping**: Organize resources (pods, services, deployments) into namespaces

**Default namespaces:**
- `default`: Default namespace for resources
- `kube-system`: System components
- `kube-public`: Publicly accessible resources
- `kube-node-lease`: Node heartbeat information

### 5. Explain the difference between Deployment, ReplicaSet, and Pod.

**Answer:**
- **Pod**: Basic unit that runs containers
- **ReplicaSet**: Ensures a specified number of pod replicas are running
- **Deployment**: Higher-level abstraction that manages ReplicaSets and provides declarative updates

**Hierarchy:**
```
Deployment → ReplicaSet → Pods
```

**Deployment** provides:
- Rolling updates and rollbacks
- ReplicaSet management
- Declarative configuration
- History of revisions

**ReplicaSet** provides:
- Pod replication
- Pod health monitoring
- Pod replacement

### 6. What is a Service in Kubernetes?

**Answer:**
A Service is an abstraction that defines a logical set of Pods and a policy to access them. It provides:
- **Stable IP address**: Even if pods are recreated
- **Load balancing**: Distributes traffic across pods
- **Service discovery**: DNS-based service discovery
- **Decoupling**: Pods can be replaced without affecting clients

**Service Types:**
- **ClusterIP**: Internal cluster IP (default)
- **NodePort**: Exposes service on each node's IP at a static port
- **LoadBalancer**: Exposes service externally using cloud provider's load balancer
- **ExternalName**: Maps service to external DNS name

### 7. What is the difference between StatefulSet and Deployment?

**Answer:**

| Feature | Deployment | StatefulSet |
|---------|-----------|-------------|
| **Pod Identity** | No stable identity | Stable network identity |
| **Pod Names** | Random | Ordered, predictable |
| **Storage** | Shared volumes | Persistent volumes per pod |
| **Scaling** | Any order | Ordered (0, 1, 2...) |
| **Updates** | Rolling update | Ordered updates |
| **Use Case** | Stateless apps | Stateful apps (databases) |

**StatefulSet** is used for:
- Databases (MySQL, PostgreSQL, MongoDB)
- Applications requiring stable network identity
- Applications needing ordered deployment/scaling
- Applications with persistent storage requirements

### 8. What is a ConfigMap and how is it used?

**Answer:**
A ConfigMap is an API object used to store non-confidential data in key-value pairs. It's used to:
- Separate configuration from container images
- Make applications portable across environments
- Store configuration data, command-line arguments, or configuration files

**Usage:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  database_url: "postgresql://localhost:5432/mydb"
  log_level: "info"
```

**Mounting in Pod:**
```yaml
spec:
  containers:
  - name: app
    envFrom:
    - configMapRef:
        name: app-config
```

### 9. What is a Secret in Kubernetes?

**Answer:**
A Secret is an object that contains sensitive data like passwords, tokens, or keys. It's similar to ConfigMap but:
- **Encrypted at rest** (if encryption at rest is enabled)
- **Base64 encoded** (not encrypted, just encoded)
- **Should not be committed** to version control
- **Can be mounted** as files or environment variables

**Types of Secrets:**
- `Opaque`: User-defined data (default)
- `kubernetes.io/dockerconfigjson`: Docker registry credentials
- `kubernetes.io/tls`: TLS certificate and key
- `kubernetes.io/service-account-token`: Service account token

**Best Practice:** Use external secret management tools (Vault, Sealed Secrets) for production.

### 10. What is a PersistentVolume (PV) and PersistentVolumeClaim (PVC)?

**Answer:**
- **PersistentVolume (PV)**: A cluster-wide storage resource provisioned by an administrator or dynamically via StorageClass
- **PersistentVolumeClaim (PVC)**: A request for storage by a user, similar to a Pod requesting CPU/memory

**Lifecycle:**
1. Admin creates PV or StorageClass
2. User creates PVC requesting storage
3. Kubernetes binds PVC to available PV
4. Pod mounts PVC as volume

**Example:**
```yaml
# PVC
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

---

## Architecture & Components

### 11. What are the main components of Kubernetes architecture?

**Answer:**
Kubernetes has two main parts:

**Control Plane (Master) Components:**
- **kube-apiserver**: Front-end for Kubernetes API
- **etcd**: Consistent and highly-available key-value store
- **kube-scheduler**: Watches for new pods and assigns them to nodes
- **kube-controller-manager**: Runs controller processes
- **cloud-controller-manager**: Links cluster to cloud provider

**Node Components:**
- **kubelet**: Agent that runs on each node
- **kube-proxy**: Network proxy maintaining network rules
- **Container Runtime**: Software for running containers (Docker, containerd, CRI-O)

### 12. What is the role of kube-apiserver?

**Answer:**
kube-apiserver is the front-end of the Kubernetes control plane. It:
- **Exposes REST API**: All operations go through API server
- **Validates requests**: Ensures requests are valid
- **Authenticates/Authorizes**: Verifies user identity and permissions
- **Processes requests**: Creates, updates, deletes resources
- **Stores state**: Writes to etcd
- **Serves API**: Provides API endpoints for kubectl and other clients

**Key Features:**
- API versioning
- Admission controllers (validating/mutating)
- API aggregation
- Rate limiting

### 13. What is etcd and why is it important?

**Answer:**
etcd is a distributed, consistent key-value store used as Kubernetes' backing store. It:
- **Stores cluster state**: All cluster data (pods, services, deployments, etc.)
- **Provides consistency**: Ensures data consistency across cluster
- **Enables coordination**: Used for leader election and distributed locking
- **Backup critical**: Must be backed up regularly

**Characteristics:**
- Highly available (typically 3 or 5 nodes)
- Strong consistency
- Watch API for change notifications
- Transaction support

### 14. What does kube-scheduler do?

**Answer:**
kube-scheduler is responsible for assigning pods to nodes. It:
- **Watches for unscheduled pods**: Monitors API server for pods without node assignment
- **Filters nodes**: Finds nodes that meet pod requirements (resources, taints, etc.)
- **Scores nodes**: Ranks nodes based on various factors
- **Binds pods**: Assigns pod to best node

**Scheduling Process:**
1. **Filtering**: Remove nodes that don't meet requirements
2. **Scoring**: Rank remaining nodes
3. **Binding**: Select highest-scoring node

**Factors considered:**
- Resource requests/limits
- Node affinity/anti-affinity
- Taints and tolerations
- Pod affinity/anti-affinity
- Node resources

### 15. What is kubelet and what are its responsibilities?

**Answer:**
kubelet is an agent that runs on each node in the cluster. It:
- **Registers node**: Registers the node with the API server
- **Monitors pods**: Watches for pods assigned to its node
- **Manages containers**: Creates, starts, stops containers via CRI
- **Reports status**: Reports node and pod status to API server
- **Executes probes**: Runs liveness and readiness probes
- **Mounts volumes**: Sets up volumes for pods
- **Manages images**: Pulls container images

**Key responsibilities:**
- Container lifecycle management
- Health monitoring
- Resource reporting
- Volume management

### 16. What is kube-proxy and how does it work?

**Answer:**
kube-proxy is a network proxy that runs on each node. It:
- **Maintains network rules**: Creates iptables/IPVS rules for Services
- **Enables service discovery**: Provides internal cluster IP for services
- **Load balances**: Distributes traffic across service endpoints
- **Supports multiple modes**: iptables, IPVS, userspace

**Modes:**
- **iptables mode** (default): Uses iptables rules for load balancing
- **IPVS mode**: Uses IPVS for better performance
- **userspace mode**: Legacy mode (deprecated)

**How it works:**
1. Watches API server for Service/Endpoint changes
2. Updates iptables/IPVS rules accordingly
3. Routes traffic to backend pods

### 17. What is the Container Runtime Interface (CRI)?

**Answer:**
CRI is a plugin interface that enables kubelet to use a variety of container runtimes without recompiling. It:
- **Standardizes interface**: Defines standard API for container runtimes
- **Enables flexibility**: Allows multiple runtime options
- **Separates concerns**: Decouples Kubernetes from runtime implementation

**CRI Components:**
- **Runtime Service**: Container lifecycle (create, start, stop, remove)
- **Image Service**: Image management (pull, list, remove)

**Supported Runtimes:**
- containerd
- CRI-O
- Docker (via CRI shim)
- Any CRI-compatible runtime

### 18. What is the difference between a Master node and a Worker node?

**Answer:**

| Feature | Master Node | Worker Node |
|---------|-------------|-------------|
| **Components** | API server, etcd, scheduler, controllers | kubelet, kube-proxy, container runtime |
| **Role** | Manages cluster | Runs workloads |
| **Pods** | Usually tainted (no user pods) | Runs application pods |
| **Resources** | Higher CPU/memory | Varies based on workload |
| **Count** | 1-5 (for HA) | Many (scales with workload) |

**Master Node:**
- Hosts control plane components
- Makes cluster-wide decisions
- Stores cluster state
- Can be made highly available

**Worker Node:**
- Runs application containers
- Reports status to master
- Executes workloads
- Scales horizontally

### 19. What is a DaemonSet?

**Answer:**
A DaemonSet ensures that a copy of a pod runs on all (or specific) nodes in the cluster. It's used for:
- **Node-level services**: Logging agents, monitoring agents
- **Network plugins**: CNI plugins
- **Storage systems**: Storage daemons
- **System services**: kube-proxy (managed as DaemonSet)

**Characteristics:**
- One pod per node (or subset)
- Pods are automatically created on new nodes
- Pods are automatically deleted when nodes are removed
- Useful for cluster-wide services

**Example:**
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
spec:
  selector:
    matchLabels:
      name: fluentd
  template:
    metadata:
      labels:
        name: fluentd
    spec:
      containers:
      - name: fluentd
        image: fluentd:latest
```

### 20. What is a Job and CronJob?

**Answer:**
- **Job**: Creates one or more pods and ensures they complete successfully
- **CronJob**: Runs Jobs on a time-based schedule (like cron)

**Job Use Cases:**
- One-time tasks
- Batch processing
- Data processing jobs
- Backup operations

**CronJob Use Cases:**
- Scheduled backups
- Periodic data synchronization
- Scheduled reports
- Maintenance tasks

**Example Job:**
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

**Example CronJob:**
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            command: ["echo", "Hello from CronJob"]
          restartPolicy: OnFailure
```

---

## Services & Networking

### 21. What are the different types of Services in Kubernetes?

**Answer:**
Kubernetes provides four types of Services:

1. **ClusterIP** (Default):
   - Exposes service on cluster-internal IP
   - Only accessible within cluster
   - Use case: Internal services

2. **NodePort**:
   - Exposes service on each node's IP at a static port
   - Accessible from outside cluster via `<NodeIP>:<NodePort>`
   - Port range: 30000-32767
   - Use case: Development, testing

3. **LoadBalancer**:
   - Exposes service externally using cloud provider's load balancer
   - Automatically creates NodePort and ClusterIP
   - Use case: Production external access

4. **ExternalName**:
   - Maps service to external DNS name
   - Returns CNAME record
   - Use case: External services

**Example:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  type: LoadBalancer
  selector:
    app: my-app
  ports:
  - port: 80
    targetPort: 8080
```

### 22. How does Service discovery work in Kubernetes?

**Answer:**
Service discovery in Kubernetes works through DNS:

1. **DNS-based**: Kubernetes has a built-in DNS server (CoreDNS)
2. **Service DNS**: `<service-name>.<namespace>.svc.cluster.local`
3. **Short names**: Can use `<service-name>` in same namespace
4. **Environment variables**: Pods get service IP/port as env vars

**DNS Resolution:**
- Same namespace: `my-service` → resolves to service IP
- Different namespace: `my-service.my-namespace` → resolves to service IP
- Full FQDN: `my-service.my-namespace.svc.cluster.local`

**Example:**
```bash
# From a pod in same namespace
curl http://my-service

# From different namespace
curl http://my-service.production.svc.cluster.local
```

### 23. What is an Ingress in Kubernetes?

**Answer:**
Ingress is an API object that manages external HTTP/HTTPS access to services. It provides:
- **URL-based routing**: Route traffic based on URL path
- **Host-based routing**: Route based on hostname
- **SSL/TLS termination**: HTTPS termination
- **Load balancing**: Distribute traffic

**Components:**
- **Ingress Controller**: Implements Ingress (nginx, traefik, etc.)
- **Ingress Resource**: Defines routing rules

**Example:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
spec:
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
```

### 24. What is the difference between Ingress and LoadBalancer?

**Answer:**

| Aspect | Ingress | LoadBalancer |
|--------|---------|--------------|
| **Protocol** | HTTP/HTTPS only | Any protocol (TCP, UDP, etc.) |
| **Layer** | Layer 7 (Application) | Layer 4 (Transport) |
| **Features** | URL routing, SSL termination | Simple load balancing |
| **Cost** | One LoadBalancer for many services | One LoadBalancer per service |
| **Use Case** | Web applications | Any TCP/UDP service |

**Ingress:**
- Single entry point for multiple services
- URL-based routing
- SSL termination
- More cost-effective

**LoadBalancer:**
- Direct external access
- Works with any protocol
- One per service
- Simpler setup

### 25. What is a NetworkPolicy?

**Answer:**
NetworkPolicy is a Kubernetes resource that controls traffic flow between pods. It:
- **Defines rules**: Allow/deny traffic based on labels
- **Pod isolation**: Isolate pods by default
- **Ingress rules**: Control incoming traffic
- **Egress rules**: Control outgoing traffic

**Requirements:**
- Network plugin must support NetworkPolicy (Calico, Cilium, etc.)
- Flannel does NOT support NetworkPolicy

**Example:**
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
```

### 26. What is CNI (Container Network Interface)?

**Answer:**
CNI is a specification and libraries for writing plugins to configure network interfaces in Linux containers. It:
- **Standardizes networking**: Common interface for network plugins
- **Plugin-based**: Multiple implementations (Flannel, Calico, Cilium)
- **Pod networking**: Configures network for pods
- **IP management**: Allocates IP addresses

**Popular CNI Plugins:**
- **Flannel**: Simple overlay network
- **Calico**: BGP-based with network policies
- **Cilium**: eBPF-based networking and security
- **Weave**: Overlay network with encryption

**How it works:**
1. kubelet calls CNI plugin when pod is created
2. Plugin configures network interface
3. Plugin allocates IP address
4. Plugin sets up routing

---

## Deployments & Scaling

### 27. What is the difference between Rolling Update and Recreate deployment strategy?

**Answer:**

| Aspect | Rolling Update | Recreate |
|--------|---------------|----------|
| **Downtime** | Zero downtime | Brief downtime |
| **Process** | Gradual replacement | Delete all, create new |
| **Resources** | Higher (old + new pods) | Lower (only new pods) |
| **Use Case** | Production | Development, stateful apps |
| **Rollback** | Easy | Requires recreation |

**Rolling Update:**
- Gradually replaces old pods with new ones
- Maintains service availability
- Can specify `maxSurge` and `maxUnavailable`
- Default strategy

**Recreate:**
- Terminates all old pods before creating new ones
- Brief service interruption
- Useful for stateful applications that can't run multiple versions

**Example:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  replicas: 3
  template:
    # pod template
```

### 28. How does Horizontal Pod Autoscaler (HPA) work?

**Answer:**
HPA automatically scales the number of pods based on observed metrics. It:
- **Monitors metrics**: CPU, memory, custom metrics
- **Scales up/down**: Increases/decreases replica count
- **Target utilization**: Maintains target metric value
- **Periodic evaluation**: Checks metrics every 15 seconds (default)

**Requirements:**
- Metrics Server installed
- Resource requests defined in pod spec
- HPA resource created

**Example:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: my-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

**Scaling Behavior:**
- CPU > 70%: Scale up
- CPU < 70%: Scale down
- Respects min/max replicas

### 29. What is Vertical Pod Autoscaler (VPA)?

**Answer:**
VPA automatically adjusts pod resource requests and limits based on historical usage. It:
- **Adjusts resources**: CPU/memory requests and limits
- **Right-sizing**: Optimizes resource allocation
- **Learning**: Analyzes historical usage patterns
- **Modes**: Off, Initial, Auto, Recreate

**VPA Modes:**
- **Off**: Only provides recommendations
- **Initial**: Sets resources at pod creation
- **Auto**: Updates resources on running pods (requires recreation)
- **Recreate**: Recreates pods with new resources

**Use Case:**
- Optimize resource usage
- Reduce over-provisioning
- Handle varying workloads

**Note:** VPA and HPA on CPU/memory cannot be used together.

### 30. What is Cluster Autoscaler?

**Answer:**
Cluster Autoscaler automatically adjusts the size of the Kubernetes cluster. It:
- **Adds nodes**: When pods can't be scheduled due to insufficient resources
- **Removes nodes**: When nodes are underutilized
- **Cloud integration**: Works with cloud providers (AWS, GCP, Azure)
- **Node groups**: Manages node groups/pools

**How it works:**
1. Detects unschedulable pods
2. Requests new nodes from cloud provider
3. Waits for nodes to become ready
4. Pods get scheduled on new nodes
5. Removes nodes when underutilized (after grace period)

**Requirements:**
- Cloud provider with autoscaling support
- Node groups configured
- Appropriate IAM permissions

---

## Storage

### 31. What is the difference between PersistentVolume and PersistentVolumeClaim?

**Answer:**

| Aspect | PersistentVolume (PV) | PersistentVolumeClaim (PVC) |
|--------|----------------------|---------------------------|
| **Level** | Cluster resource | Namespace resource |
| **Created by** | Admin or StorageClass | User/Developer |
| **Binding** | Bound to PVC | Requests PV |
| **Lifecycle** | Independent | Depends on PV |

**PersistentVolume (PV):**
- Cluster-wide storage resource
- Provisioned by admin or dynamically
- Has access modes and storage capacity
- Can be statically or dynamically provisioned

**PersistentVolumeClaim (PVC):**
- Request for storage by user
- Specifies size and access mode
- Kubernetes binds to available PV
- Used in pod specs

**Example:**
```yaml
# PV (created by admin)
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-1
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  hostPath:
    path: /data/pv1

# PVC (created by user)
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-1
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

### 32. What is a StorageClass?

**Answer:**
StorageClass provides a way to describe different "classes" of storage. It:
- **Dynamic provisioning**: Automatically creates PVs when PVCs are created
- **Storage types**: Different classes (SSD, HDD, fast, slow)
- **Provisioner**: Defines which volume plugin to use
- **Parameters**: Storage-specific parameters

**Benefits:**
- No manual PV creation
- Automatic provisioning
- Different storage tiers
- Cost optimization

**Example:**
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
  fsType: ext4
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  storageClassName: fast-ssd
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
```

### 33. What are the different access modes for PersistentVolumes?

**Answer:**
Access modes define how a volume can be mounted:

1. **ReadWriteOnce (RWO)**:
   - Can be mounted as read-write by single node
   - Most common mode
   - Use case: Database storage

2. **ReadOnlyMany (ROX)**:
   - Can be mounted read-only by many nodes
   - Use case: Shared configuration, read-only data

3. **ReadWriteMany (RWX)**:
   - Can be mounted as read-write by many nodes
   - Requires storage that supports concurrent access
   - Use case: Shared file systems

4. **ReadWriteOncePod (RWOP)** (Kubernetes 1.22+):
   - Can be mounted as read-write by single pod
   - More restrictive than RWO
   - Use case: Pod-specific storage

**Note:** Not all storage backends support all modes.

---

## Security

### 34. What is RBAC in Kubernetes?

**Answer:**
RBAC (Role-Based Access Control) is a method of regulating access to resources based on roles. It provides:
- **Fine-grained control**: Control who can do what
- **Role**: Defines permissions (what can be done)
- **RoleBinding**: Binds role to users/groups/service accounts
- **Namespaced/Cluster**: Roles (namespaced) or ClusterRoles (cluster-wide)

**Components:**
- **Role/ClusterRole**: Defines permissions
- **RoleBinding/ClusterRoleBinding**: Grants permissions

**Example:**
```yaml
# Role
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: default
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]

# RoleBinding
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

### 35. What is a ServiceAccount in Kubernetes?

**Answer:**
ServiceAccount provides an identity for pods. It:
- **Pod identity**: Identifies pods to API server
- **RBAC**: Used with RBAC for authorization
- **Secrets**: Automatically mounts secrets
- **Default**: Each namespace has default service account

**Use Cases:**
- Pod authentication
- RBAC authorization
- Image pull secrets
- External API access

**Example:**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-sa
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: my-sa
  containers:
  - name: app
    image: nginx
```

### 36. What are Pod Security Standards?

**Answer:**
Pod Security Standards define three policies for pod security:

1. **Privileged** (Most Permissive):
   - No restrictions
   - Use case: System components

2. **Baseline** (Recommended):
   - Prevents known privilege escalations
   - Use case: Most applications

3. **Restricted** (Most Secure):
   - Hardened security
   - Use case: Security-sensitive applications

**Enforcement Modes:**
- **enforce**: Policy violations are rejected
- **audit**: Policy violations are logged
- **warn**: Policy violations trigger warnings

**Example:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### 37. What is a SecurityContext?

**Answer:**
SecurityContext defines security settings for pods and containers. It controls:
- **User/Group**: Run as non-root user
- **Capabilities**: Linux capabilities
- **SELinux/AppArmor**: Security modules
- **Read-only filesystem**: Immutable filesystem
- **Privilege escalation**: Allow/deny

**Pod-level vs Container-level:**
- **Pod-level**: Applies to all containers
- **Container-level**: Overrides pod-level for specific container

**Example:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-pod
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

### 38. What is the difference between Secrets and ConfigMaps?

**Answer:**

| Aspect | Secrets | ConfigMaps |
|--------|---------|------------|
| **Data Type** | Sensitive data | Non-sensitive data |
| **Encoding** | Base64 encoded | Plain text |
| **Use Case** | Passwords, tokens, keys | Configuration, env vars |
| **Size Limit** | 1MB | 1MB |
| **Storage** | Can be encrypted at rest | Plain text |

**Secrets:**
- For sensitive data
- Base64 encoded (not encrypted)
- Should use external secret management in production
- Types: Opaque, docker-registry, tls

**ConfigMaps:**
- For non-sensitive configuration
- Plain text
- Can be mounted as files or env vars
- Version controlled

**Best Practice:**
- Use Secrets for sensitive data
- Use ConfigMaps for configuration
- Consider external secret management (Vault) for production

---

## Troubleshooting

### 39. How do you debug a pod that is in CrashLoopBackOff?

**Answer:**
Steps to debug CrashLoopBackOff:

1. **Check pod logs:**
   ```bash
   kubectl logs <pod-name> -n <namespace>
   kubectl logs <pod-name> -n <namespace> --previous
   ```

2. **Describe pod:**
   ```bash
   kubectl describe pod <pod-name> -n <namespace>
   ```

3. **Check events:**
   ```bash
   kubectl get events --sort-by='.lastTimestamp' -n <namespace>
   ```

4. **Check container exit code:**
   ```bash
   kubectl get pod <pod-name> -o jsonpath='{.status.containerStatuses[0].lastState.terminated.exitCode}'
   ```

5. **Exec into container (if running):**
   ```bash
   kubectl exec -it <pod-name> -n <namespace> -- /bin/sh
   ```

6. **Check resource limits:**
   ```bash
   kubectl describe pod <pod-name> | grep -A 5 "Limits:"
   ```

**Common Causes:**
- Application errors (check logs)
- Configuration errors
- Resource limits exceeded
- Health check failures
- Missing dependencies

### 40. How do you troubleshoot a pod stuck in Pending state?

**Answer:**
Steps to troubleshoot Pending pods:

1. **Describe pod:**
   ```bash
   kubectl describe pod <pod-name> -n <namespace>
   ```
   Look for "Events" section

2. **Check node resources:**
   ```bash
   kubectl top nodes
   kubectl describe node <node-name>
   ```

3. **Check node status:**
   ```bash
   kubectl get nodes
   kubectl get nodes -o wide
   ```

4. **Check taints and tolerations:**
   ```bash
   kubectl describe node | grep Taints
   kubectl describe pod <pod-name> | grep -A 5 "Tolerations:"
   ```

5. **Check pod resource requests:**
   ```bash
   kubectl describe pod <pod-name> | grep -A 5 "Requests:"
   ```

**Common Causes:**
- Insufficient resources (CPU/memory)
- Node not ready
- Taints preventing scheduling
- No nodes match node selector/affinity
- PVC not bound
- Image pull errors

---

## Advanced Topics

### 41. What are Taints and Tolerations?

**Answer:**
Taints and Tolerations work together to ensure pods are not scheduled on inappropriate nodes.

**Taints:**
- Applied to nodes
- Prevents pods from being scheduled (unless they tolerate the taint)
- Three effects: `NoSchedule`, `PreferNoSchedule`, `NoExecute`

**Tolerations:**
- Applied to pods
- Allows pods to be scheduled on tainted nodes
- Must match taint key, value, and effect

**Example:**
```yaml
# Taint a node
kubectl taint nodes node1 key=value:NoSchedule

# Pod with toleration
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  tolerations:
  - key: "key"
    operator: "Equal"
    value: "value"
    effect: "NoSchedule"
  containers:
  - name: app
    image: nginx
```

**Use Cases:**
- Dedicated nodes for specific workloads
- Nodes with special hardware (GPU, high memory)
- Preventing user pods on master nodes

### 42. What are Node Affinity and Pod Affinity?

**Answer:**

**Node Affinity:**
- Constrains which nodes a pod can be scheduled on
- Based on node labels
- Types: `requiredDuringSchedulingIgnoredDuringExecution`, `preferredDuringSchedulingIgnoredDuringExecution`

**Pod Affinity:**
- Schedules pods relative to other pods
- Pod Affinity: Co-locate pods
- Pod Anti-Affinity: Separate pods

**Example:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: disktype
            operator: In
            values:
            - ssd
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - database
        topologyKey: kubernetes.io/hostname
  containers:
  - name: app
    image: nginx
```

### 43. What is a Headless Service?

**Answer:**
A Headless Service is a service with `clusterIP: None`. It:
- **No load balancing**: Returns individual pod IPs
- **DNS resolution**: Returns A records for each pod
- **StatefulSets**: Used with StatefulSets for stable network identity
- **Direct pod access**: Access pods directly by DNS name

**Example:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: headless-service
spec:
  clusterIP: None
  selector:
    app: my-app
  ports:
  - port: 80
```

**DNS Resolution:**
- Service: `headless-service.namespace.svc.cluster.local`
- Pods: `pod-0.headless-service.namespace.svc.cluster.local`

**Use Cases:**
- StatefulSets
- Service discovery
- Direct pod-to-pod communication

### 44. What is a Custom Resource Definition (CRD)?

**Answer:**
CRD extends Kubernetes API to add custom resources. It:
- **Extends API**: Adds new resource types
- **Custom controllers**: Can be managed by custom controllers
- **Operators**: Foundation for Kubernetes operators
- **Domain-specific**: Represents domain-specific concepts

**Example:**
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: databases.example.com
spec:
  group: example.com
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              databaseName:
                type: string
  scope: Namespaced
  names:
    plural: databases
    singular: database
    kind: Database
```

**Use Cases:**
- Operators (Prometheus, Istio)
- Custom abstractions
- Domain-specific resources

### 45. What is a Kubernetes Operator?

**Answer:**
A Kubernetes Operator is a method of packaging, deploying, and managing a Kubernetes application. It:
- **Custom controllers**: Uses custom controllers
- **CRDs**: Defines Custom Resources
- **Automation**: Automates complex application management
- **Domain knowledge**: Encodes operational knowledge

**Components:**
- Custom Resource Definition (CRD)
- Custom Controller
- Operator logic

**Example Operators:**
- Prometheus Operator
- Istio Operator
- Elasticsearch Operator
- PostgreSQL Operator

**Benefits:**
- Automates complex operations
- Encodes best practices
- Reduces operational burden

### 46. What is Helm and how is it used?

**Answer:**
Helm is a package manager for Kubernetes. It:
- **Packages**: Packages Kubernetes applications as charts
- **Manages**: Manages application lifecycle
- **Templates**: Uses templates for configuration
- **Releases**: Manages releases and versions

**Key Concepts:**
- **Chart**: Package of Kubernetes resources
- **Release**: Instance of a chart deployed to cluster
- **Repository**: Collection of charts

**Example:**
```bash
# Install Helm chart
helm install my-release stable/nginx

# Upgrade release
helm upgrade my-release stable/nginx

# List releases
helm list

# Uninstall
helm uninstall my-release
```

**Chart Structure:**
```
mychart/
  Chart.yaml
  values.yaml
  templates/
    deployment.yaml
    service.yaml
```

### 47. What are Init Containers?

**Answer:**
Init Containers are containers that run before the main containers in a pod. They:
- **Run sequentially**: Run one after another
- **Must succeed**: All must complete successfully
- **Setup tasks**: Perform setup/initialization
- **Separate from app**: Different from application containers

**Use Cases:**
- Wait for dependencies
- Initialize data
- Setup configuration
- Database migrations

**Example:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  initContainers:
  - name: init-db
    image: busybox
    command: ['sh', '-c', 'until nslookup database; do sleep 2; done']
  containers:
  - name: app
    image: nginx
```

### 48. What are Liveness and Readiness Probes?

**Answer:**

**Liveness Probe:**
- Determines if container is alive
- If fails, container is restarted
- Use when container can become deadlocked

**Readiness Probe:**
- Determines if container is ready to serve traffic
- If fails, pod is removed from service endpoints
- Use when container needs time to start

**Probe Types:**
- `httpGet`: HTTP GET request
- `tcpSocket`: TCP connection
- `exec`: Execute command

**Example:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: app
    image: nginx
    livenessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 30
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
```

### 49. What is Resource Quota?

**Answer:**
ResourceQuota restricts resource consumption per namespace. It limits:
- **Compute resources**: CPU, memory
- **Storage resources**: PersistentVolumeClaims
- **Object counts**: Pods, Services, ConfigMaps
- **Extended resources**: GPU, custom resources

**Example:**
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-quota
  namespace: production
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
    persistentvolumeclaims: "10"
    pods: "10"
```

**Effects:**
- Prevents resource exhaustion
- Enforces limits per namespace
- Requires resource requests/limits in pods

### 50. What is LimitRange?

**Answer:**
LimitRange sets default resource limits and requests for pods in a namespace. It:
- **Default limits**: Sets default CPU/memory limits
- **Default requests**: Sets default CPU/memory requests
- **Min/Max constraints**: Enforces min/max resource values
- **Storage limits**: Limits storage per PVC

**Example:**
```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: mem-limit-range
  namespace: default
spec:
  limits:
  - default:
      memory: "512Mi"
      cpu: "500m"
    defaultRequest:
      memory: "256Mi"
      cpu: "250m"
    max:
      memory: "1Gi"
      cpu: "1000m"
    min:
      memory: "128Mi"
      cpu: "100m"
    type: Container
```

**Benefits:**
- Prevents resource exhaustion
- Sets defaults for pods without limits
- Enforces resource policies

### 51. What is the difference between StatefulSet and Deployment?

**Answer:**

| Feature | Deployment | StatefulSet |
|---------|-----------|-------------|
| **Pod Identity** | No stable identity | Stable network identity |
| **Pod Names** | Random | Ordered, predictable |
| **Storage** | Shared volumes | Persistent volumes per pod |
| **Scaling** | Any order | Ordered (0, 1, 2...) |
| **Updates** | Rolling update | Ordered updates |
| **DNS** | Service DNS | Stable pod DNS |
| **Use Case** | Stateless apps | Stateful apps |

**StatefulSet Features:**
- Stable network identity: `pod-0`, `pod-1`, etc.
- Stable storage: Each pod gets its own PV
- Ordered deployment: Pods created in order
- Ordered scaling: Pods scaled in order
- Ordered deletion: Pods deleted in reverse order

### 52. What is a Pod Disruption Budget (PDB)?

**Answer:**
PDB limits the number of pods that can be voluntarily disrupted. It:
- **Protects availability**: Ensures minimum number of pods
- **Voluntary disruptions**: Applies to evictions, drain operations
- **Min available**: Minimum pods that must be available
- **Max unavailable**: Maximum pods that can be unavailable

**Example:**
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: my-app
```

**Use Cases:**
- High availability applications
- During node maintenance
- Cluster upgrades

### 53. What is a ConfigMap and how is it different from a Secret?

**Answer:**

| Aspect | ConfigMap | Secret |
|--------|-----------|--------|
| **Data Type** | Non-sensitive | Sensitive |
| **Encoding** | Plain text | Base64 |
| **Size Limit** | 1MB | 1MB |
| **Use Case** | Configuration | Credentials, tokens |

**ConfigMap:**
- Stores non-sensitive configuration
- Plain text
- Can be mounted as files or env vars
- Version controlled

**Secret:**
- Stores sensitive data
- Base64 encoded (not encrypted)
- Should use external secret management in production

### 54. What is the Kubernetes API Server?

**Answer:**
kube-apiserver is the front-end of the Kubernetes control plane. It:
- **REST API**: Exposes Kubernetes API
- **Validation**: Validates and processes API requests
- **Authentication**: Verifies user identity
- **Authorization**: Checks permissions (RBAC)
- **Admission Control**: Mutating/validating admission controllers
- **State storage**: Writes to etcd

**Key Features:**
- API versioning
- Rate limiting
- API aggregation
- Watch API for change notifications

### 55. What is etcd and why is it critical?

**Answer:**
etcd is a distributed, consistent key-value store used as Kubernetes' backing store. It:
- **Stores cluster state**: All cluster data
- **Consistency**: Ensures data consistency
- **High availability**: Typically 3 or 5 nodes
- **Watch API**: Notifies on changes
- **Transactions**: Supports transactions

**Critical for:**
- Cluster state persistence
- Configuration storage
- Coordination
- Leader election

**Backup:**
- Must be backed up regularly
- Critical for disaster recovery
- Contains all cluster configuration

### 56. What is the Container Storage Interface (CSI)?

**Answer:**
CSI is a standard for exposing storage systems to containerized workloads. It:
- **Standardizes storage**: Common interface for storage providers
- **Dynamic provisioning**: Automatically provisions storage
- **Vendor independence**: Storage vendors develop plugins
- **Supports**: Block, file, object storage

**CSI Components:**
- **CSI Driver**: Storage vendor implementation
- **External Provisioner**: Handles volume provisioning
- **External Attacher**: Handles volume attachment
- **External Resizer**: Handles volume resizing

**Benefits:**
- Standardized interface
- Vendor independence
- Dynamic provisioning
- Advanced features (snapshots, cloning)

### 57. What is a Service Mesh?

**Answer:**
A Service Mesh is a dedicated infrastructure layer for managing service-to-service communication. It provides:
- **Traffic management**: Load balancing, routing
- **Security**: mTLS, authentication
- **Observability**: Metrics, tracing, logging
- **Policy enforcement**: Rate limiting, access control

**Popular Service Meshes:**
- **Istio**: Most popular, feature-rich
- **Linkerd**: Lightweight, simple
- **Consul Connect**: HashiCorp's solution
- **Kuma**: Universal service mesh

**Components:**
- **Control Plane**: Manages configuration
- **Data Plane**: Proxies (Envoy, Linkerd-proxy)

### 58. What is GitOps?

**Answer:**
GitOps is a methodology for managing infrastructure and applications using Git as the source of truth. It:
- **Git as source of truth**: All changes in Git
- **Declarative**: Desired state in Git
- **Automated**: Automated sync and deployment
- **Observable**: Full audit trail

**Tools:**
- **ArgoCD**: GitOps continuous delivery
- **Flux**: GitOps toolkit
- **Jenkins X**: CI/CD with GitOps

**Benefits:**
- Version control
- Audit trail
- Rollback capability
- Collaboration

### 59. What is the difference between kubectl apply and kubectl create?

**Answer:**

| Aspect | kubectl apply | kubectl create |
|--------|---------------|----------------|
| **Mode** | Declarative | Imperative |
| **Idempotent** | Yes | No (fails if exists) |
| **Manages** | Creates and updates | Only creates |
| **Use Case** | Configuration files | Quick creation |

**kubectl apply:**
- Declarative approach
- Creates or updates resources
- Idempotent (can run multiple times)
- Uses configuration files
- Recommended for production

**kubectl create:**
- Imperative approach
- Only creates resources
- Fails if resource exists
- Quick for testing

**Example:**
```bash
# Apply (declarative)
kubectl apply -f deployment.yaml

# Create (imperative)
kubectl create deployment nginx --image=nginx
```

### 60. What is a Kubernetes Admission Controller?

**Answer:**
Admission Controllers are plugins that intercept requests to the API server. They:
- **Intercept requests**: Before object is stored
- **Mutate**: Modify objects (MutatingAdmissionWebhook)
- **Validate**: Reject invalid objects (ValidatingAdmissionWebhook)
- **Enforce policies**: Security, resource policies

**Types:**
- **Mutating**: Modify objects (e.g., add default values)
- **Validating**: Validate objects (e.g., enforce policies)

**Built-in Controllers:**
- ResourceQuota
- LimitRanger
- PodSecurityPolicy (deprecated)
- NamespaceLifecycle

**Webhooks:**
- MutatingAdmissionWebhook
- ValidatingAdmissionWebhook

### 61. What is the difference between requests and limits?

**Answer:**

| Aspect | Requests | Limits |
|--------|----------|--------|
| **Purpose** | Resource reservation | Maximum resource usage |
| **Scheduling** | Used for scheduling | Not used for scheduling |
| **Guarantee** | Guaranteed minimum | Maximum allowed |
| **OOM Kill** | No | Yes (if exceeded) |

**Requests:**
- Minimum resources guaranteed
- Used by scheduler to place pods
- Pod gets at least this amount
- Can be throttled if not available

**Limits:**
- Maximum resources allowed
- Pod cannot exceed this
- OOM kill if memory limit exceeded
- CPU throttling if CPU limit exceeded

**Example:**
```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### 62. What is a Kubernetes Namespace and when to use it?

**Answer:**
A Namespace is a virtual cluster within a physical cluster. Use it for:
- **Resource organization**: Group related resources
- **Access control**: Apply RBAC per namespace
- **Resource quotas**: Limit resources per namespace
- **Environment separation**: dev, staging, prod
- **Multi-tenancy**: Separate teams/projects

**Default Namespaces:**
- `default`: Default namespace
- `kube-system`: System components
- `kube-public`: Public resources
- `kube-node-lease`: Node heartbeat

**Example:**
```bash
# Create namespace
kubectl create namespace production

# Apply resource to namespace
kubectl apply -f deployment.yaml -n production

# Set default namespace
kubectl config set-context --current --namespace=production
```

### 63. What is the difference between ClusterRole and Role?

**Answer:**

| Aspect | Role | ClusterRole |
|--------|------|-------------|
| **Scope** | Namespace | Cluster-wide |
| **Resources** | Namespaced resources | Cluster and namespaced resources |
| **Binding** | RoleBinding | ClusterRoleBinding |
| **Use Case** | Namespace-specific | Cluster-wide permissions |

**Role:**
- Namespace-scoped
- Applies to resources in namespace
- Bound with RoleBinding

**ClusterRole:**
- Cluster-scoped
- Can grant permissions to cluster resources
- Can grant permissions to all namespaces
- Bound with ClusterRoleBinding

**Example:**
```yaml
# Role (namespaced)
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]

# ClusterRole (cluster-wide)
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cluster-admin
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
```

### 64. What is a Kubernetes Controller?

**Answer:**
A Controller is a control loop that watches cluster state and makes changes to move current state toward desired state. It:
- **Watches resources**: Monitors API server for changes
- **Reconciles**: Compares desired vs current state
- **Takes action**: Creates/updates/deletes resources
- **Maintains state**: Ensures desired state

**Built-in Controllers:**
- Deployment Controller
- ReplicaSet Controller
- StatefulSet Controller
- Job Controller
- Service Controller

**How it works:**
1. Watch API server for changes
2. Compare desired vs current state
3. Take corrective action
4. Repeat

### 65. What is the difference between kubectl get and kubectl describe?

**Answer:**

| Aspect | kubectl get | kubectl describe |
|--------|-------------|------------------|
| **Output** | Table/list | Detailed information |
| **Use Case** | Quick overview | Detailed debugging |
| **Information** | Basic fields | All fields + events |

**kubectl get:**
- Shows basic information in table format
- Quick overview
- Can filter and format output
- Use for listing resources

**kubectl describe:**
- Shows detailed information
- Includes events
- Shows all fields
- Use for debugging

**Example:**
```bash
# Get (quick overview)
kubectl get pods

# Describe (detailed)
kubectl describe pod my-pod
```

---

## Monitoring & Observability

### 66. How do you monitor Kubernetes clusters?

**Answer:**

**Metrics Collection:**
- **Metrics Server**: Core metrics (CPU, memory)
- **Prometheus**: Comprehensive metrics collection
- **cAdvisor**: Container metrics
- **Node Exporter**: Node-level metrics

**Tools:**
- **Prometheus + Grafana**: Metrics and visualization
- **Datadog**: Commercial monitoring
- **New Relic**: APM and monitoring
- **Dynatrace**: Full-stack monitoring

**Key Metrics:**
- Cluster: Node status, pod count, resource usage
- Application: Request rate, error rate, latency
- Infrastructure: CPU, memory, disk, network

**Example:**
```bash
# Install metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# View node metrics
kubectl top nodes

# View pod metrics
kubectl top pods
```

### 67. What is the difference between Liveness, Readiness, and Startup Probes?

**Answer:**

| Probe | Purpose | Action on Failure | When to Use |
|-------|---------|------------------|-------------|
| **Startup** | Container startup | Restarts container | Slow-starting containers |
| **Liveness** | Container alive | Restarts container | Deadlock detection |
| **Readiness** | Ready for traffic | Removes from endpoints | Initialization needed |

**Startup Probe:**
- Disables liveness/readiness until success
- Gives container time to start
- Prevents premature restarts

**Liveness Probe:**
- Checks if container is running
- Restarts if fails
- Use for deadlock detection

**Readiness Probe:**
- Checks if container can serve traffic
- Removes from service if fails
- Use when container needs initialization

### 68. How do you implement logging in Kubernetes?

**Answer:**

**Log Collection:**
- **kubectl logs**: View pod logs
- **Fluentd/Fluent Bit**: Log aggregation
- **Filebeat**: Log shipper
- **Loki**: Log aggregation system

**Centralized Logging:**
- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **EFK Stack**: Elasticsearch, Fluentd, Kibana
- **Loki + Grafana**: Lightweight logging
- **Splunk**: Enterprise logging

**Best Practices:**
- Use structured logging (JSON)
- Include correlation IDs
- Set log levels appropriately
- Rotate logs to prevent disk fill

**Example:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
  - name: app
    image: nginx
    # Logs go to stdout/stderr
    # Collected by kubelet
```

### 69. What is Prometheus Operator?

**Answer:**
Prometheus Operator manages Prometheus and related monitoring components. It:
- **CRDs**: Defines Custom Resources (Prometheus, ServiceMonitor)
- **Automation**: Automates Prometheus configuration
- **Service Discovery**: Auto-discovers targets
- **High Availability**: Manages HA Prometheus

**Components:**
- **Prometheus CRD**: Prometheus instance
- **ServiceMonitor**: Discovers services to monitor
- **Alertmanager**: Handles alerts
- **PrometheusRule**: Defines alerting rules

**Benefits:**
- Declarative configuration
- Automatic target discovery
- Simplified management
- Kubernetes-native

### 70. How do you backup and restore etcd?

**Answer:**

**Backup:**
```bash
# Snapshot etcd
ETCDCTL_API=3 etcdctl snapshot save /backup/etcd-snapshot.db \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key

# Check snapshot
etcdctl snapshot status /backup/etcd-snapshot.db
```

**Restore:**
```bash
# Stop kube-apiserver
systemctl stop kube-apiserver

# Restore from snapshot
ETCDCTL_API=3 etcdctl snapshot restore /backup/etcd-snapshot.db \
  --data-dir=/var/lib/etcd-backup

# Update etcd data directory
# Restart etcd and kube-apiserver
```

**Best Practices:**
- Backup regularly (daily)
- Test restore procedures
- Store backups securely
- Document restore process

---

## Advanced Networking

### 71. What is CoreDNS and how does it work?

**Answer:**
CoreDNS is the default DNS server in Kubernetes. It:
- **Service Discovery**: Resolves service names to IPs
- **Pod DNS**: Provides DNS for pods
- **Plugins**: Extensible plugin architecture
- **Lightweight**: Fast and efficient

**How it works:**
1. Pods query CoreDNS for service names
2. CoreDNS checks Kubernetes API
3. Returns service IP or pod IPs
4. Caches responses

**Configuration:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        health
        kubernetes cluster.local in-addr.arpa ip6.arpa {
            pods insecure
            fallthrough in-addr.arpa ip6.arpa
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
```

### 72. What is the difference between ClusterIP, NodePort, and LoadBalancer Services?

**Answer:**

| Type | Access | IP | Use Case |
|------|--------|-----|----------|
| **ClusterIP** | Internal only | Cluster IP | Internal services |
| **NodePort** | External via node IP | Node IP + port (30000-32767) | Development, testing |
| **LoadBalancer** | External via LB | Cloud LB IP | Production external access |

**ClusterIP:**
- Default service type
- Only accessible within cluster
- No external access

**NodePort:**
- Exposes on each node's IP
- Port range: 30000-32767
- Accessible from outside cluster

**LoadBalancer:**
- Creates cloud load balancer
- External IP assigned
- Best for production

### 73. What is kube-proxy and how does it work?

**Answer:**
kube-proxy is a network proxy that runs on each node. It:
- **Maintains rules**: Creates iptables/IPVS rules
- **Load balancing**: Distributes traffic to pods
- **Service discovery**: Enables service access
- **Modes**: iptables (default), IPVS, userspace

**How it works:**
1. Watches API server for Service/Endpoint changes
2. Updates iptables/IPVS rules
3. Routes traffic to backend pods
4. Load balances across endpoints

**Modes:**
- **iptables**: Default, uses iptables rules
- **IPVS**: Better performance, uses IPVS
- **userspace**: Legacy (deprecated)

### 74. What is a Service Mesh and why use it?

**Answer:**
A Service Mesh is a dedicated infrastructure layer for service-to-service communication. Benefits:
- **Traffic Management**: Load balancing, routing, retries
- **Security**: mTLS, authentication, authorization
- **Observability**: Metrics, tracing, logging
- **Policy**: Rate limiting, circuit breakers

**Popular Meshes:**
- **Istio**: Most feature-rich
- **Linkerd**: Lightweight, simple
- **Consul Connect**: HashiCorp solution
- **Kuma**: Universal mesh

**Components:**
- **Control Plane**: Manages configuration
- **Data Plane**: Sidecar proxies (Envoy)

**When to use:**
- Microservices architecture
- Need for advanced traffic management
- Security requirements (mTLS)
- Observability needs

### 75. What is the difference between Ingress and Ingress Controller?

**Answer:**

| Component | Purpose | Implementation |
|-----------|---------|----------------|
| **Ingress** | Defines routing rules | Kubernetes resource |
| **Ingress Controller** | Implements Ingress | Application (nginx, traefik) |

**Ingress:**
- Kubernetes resource
- Defines routing rules
- Declarative configuration
- Doesn't do actual routing

**Ingress Controller:**
- Application that implements Ingress
- Reads Ingress resources
- Performs actual routing
- Examples: nginx, traefik, Istio Gateway

**Example:**
```yaml
# Ingress (rules)
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
spec:
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        backend:
          service:
            name: app-service
            port:
              number: 80
```

---

## Storage & Data Management

### 76. What is the difference between static and dynamic volume provisioning?

**Answer:**

| Type | Provisioning | Management | Use Case |
|------|-------------|------------|----------|
| **Static** | Manual | Admin creates PVs | Predictable workloads |
| **Dynamic** | Automatic | StorageClass creates PVs | On-demand storage |

**Static Provisioning:**
- Admin creates PVs manually
- User creates PVC
- Kubernetes binds PVC to PV
- More control, manual work

**Dynamic Provisioning:**
- Admin creates StorageClass
- User creates PVC
- StorageClass automatically creates PV
- Less manual work, automated

**Example:**
```yaml
# StorageClass (dynamic)
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3

# PVC uses StorageClass
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  storageClassName: fast-ssd
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
```

### 77. What is a Volume Snapshot?

**Answer:**
Volume Snapshots provide the ability to create snapshots of persistent volumes. It:
- **Point-in-time**: Creates point-in-time copy
- **Backup**: Used for backups
- **Cloning**: Can clone volumes from snapshots
- **CSI**: Requires CSI driver with snapshot support

**Components:**
- **VolumeSnapshotClass**: Defines snapshot parameters
- **VolumeSnapshot**: Request for snapshot
- **VolumeSnapshotContent**: Actual snapshot

**Example:**
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotClass
metadata:
  name: csi-snapshotter
driver: ebs.csi.aws.com
deletionPolicy: Delete
---
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: my-snapshot
spec:
  volumeSnapshotClassName: csi-snapshotter
  source:
    persistentVolumeClaimName: my-pvc
```

### 78. What is the difference between EmptyDir and PersistentVolume?

**Answer:**

| Aspect | EmptyDir | PersistentVolume |
|--------|-----------|-----------------|
| **Lifecycle** | Pod lifecycle | Independent |
| **Persistence** | Ephemeral | Persistent |
| **Scope** | Pod | Cluster |
| **Use Case** | Temporary storage | Long-term storage |

**EmptyDir:**
- Created when pod starts
- Deleted when pod terminates
- Shared between containers in pod
- Use for temporary files, cache

**PersistentVolume:**
- Independent of pod lifecycle
- Data persists across pod restarts
- Can be shared or exclusive
- Use for databases, stateful apps

### 79. How do you manage secrets in Kubernetes?

**Answer:**

**Built-in Secrets:**
- **Opaque**: User-defined data
- **docker-registry**: Docker registry credentials
- **tls**: TLS certificates
- **service-account-token**: Service account tokens

**Best Practices:**
- Use external secret management (Vault, AWS Secrets Manager)
- Enable encryption at rest
- Rotate secrets regularly
- Use RBAC to limit access
- Never commit secrets to Git

**External Tools:**
- **HashiCorp Vault**: Secret management
- **Sealed Secrets**: Encrypted secrets in Git
- **External Secrets Operator**: Syncs external secrets
- **AWS Secrets Manager**: AWS managed

**Example:**
```bash
# Create secret
kubectl create secret generic my-secret \
  --from-literal=username=admin \
  --from-literal=password=secret

# Use in pod
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    env:
    - name: USERNAME
      valueFrom:
        secretKeyRef:
          name: my-secret
          key: username
```

### 80. What is a ConfigMap and when should you use it?

**Answer:**
ConfigMap stores non-sensitive configuration data. Use it for:
- **Configuration files**: Application config
- **Environment variables**: Non-sensitive env vars
- **Command-line arguments**: App arguments
- **Configuration data**: Any non-sensitive data

**When to use:**
- Application configuration
- Feature flags
- Environment-specific settings
- Non-sensitive data

**When NOT to use:**
- Sensitive data (use Secrets)
- Large files (>1MB)
- Binary data

**Example:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  database_url: "postgresql://localhost:5432/mydb"
  log_level: "info"
  config.yaml: |
    server:
      port: 8080
      host: 0.0.0.0
```

---

## Security & Compliance

### 81. What is Pod Security Policy (PSP) and why is it deprecated?

**Answer:**
Pod Security Policy was a cluster-level resource that controlled security-sensitive aspects of pod specification. It's deprecated because:
- **Complexity**: Difficult to use correctly
- **Confusion**: Multiple policies per pod
- **Replacement**: Pod Security Standards (simpler)

**Replacement:**
- **Pod Security Standards**: Simpler, namespace-level
- **Admission Controllers**: Custom validation
- **OPA/Gatekeeper**: Policy enforcement

**Pod Security Standards:**
- **Privileged**: No restrictions
- **Baseline**: Minimal restrictions
- **Restricted**: Maximum restrictions

### 82. What is NetworkPolicy and how does it work?

**Answer:**
NetworkPolicy controls traffic flow between pods. It:
- **Ingress rules**: Controls incoming traffic
- **Egress rules**: Controls outgoing traffic
- **Label-based**: Uses pod labels for selection
- **CNI requirement**: Requires CNI plugin support

**Requirements:**
- CNI plugin must support NetworkPolicy
- Flannel does NOT support
- Calico, Cilium, Weave support it

**Example:**
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
```

### 83. How do you implement RBAC in Kubernetes?

**Answer:**

**Components:**
- **Role/ClusterRole**: Defines permissions
- **RoleBinding/ClusterRoleBinding**: Grants permissions
- **Subjects**: Users, groups, service accounts

**Example:**
```yaml
# Role
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]

# RoleBinding
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

**Best Practices:**
- Principle of least privilege
- Use namespaced roles when possible
- Regular access reviews
- Document permissions

### 84. What is OPA (Open Policy Agent) and Gatekeeper?

**Answer:**
OPA is a general-purpose policy engine. Gatekeeper is a Kubernetes admission controller that uses OPA. They:
- **Policy enforcement**: Enforce policies on resources
- **Admission control**: Validate/mutate resources
- **Declarative**: Policies as code
- **Flexible**: Custom policy rules

**Use Cases:**
- Security policies
- Compliance requirements
- Resource constraints
- Best practices enforcement

**Example:**
```yaml
# Gatekeeper ConstraintTemplate
apiVersion: templates.gatekeeper.sh/v1beta1
kind: ConstraintTemplate
metadata:
  name: k8srequiredlabels
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredLabels
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8srequiredlabels
        violation[{"msg": msg}] {
          required := input.parameters.labels
          provided := input.review.object.metadata.labels
          missing := required - provided
          count(missing) > 0
          msg := sprintf("Missing required labels: %v", [missing])
        }
```

### 85. What is Pod Security Standards?

**Answer:**
Pod Security Standards define three policies for pod security:
- **Privileged**: No restrictions
- **Baseline**: Prevents known privilege escalations
- **Restricted**: Hardened security

**Enforcement Modes:**
- **enforce**: Rejects violating pods
- **audit**: Logs violations
- **warn**: Warns on violations

**Example:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## Deployment & Operations

### 86. What is the difference between Deployment and StatefulSet?

**Answer:**

| Feature | Deployment | StatefulSet |
|---------|-----------|-------------|
| **Identity** | No stable identity | Stable network identity |
| **Storage** | Shared volumes | Persistent volumes per pod |
| **Scaling** | Any order | Ordered (0, 1, 2...) |
| **DNS** | Service DNS | Stable pod DNS |
| **Use Case** | Stateless | Stateful |

**Deployment:**
- Stateless applications
- Random pod names
- Shared storage
- Any scaling order

**StatefulSet:**
- Stateful applications
- Ordered pod names (pod-0, pod-1)
- Individual persistent volumes
- Ordered scaling/deletion

### 87. What is a DaemonSet and when to use it?

**Answer:**
DaemonSet ensures a copy of a pod runs on all (or specific) nodes. Use it for:
- **Node-level services**: Logging, monitoring agents
- **Network plugins**: CNI plugins
- **Storage systems**: Storage daemons
- **System services**: kube-proxy

**Characteristics:**
- One pod per node (or subset)
- Pods created on new nodes automatically
- Pods deleted when nodes removed
- Useful for cluster-wide services

**Example:**
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
spec:
  selector:
    matchLabels:
      name: fluentd
  template:
    metadata:
      labels:
        name: fluentd
    spec:
      containers:
      - name: fluentd
        image: fluentd:latest
```

### 88. What is a Job and CronJob?

**Answer:**

**Job:**
- Creates pods that run to completion
- Ensures specified number of pods complete successfully
- Use for one-time tasks, batch processing

**CronJob:**
- Runs Jobs on a time-based schedule
- Like cron, but for Kubernetes
- Use for scheduled tasks, periodic jobs

**Example:**
```yaml
# Job
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4

# CronJob
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            command: ["echo", "Hello from CronJob"]
          restartPolicy: OnFailure
```

### 89. How do you perform rolling updates in Kubernetes?

**Answer:**

**Rolling Update Strategy:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # Can have 1 extra pod
      maxUnavailable: 0  # Must have all pods available
  template:
    spec:
      containers:
      - name: app
        image: nginx:1.21
```

**Process:**
1. Creates new ReplicaSet with new image
2. Gradually scales up new ReplicaSet
3. Scales down old ReplicaSet
4. Maintains service availability

**Commands:**
```bash
# Update image
kubectl set image deployment/my-app app=nginx:1.22

# Check rollout status
kubectl rollout status deployment/my-app

# Rollback
kubectl rollout undo deployment/my-app

# View rollout history
kubectl rollout history deployment/my-app
```

### 90. What is the difference between Recreate and RollingUpdate?

**Answer:**

| Strategy | Process | Downtime | Use Case |
|---------|---------|----------|----------|
| **Recreate** | Delete all, create new | Brief downtime | Stateful apps |
| **RollingUpdate** | Gradual replacement | Zero downtime | Stateless apps |

**Recreate:**
- Terminates all old pods
- Creates new pods
- Brief service interruption
- Use for stateful apps that can't run multiple versions

**RollingUpdate:**
- Gradually replaces pods
- Maintains service availability
- Zero downtime
- Default strategy

---

## Troubleshooting & Debugging

### 91. How do you debug a pod that keeps restarting?

**Answer:**

**Steps:**
1. **Check pod status:**
   ```bash
   kubectl get pods
   kubectl describe pod <pod-name>
   ```

2. **Check logs:**
   ```bash
   kubectl logs <pod-name>
   kubectl logs <pod-name> --previous
   ```

3. **Check events:**
   ```bash
   kubectl get events --sort-by='.lastTimestamp'
   ```

4. **Check resource limits:**
   ```bash
   kubectl describe pod <pod-name> | grep -A 5 "Limits:"
   ```

5. **Check liveness probe:**
   ```bash
   kubectl describe pod <pod-name> | grep -A 10 "Liveness:"
   ```

**Common Causes:**
- Application crashes
- Liveness probe failures
- Resource limits exceeded
- Configuration errors
- Image pull failures

### 92. How do you troubleshoot network connectivity issues?

**Answer:**

**Diagnosis:**
1. **Check service endpoints:**
   ```bash
   kubectl get endpoints <service-name>
   ```

2. **Test DNS resolution:**
   ```bash
   kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup <service-name>
   ```

3. **Test connectivity:**
   ```bash
   kubectl run -it --rm debug --image=busybox --restart=Never -- wget -O- <service-url>
   ```

4. **Check network policies:**
   ```bash
   kubectl get networkpolicies
   ```

5. **Check CNI pods:**
   ```bash
   kubectl get pods -n kube-system | grep -E "flannel|calico|cilium"
   ```

**Common Issues:**
- Service selector mismatch
- NetworkPolicy blocking traffic
- CNI plugin issues
- DNS resolution failures

### 93. How do you troubleshoot storage issues?

**Answer:**

**Steps:**
1. **Check PVC status:**
   ```bash
   kubectl get pvc
   kubectl describe pvc <pvc-name>
   ```

2. **Check PV status:**
   ```bash
   kubectl get pv
   kubectl describe pv <pv-name>
   ```

3. **Check StorageClass:**
   ```bash
   kubectl get storageclass
   kubectl describe storageclass <sc-name>
   ```

4. **Check pod volume mounts:**
   ```bash
   kubectl describe pod <pod-name> | grep -A 10 "Mounts:"
   ```

5. **Test volume access:**
   ```bash
   kubectl exec <pod-name> -- ls -la /mount/point
   ```

**Common Issues:**
- PVC not bound
- StorageClass not found
- Insufficient storage
- Permission issues
- Volume mount failures

### 94. How do you debug API server issues?

**Answer:**

**Diagnosis:**
1. **Check API server status:**
   ```bash
   kubectl get componentstatuses
   kubectl get nodes
   ```

2. **Check API server logs:**
   ```bash
   # On control plane node
   journalctl -u kube-apiserver -f
   ```

3. **Test API connectivity:**
   ```bash
   curl -k https://localhost:6443/healthz
   ```

4. **Check etcd:**
   ```bash
   ETCDCTL_API=3 etcdctl endpoint health
   ```

5. **Check certificates:**
   ```bash
   openssl x509 -in /etc/kubernetes/pki/apiserver.crt -text -noout
   ```

**Common Issues:**
- etcd connectivity
- Certificate expiration
- Resource exhaustion
- Network issues

### 95. How do you troubleshoot scheduler issues?

**Answer:**

**Diagnosis:**
1. **Check scheduler logs:**
   ```bash
   journalctl -u kube-scheduler -f
   ```

2. **Check unschedulable pods:**
   ```bash
   kubectl get pods --field-selector=status.phase=Pending
   kubectl describe pod <pending-pod>
   ```

3. **Check node resources:**
   ```bash
   kubectl describe node <node-name>
   kubectl top nodes
   ```

4. **Check taints:**
   ```bash
   kubectl describe node | grep Taints
   ```

5. **Check node conditions:**
   ```bash
   kubectl get nodes -o wide
   kubectl describe node <node-name>
   ```

**Common Issues:**
- Insufficient resources
- Node taints
- Node not ready
- Affinity/anti-affinity rules
- Pod resource requests too high

---

## Advanced Concepts

### 96. What is a Custom Resource Definition (CRD)?

**Answer:**
CRD extends Kubernetes API to add custom resources. It:
- **Extends API**: Adds new resource types
- **Custom controllers**: Managed by custom controllers
- **Operators**: Foundation for operators
- **Domain-specific**: Represents domain concepts

**Example:**
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: databases.example.com
spec:
  group: example.com
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              databaseName:
                type: string
  scope: Namespaced
  names:
    plural: databases
    singular: database
    kind: Database
```

### 97. What is a Kubernetes Operator?

**Answer:**
A Kubernetes Operator is a method of packaging, deploying, and managing Kubernetes applications. It:
- **Custom controllers**: Uses custom controllers
- **CRDs**: Defines Custom Resources
- **Automation**: Automates complex operations
- **Domain knowledge**: Encodes operational knowledge

**Example Operators:**
- Prometheus Operator
- Istio Operator
- Elasticsearch Operator
- PostgreSQL Operator

**Benefits:**
- Automates complex operations
- Encodes best practices
- Reduces operational burden
- Self-healing applications

### 98. What is Helm and how does it work?

**Answer:**
Helm is a package manager for Kubernetes. It:
- **Charts**: Packages Kubernetes applications
- **Releases**: Manages application instances
- **Templates**: Uses Go templates for configuration
- **Repositories**: Stores and shares charts

**Key Concepts:**
- **Chart**: Package of Kubernetes resources
- **Release**: Instance of chart deployed to cluster
- **Repository**: Collection of charts

**Example:**
```bash
# Install chart
helm install my-release stable/nginx

# Upgrade release
helm upgrade my-release stable/nginx

# List releases
helm list

# Uninstall
helm uninstall my-release
```

### 99. What is GitOps?

**Answer:**
GitOps is a methodology for managing infrastructure and applications using Git as the source of truth. It:
- **Git as source**: All changes in Git
- **Declarative**: Desired state in Git
- **Automated**: Automated sync and deployment
- **Observable**: Full audit trail

**Tools:**
- **ArgoCD**: GitOps continuous delivery
- **Flux**: GitOps toolkit
- **Jenkins X**: CI/CD with GitOps

**Benefits:**
- Version control
- Audit trail
- Rollback capability
- Collaboration

### 100. What is the difference between kubectl apply and kubectl create?

**Answer:**

| Aspect | kubectl apply | kubectl create |
|--------|---------------|----------------|
| **Mode** | Declarative | Imperative |
| **Idempotent** | Yes | No |
| **Manages** | Creates and updates | Only creates |
| **Use Case** | Configuration files | Quick creation |

**kubectl apply:**
- Declarative approach
- Creates or updates resources
- Idempotent
- Recommended for production

**kubectl create:**
- Imperative approach
- Only creates resources
- Fails if exists
- Quick for testing

**Example:**
```bash
# Apply (declarative)
kubectl apply -f deployment.yaml

# Create (imperative)
kubectl create deployment nginx --image=nginx
```

---

## Performance & Optimization

### 101. How do you optimize Kubernetes cluster performance?

**Answer:**

**Node Optimization:**
- Right-size nodes
- Use appropriate instance types
- Enable node autoscaling
- Optimize kernel parameters

**Pod Optimization:**
- Set appropriate resource requests/limits
- Use HPA for scaling
- Optimize images (multi-stage builds)
- Use readiness/liveness probes

**Network Optimization:**
- Use IPVS mode for kube-proxy
- Optimize CNI plugin
- Use service mesh for advanced routing
- Enable connection pooling

**Storage Optimization:**
- Use appropriate storage classes
- Enable volume snapshots
- Optimize I/O performance
- Use local storage when appropriate

### 102. How do you implement resource quotas effectively?

**Answer:**

**Best Practices:**
1. **Set appropriate limits:**
   ```yaml
   apiVersion: v1
   kind: ResourceQuota
   metadata:
     name: compute-quota
   spec:
     hard:
       requests.cpu: "10"
       requests.memory: 20Gi
       limits.cpu: "20"
       limits.memory: 40Gi
   ```

2. **Use LimitRange for defaults:**
   - Sets default requests/limits
   - Prevents resource exhaustion

3. **Monitor usage:**
   - Use metrics-server
   - Set up alerts
   - Regular reviews

4. **Adjust based on usage:**
   - Review actual usage
   - Adjust quotas accordingly
   - Balance between teams

### 103. What is Vertical Pod Autoscaler (VPA)?

**Answer:**
VPA automatically adjusts pod resource requests and limits based on historical usage. It:
- **Right-sizing**: Optimizes resource allocation
- **Learning**: Analyzes historical patterns
- **Modes**: Off, Initial, Auto, Recreate

**Modes:**
- **Off**: Only recommendations
- **Initial**: Sets resources at creation
- **Auto**: Updates resources (requires recreation)
- **Recreate**: Recreates pods with new resources

**Note:** VPA and HPA on CPU/memory cannot be used together.

### 104. How do you implement cluster autoscaling?

**Answer:**

**Cluster Autoscaler:**
- Automatically adjusts cluster size
- Adds nodes when pods can't be scheduled
- Removes nodes when underutilized
- Works with cloud providers

**Node Autoscaler:**
- Scales node groups
- Based on unschedulable pods
- Respects node group constraints

**Configuration:**
```yaml
# Cluster Autoscaler deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cluster-autoscaler
spec:
  template:
    spec:
      containers:
      - name: cluster-autoscaler
        image: k8s.gcr.io/autoscaling/cluster-autoscaler
        command:
        - ./cluster-autoscaler
        - --nodes=1:10:node-group-name
        - --cloud-provider=aws
```

### 105. What are the best practices for Kubernetes security?

**Answer:**

**Best Practices:**
1. **Enable RBAC**: Use Role-Based Access Control
2. **Pod Security Standards**: Enforce security policies
3. **Network Policies**: Restrict pod-to-pod communication
4. **Image Security**: Scan images, use trusted sources
5. **Secrets Management**: Use external secret management
6. **Encryption**: Enable encryption at rest and in transit
7. **Regular Updates**: Keep cluster components updated
8. **Audit Logging**: Enable audit logs
9. **Least Privilege**: Grant minimum required permissions
10. **Security Scanning**: Regular vulnerability scans

**Example:**
```yaml
# Pod Security Standards
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## Real-World Scenarios

### 106. How do you handle zero-downtime deployments in Kubernetes?

**Answer:**

**Strategies:**
1. **Rolling Updates:**
   ```yaml
   strategy:
     type: RollingUpdate
     rollingUpdate:
       maxSurge: 1
       maxUnavailable: 0
   ```

2. **Readiness Probes:**
   - Ensure pods are ready before receiving traffic
   - Prevents traffic to unready pods

3. **PreStop Hooks:**
   ```yaml
   lifecycle:
     preStop:
       exec:
         command: ["/bin/sh", "-c", "sleep 15"]
   ```

4. **Service Mesh:**
   - Gradual traffic shifting
   - Canary deployments
   - Traffic mirroring

**Best Practices:**
- Use rolling updates
- Configure proper probes
- Test in staging first
- Monitor during deployment
- Have rollback plan ready

### 107. How do you implement canary deployments in Kubernetes?

**Answer:**

**Method 1: Multiple Deployments**
```yaml
# Stable deployment (90% traffic)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-stable
spec:
  replicas: 9
---
# Canary deployment (10% traffic)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-canary
spec:
  replicas: 1
---
# Service with selector
apiVersion: v1
kind: Service
metadata:
  name: app-service
spec:
  selector:
    app: my-app
  ports:
  - port: 80
```

**Method 2: Using Service Mesh (Istio)**
```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: app
spec:
  hosts:
  - app.example.com
  http:
  - match:
    - headers:
        canary:
          exact: "true"
    route:
    - destination:
        host: app
        subset: canary
      weight: 100
  - route:
    - destination:
        host: app
        subset: stable
      weight: 90
    - destination:
        host: app
        subset: canary
      weight: 10
```

### 108. How do you manage configuration across multiple environments?

**Answer:**

**Methods:**
1. **Separate Namespaces:**
   ```bash
   kubectl create namespace dev
   kubectl create namespace staging
   kubectl create namespace prod
   ```

2. **ConfigMaps per Environment:**
   ```yaml
   # dev-configmap.yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: app-config
     namespace: dev
   data:
     env: "development"
   ```

3. **Helm with Values:**
   ```yaml
   # values-dev.yaml
   environment: development
   replicas: 1
   
   # values-prod.yaml
   environment: production
   replicas: 5
   ```

4. **Kustomize:**
   ```
   base/
   overlays/
     dev/
     staging/
     prod/
   ```

**Best Practices:**
- Use namespaces for isolation
- Version control configurations
- Use templating (Helm/Kustomize)
- Avoid hardcoding values
- Use external config management

### 109. How do you implement disaster recovery for Kubernetes?

**Answer:**

**Backup Strategy:**
1. **etcd Backup:**
   ```bash
   ETCDCTL_API=3 etcdctl snapshot save /backup/etcd-snapshot.db
   ```

2. **Resource Backup:**
   ```bash
   # Backup all resources
   kubectl get all --all-namespaces -o yaml > backup.yaml
   ```

3. **PV Backup:**
   - Use volume snapshots
   - External backup solutions
   - Cloud provider snapshots

**Recovery Process:**
1. Restore etcd from backup
2. Restore cluster configuration
3. Restore application resources
4. Restore persistent volumes
5. Verify cluster health

**Best Practices:**
- Regular automated backups
- Test restore procedures
- Document recovery process
- Off-site backup storage
- Define RTO/RPO

### 110. How do you implement multi-cluster Kubernetes?

**Answer:**

**Approaches:**
1. **Federation (Deprecated):**
   - Kubernetes Federation v1 (deprecated)
   - Not recommended

2. **Multi-Cluster Management:**
   - **Rancher**: Multi-cluster management
   - **Google Anthos**: Multi-cloud management
   - **Azure Arc**: Multi-cloud management

3. **GitOps:**
   - ArgoCD with multiple clusters
   - Flux with multi-cluster
   - Single source of truth

4. **Service Mesh:**
   - Istio multi-cluster
   - Linkerd multi-cluster
   - Cross-cluster communication

**Use Cases:**
- High availability
- Geographic distribution
- Multi-cloud deployments
- Disaster recovery
- Workload isolation

### 111. What is the difference between kubectl port-forward and kubectl proxy?

**Answer:**

| Aspect | kubectl port-forward | kubectl proxy |
|--------|---------------------|---------------|
| **Purpose** | Forward specific port | Proxy API server |
| **Port** | Specific port | Fixed port (8001) |
| **Use Case** | Access pod/service | Access API server |
| **Protocol** | TCP | HTTP |

**kubectl port-forward:**
- Forwards pod/service port to local machine
- Access specific service
- Use for debugging, testing

**kubectl proxy:**
- Proxies API server to local machine
- Access Kubernetes API
- Use for API exploration

**Example:**
```bash
# Port forward
kubectl port-forward pod/my-pod 8080:80

# Proxy
kubectl proxy
# Access API at http://localhost:8001
```

### 112. How do you implement pod disruption budgets effectively?

**Answer:**

**PDB Configuration:**
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: my-app
```

**Best Practices:**
1. **Set appropriate minAvailable:**
   - Consider application requirements
   - Balance availability and updates

2. **Use with Deployments:**
   - PDB works with Deployments
   - Ensures availability during updates

3. **Monitor PDB:**
   - Check PDB status
   - Monitor disruptions
   - Adjust as needed

**Use Cases:**
- High availability applications
- During node maintenance
- Cluster upgrades
- Rolling updates

### 113. How do you implement resource quotas per team?

**Answer:**

**Strategy:**
1. **Create namespace per team:**
   ```bash
   kubectl create namespace team-a
   kubectl create namespace team-b
   ```

2. **Apply ResourceQuota:**
   ```yaml
   apiVersion: v1
   kind: ResourceQuota
   metadata:
     name: team-a-quota
     namespace: team-a
   spec:
     hard:
       requests.cpu: "10"
       requests.memory: 20Gi
       limits.cpu: "20"
       limits.memory: 40Gi
       pods: "20"
   ```

3. **Apply RBAC:**
   - Limit team access to their namespace
   - Use RoleBinding per namespace

**Benefits:**
- Resource isolation
- Cost control
- Prevents resource exhaustion
- Team autonomy

### 114. How do you troubleshoot image pull errors?

**Answer:**

**Common Issues:**
1. **Private Registry:**
   ```bash
   # Create image pull secret
   kubectl create secret docker-registry regcred \
     --docker-server=<registry> \
     --docker-username=<username> \
     --docker-password=<password>
   
   # Use in pod
   spec:
     imagePullSecrets:
     - name: regcred
   ```

2. **Network Issues:**
   - Check network connectivity
   - Check DNS resolution
   - Check firewall rules

3. **Authentication:**
   - Verify credentials
   - Check secret exists
   - Verify secret is correct

4. **Image Not Found:**
   - Verify image exists
   - Check image tag
   - Verify registry URL

**Diagnosis:**
```bash
# Check pod events
kubectl describe pod <pod-name>

# Check image pull secret
kubectl get secret <secret-name> -o yaml

# Test image pull manually
docker pull <image>
```

### 115. How do you implement custom schedulers in Kubernetes?

**Answer:**

**Custom Scheduler:**
1. **Create scheduler binary:**
   - Implements scheduler interface
   - Custom scheduling logic

2. **Deploy as Deployment:**
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: custom-scheduler
   spec:
     template:
       spec:
         containers:
         - name: scheduler
           image: custom-scheduler:latest
           command:
           - /usr/local/bin/kube-scheduler
           - --config=/etc/kubernetes/scheduler-config.yaml
   ```

3. **Use in Pod:**
   ```yaml
   spec:
     schedulerName: custom-scheduler
   ```

**Use Cases:**
- Special scheduling requirements
- Workload-specific logic
- Multi-tenant scheduling
- Cost optimization

### 116. What is the difference between kubectl exec and kubectl attach?

**Answer:**

| Aspect | kubectl exec | kubectl attach |
|--------|--------------|----------------|
| **Process** | New process | Main process (PID 1) |
| **Sessions** | Multiple | Single |
| **Exit** | Doesn't affect pod | May stop pod |
| **Use Case** | Debugging, commands | Interactive apps |

**kubectl exec:**
- Creates new process
- Multiple sessions possible
- Exiting doesn't affect pod
- Use for debugging

**kubectl attach:**
- Attaches to main process
- Single session
- Exiting may stop pod
- Use for interactive apps

### 117. How do you implement pod autoscaling based on custom metrics?

**Answer:**

**Custom Metrics HPA:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: my-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metric:
        name: http_requests_per_second
      target:
        type: AverageValue
        averageValue: "100"
```

**Requirements:**
- Metrics adapter (Prometheus Adapter)
- Custom metrics API
- Metrics collection

**Example Metrics:**
- Request rate
- Queue length
- Custom business metrics
- Application-specific metrics

### 118. How do you implement pod priority and preemption?

**Answer:**

**Priority Classes:**
```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: high-priority
value: 1000
globalDefault: false
description: "High priority class"
---
apiVersion: v1
kind: Pod
metadata:
  name: important-pod
spec:
  priorityClassName: high-priority
  containers:
  - name: app
    image: nginx
```

**How it works:**
- Higher priority pods scheduled first
- Lower priority pods can be preempted
- Ensures critical workloads run

**Use Cases:**
- Critical applications
- System components
- Batch jobs with different priorities

### 119. How do you implement pod anti-affinity for high availability?

**Answer:**

**Pod Anti-Affinity:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3
  template:
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - my-app
            topologyKey: kubernetes.io/hostname
      containers:
      - name: app
        image: nginx
```

**Benefits:**
- Pods on different nodes
- High availability
- Fault tolerance
- Prevents single point of failure

### 120. How do you troubleshoot etcd performance issues?

**Answer:**

**Diagnosis:**
1. **Check etcd metrics:**
   ```bash
   ETCDCTL_API=3 etcdctl endpoint health
   ```

2. **Check etcd logs:**
   ```bash
   journalctl -u etcd -f
   ```

3. **Check disk I/O:**
   ```bash
   iostat -x 1
   ```

4. **Check etcd size:**
   ```bash
   ETCDCTL_API=3 etcdctl endpoint status
   ```

**Common Issues:**
- Disk I/O bottlenecks
- Large etcd database
- Network latency
- Resource constraints

**Solutions:**
- Use fast storage (SSD)
- Compact etcd regularly
- Defragment etcd
- Increase resources
- Optimize API server requests

---

## Advanced Networking & Security

### 121. How do you implement network policies for microservices?

**Answer:**

**Network Policy Strategy:**
```yaml
# Deny all by default
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress

# Allow frontend to backend
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
```

**Best Practices:**
- Deny all by default
- Allow specific traffic
- Use labels for selection
- Test policies thoroughly

### 122. How do you implement TLS/SSL in Kubernetes?

**Answer:**

**Methods:**
1. **Ingress TLS:**
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: tls-ingress
   spec:
     tls:
     - hosts:
       - app.example.com
       secretName: tls-secret
     rules:
     - host: app.example.com
       http:
         paths:
         - path: /
           backend:
             service:
               name: app-service
               port:
                 number: 80
   ```

2. **Cert-Manager:**
   - Automatic certificate management
   - Let's Encrypt integration
   - Certificate rotation

3. **Service Mesh:**
   - mTLS between services
   - Automatic certificate management
   - Istio, Linkerd

### 123. What is the difference between ClusterIP None and regular ClusterIP?

**Answer:**

| Aspect | ClusterIP (regular) | ClusterIP: None (Headless) |
|--------|-------------------|---------------------------|
| **IP** | Single cluster IP | No cluster IP |
| **DNS** | Returns service IP | Returns pod IPs |
| **Load Balancing** | Yes | No |
| **Use Case** | Regular services | StatefulSets, direct pod access |

**Regular ClusterIP:**
- Single IP address
- Load balances to pods
- Standard service behavior

**Headless Service (ClusterIP: None):**
- No cluster IP
- Returns individual pod IPs
- Direct pod access
- Used with StatefulSets

### 124. How do you implement service mesh in Kubernetes?

**Answer:**

**Istio Example:**
1. **Install Istio:**
   ```bash
   istioctl install
   ```

2. **Enable sidecar injection:**
   ```yaml
   apiVersion: v1
   kind: Namespace
   metadata:
     name: production
     labels:
       istio-injection: enabled
   ```

3. **Deploy application:**
   - Pods get sidecar automatically
   - mTLS enabled
   - Traffic management

**Benefits:**
- Traffic management
- Security (mTLS)
- Observability
- Policy enforcement

### 125. How do you troubleshoot DNS issues in Kubernetes?

**Answer:**

**Diagnosis:**
1. **Test DNS resolution:**
   ```bash
   kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup kubernetes.default
   ```

2. **Check CoreDNS pods:**
   ```bash
   kubectl get pods -n kube-system | grep coredns
   kubectl logs -n kube-system <coredns-pod>
   ```

3. **Check CoreDNS config:**
   ```bash
   kubectl get configmap coredns -n kube-system -o yaml
   ```

4. **Check service endpoints:**
   ```bash
   kubectl get endpoints
   ```

**Common Issues:**
- CoreDNS pods not running
- Incorrect DNS configuration
- Network policies blocking
- Service selector mismatch

---

## Storage & Data Management Advanced

### 126. How do you implement volume snapshots for backup?

**Answer:**

**Volume Snapshot:**
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotClass
metadata:
  name: csi-snapshotter
driver: ebs.csi.aws.com
deletionPolicy: Delete
---
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: my-snapshot
spec:
  volumeSnapshotClassName: csi-snapshotter
  source:
    persistentVolumeClaimName: my-pvc
```

**Restore from Snapshot:**
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: restored-pvc
spec:
  dataSource:
    name: my-snapshot
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

### 127. How do you implement dynamic volume provisioning?

**Answer:**

**StorageClass:**
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
  fsType: ext4
  encrypted: "true"
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
```

**PVC:**
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  storageClassName: fast-ssd
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
```

**Benefits:**
- Automatic provisioning
- No manual PV creation
- Flexible storage classes
- Cost optimization

### 128. What is the difference between StatefulSet and Deployment for databases?

**Answer:**

| Aspect | Deployment | StatefulSet |
|--------|-----------|-------------|
| **Identity** | No stable identity | Stable identity |
| **Storage** | Shared | Individual per pod |
| **Scaling** | Any order | Ordered |
| **DNS** | Service DNS | Stable pod DNS |
| **Use Case** | Stateless | Databases, stateful |

**StatefulSet for Databases:**
- Each pod has own storage
- Stable network identity
- Ordered scaling
- Predictable pod names
- Perfect for databases

**Example:**
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql
  replicas: 3
  template:
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 10Gi
```

### 129. How do you implement backup for StatefulSets?

**Answer:**

**Backup Strategy:**
1. **Volume Snapshots:**
   - Use CSI volume snapshots
   - Point-in-time backups
   - Automated via CronJob

2. **Application-level Backup:**
   ```yaml
   apiVersion: batch/v1
   kind: CronJob
   metadata:
     name: db-backup
   spec:
     schedule: "0 2 * * *"
     jobTemplate:
       spec:
         template:
           spec:
             containers:
             - name: backup
               image: mysql:8.0
               command:
               - mysqldump
               - -h mysql-0.mysql
               - --all-databases
               - > /backup/backup.sql
               volumeMounts:
               - name: backup-volume
                 mountPath: /backup
             volumes:
             - name: backup-volume
               persistentVolumeClaim:
                 claimName: backup-pvc
   ```

3. **External Backup Tools:**
   - Velero for cluster backups
   - Cloud provider snapshots
   - Custom backup solutions

### 130. How do you implement ReadWriteMany volumes?

**Answer:**

**StorageClass with RWX:**
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nfs-storage
provisioner: example.com/nfs
parameters:
  server: nfs-server.example.com
  path: /exports
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-pvc
spec:
  storageClassName: nfs-storage
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 100Gi
```

**Supported Storage:**
- NFS
- GlusterFS
- CephFS
- Azure File
- Some cloud storage

**Use Cases:**
- Shared configuration
- Content management
- Log aggregation
- Shared file systems

---

## Advanced Operations

### 131. How do you implement cluster upgrades?

**Answer:**

**Upgrade Process:**
1. **Backup:**
   - Backup etcd
   - Backup cluster configuration
   - Document current state

2. **Upgrade Control Plane:**
   ```bash
   # Upgrade kubeadm
   apt-mark unhold kubeadm
   apt-get update && apt-get install -y kubeadm=1.28.0-00
   apt-mark hold kubeadm
   
   # Plan upgrade
   kubeadm upgrade plan
   
   # Apply upgrade
   kubeadm upgrade apply v1.28.0
   ```

3. **Upgrade Nodes:**
   ```bash
   # Drain node
   kubectl drain <node-name> --ignore-daemonsets
   
   # Upgrade kubelet
   apt-get update && apt-get install -y kubelet=1.28.0-00
   
   # Restart kubelet
   systemctl daemon-reload
   systemctl restart kubelet
   
   # Uncordon node
   kubectl uncordon <node-name>
   ```

**Best Practices:**
- Test in non-production first
- Upgrade one version at a time
- Backup before upgrade
- Document upgrade process
- Have rollback plan

### 132. How do you implement node maintenance?

**Answer:**

**Process:**
1. **Cordon node:**
   ```bash
   kubectl cordon <node-name>
   ```

2. **Drain node:**
   ```bash
   kubectl drain <node-name> \
     --ignore-daemonsets \
     --delete-emptydir-data \
     --force
   ```

3. **Perform maintenance:**
   - Update OS
   - Update packages
   - Hardware maintenance

4. **Uncordon node:**
   ```bash
   kubectl uncordon <node-name>
   ```

**Best Practices:**
- Use during maintenance windows
- Ensure pod disruption budgets
- Monitor during drain
- Test in staging first

### 133. How do you implement cluster backup automation?

**Answer:**

**Velero:**
```bash
# Install Velero
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws \
  --bucket my-backup-bucket \
  --secret-file ./credentials-velero

# Create backup
velero backup create my-backup

# Schedule backups
velero schedule create daily-backup --schedule="0 2 * * *"
```

**Custom Script:**
```bash
#!/bin/bash
# Backup etcd
ETCDCTL_API=3 etcdctl snapshot save /backup/etcd-$(date +%Y%m%d).db

# Backup resources
kubectl get all --all-namespaces -o yaml > /backup/resources-$(date +%Y%m%d).yaml

# Upload to S3
aws s3 cp /backup s3://my-backup-bucket/ --recursive
```

### 134. How do you implement multi-region Kubernetes?

**Answer:**

**Approaches:**
1. **Separate Clusters:**
   - One cluster per region
   - Use federation or management tools
   - Independent operations

2. **Global Load Balancer:**
   - Route traffic to nearest region
   - DNS-based routing
   - Health checks

3. **Service Mesh:**
   - Cross-region communication
   - Traffic management
   - Service discovery

4. **GitOps:**
   - Single source of truth
   - Deploy to multiple clusters
   - Consistent configuration

**Considerations:**
- Data locality
- Latency
- Compliance
- Disaster recovery
- Cost

### 135. How do you troubleshoot API server performance issues?

**Answer:**

**Diagnosis:**
1. **Check API server metrics:**
   ```bash
   curl http://localhost:6443/metrics | grep apiserver
   ```

2. **Check etcd performance:**
   ```bash
   ETCDCTL_API=3 etcdctl endpoint status
   ```

3. **Check resource usage:**
   ```bash
   kubectl top nodes
   top
   ```

4. **Check API server logs:**
   ```bash
   journalctl -u kube-apiserver -f
   ```

**Optimization:**
- Increase API server resources
- Optimize etcd
- Use API server caching
- Limit watch operations
- Use pagination

---

*This is the fifth batch of 30 questions (106-135). More questions will be added in subsequent batches.*

