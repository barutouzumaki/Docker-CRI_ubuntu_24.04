package cri

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/cri-api/pkg/apis/runtime/v1"

	"github.com/my-experiments/docker-cri/pkg/docker"
	"github.com/my-experiments/docker-cri/pkg/utils"
)

// RuntimeService implements the CRI RuntimeService
type RuntimeService struct {
	dockerClient *docker.Client
	rootDir      string
}

// NewRuntimeService creates a new RuntimeService
func NewRuntimeService(dockerClient *docker.Client, rootDir string) *RuntimeService {
	return &RuntimeService{
		dockerClient: dockerClient,
		rootDir:      rootDir,
	}
}

// Version returns the runtime name, runtime name, runtime version and runtime API version
func (r *RuntimeService) Version(ctx context.Context, req *v1.VersionRequest) (*v1.VersionResponse, error) {
	logrus.Debugf("Version request: %+v", req)

	// Get Docker version info
	version, err := r.dockerClient.ContainerInspect(ctx, "")
	if err != nil {
		// If we can't inspect, we'll return default values
		logrus.Warnf("Could not get Docker version info: %v", err)
	}

	return &v1.VersionResponse{
		Version:           "0.1.0",
		RuntimeName:       "docker-cri",
		RuntimeVersion:    "1.0.0",
		RuntimeApiVersion: "v1",
	}, nil
}

// RunPodSandbox creates and starts a pod-level sandbox
func (r *RuntimeService) RunPodSandbox(ctx context.Context, req *v1.RunPodSandboxRequest) (*v1.RunPodSandboxResponse, error) {
	logrus.Infof("RunPodSandbox request: %+v", req)

	if req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "config is required")
	}

	// Create pause container for the sandbox
	sandboxID := req.Config.Metadata.Uid
	if sandboxID == "" {
		sandboxID = utils.GenerateID()
	}

	// Use pause image (typically gcr.io/google-containers/pause:3.9 or similar)
	pauseImage := "registry.k8s.io/pause:3.9"
	if req.Config.Image != nil && req.Config.Image.Image != "" {
		pauseImage = req.Config.Image.Image
	}

	// Create container config
	containerConfig := &container.Config{
		Image: pauseImage,
		Labels: map[string]string{
			"io.kubernetes.pod.uid":        sandboxID,
			"io.kubernetes.pod.name":       req.Config.Metadata.Name,
			"io.kubernetes.pod.namespace":  req.Config.Metadata.Namespace,
			"io.kubernetes.sandbox.id":     sandboxID,
			"io.kubernetes.container.type": "sandbox",
		},
		Entrypoint: []string{"/pause"},
	}

	// Create host config
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode("bridge"),
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	// Apply Linux-specific configurations
	if req.Config.Linux != nil {
		if req.Config.Linux.SecurityContext != nil {
			if req.Config.Linux.SecurityContext.NamespaceOptions != nil {
				if req.Config.Linux.SecurityContext.NamespaceOptions.Network == v1.NamespaceMode_NODE {
					hostConfig.NetworkMode = container.NetworkMode("host")
				}
			}
		}
	}

	// Create the container
	createResp, err := r.dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		"",
		fmt.Sprintf("k8s_POD_%s_%s_%s_0", req.Config.Metadata.Name, req.Config.Metadata.Namespace, sandboxID),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create sandbox container: %v", err)
	}

	// Start the container
	if err := r.dockerClient.ContainerStart(ctx, createResp.ID, container.StartOptions{}); err != nil {
		// Clean up on failure
		r.dockerClient.ContainerRemove(ctx, createResp.ID, container.RemoveOptions{Force: true})
		return nil, status.Errorf(codes.Internal, "failed to start sandbox container: %v", err)
	}

	logrus.Infof("Created and started sandbox: %s", createResp.ID)

	return &v1.RunPodSandboxResponse{
		PodSandboxId: createResp.ID,
	}, nil
}

// StopPodSandbox stops the sandbox
func (r *RuntimeService) StopPodSandbox(ctx context.Context, req *v1.StopPodSandboxRequest) (*v1.StopPodSandboxResponse, error) {
	logrus.Infof("StopPodSandbox request: %+v", req)

	if req.PodSandboxId == "" {
		return nil, status.Error(codes.InvalidArgument, "pod sandbox ID is required")
	}

	timeout := int(time.Second * 10)
	if err := r.dockerClient.ContainerStop(ctx, req.PodSandboxId, container.StopOptions{Timeout: &timeout}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to stop sandbox: %v", err)
	}

	return &v1.StopPodSandboxResponse{}, nil
}

// RemovePodSandbox removes the sandbox
func (r *RuntimeService) RemovePodSandbox(ctx context.Context, req *v1.RemovePodSandboxRequest) (*v1.RemovePodSandboxResponse, error) {
	logrus.Infof("RemovePodSandbox request: %+v", req)

	if req.PodSandboxId == "" {
		return nil, status.Error(codes.InvalidArgument, "pod sandbox ID is required")
	}

	if err := r.dockerClient.ContainerRemove(ctx, req.PodSandboxId, container.RemoveOptions{Force: true}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove sandbox: %v", err)
	}

	return &v1.RemovePodSandboxResponse{}, nil
}

// PodSandboxStatus returns the status of the PodSandbox
func (r *RuntimeService) PodSandboxStatus(ctx context.Context, req *v1.PodSandboxStatusRequest) (*v1.PodSandboxStatusResponse, error) {
	logrus.Debugf("PodSandboxStatus request: %+v", req)

	if req.PodSandboxId == "" {
		return nil, status.Error(codes.InvalidArgument, "pod sandbox ID is required")
	}

	containerJSON, err := r.dockerClient.ContainerInspect(ctx, req.PodSandboxId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sandbox not found: %v", err)
	}

	// Convert Docker container state to CRI sandbox state
	state := v1.PodSandboxState_SANDBOX_NOTREADY
	if containerJSON.State.Running {
		state = v1.PodSandboxState_SANDBOX_READY
	}

	status := &v1.PodSandboxStatus{
		Id:       req.PodSandboxId,
		State:    state,
		CreatedAt: containerJSON.Created.UnixNano(),
		Network: &v1.PodSandboxNetworkStatus{
			Ip: containerJSON.NetworkSettings.IPAddress,
		},
	}

	// Extract metadata from labels
	if containerJSON.Config != nil && containerJSON.Config.Labels != nil {
		labels := containerJSON.Config.Labels
		status.Metadata = &v1.PodSandboxMetadata{
			Name:      labels["io.kubernetes.pod.name"],
			Uid:       labels["io.kubernetes.pod.uid"],
			Namespace: labels["io.kubernetes.pod.namespace"],
		}
	}

	return &v1.PodSandboxStatusResponse{
		Status: status,
	}, nil
}

// ListPodSandbox lists PodSandboxes
func (r *RuntimeService) ListPodSandbox(ctx context.Context, req *v1.ListPodSandboxRequest) (*v1.ListPodSandboxResponse, error) {
	logrus.Debugf("ListPodSandbox request: %+v", req)

	containers, err := r.dockerClient.ContainerList(ctx, container.ListOptions{
		All: true,
		Filters: map[string][]string{
			"label": {"io.kubernetes.container.type=sandbox"},
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list containers: %v", err)
	}

	var sandboxes []*v1.PodSandbox
	for _, c := range containers {
		sandbox := &v1.PodSandbox{
			Id:       c.ID,
			State:    utils.ConvertContainerState(c.Status),
			CreatedAt: c.Created,
		}

		// Extract metadata from labels
		if labels := c.Labels; labels != nil {
			sandbox.Metadata = &v1.PodSandboxMetadata{
				Name:      labels["io.kubernetes.pod.name"],
				Uid:       labels["io.kubernetes.pod.uid"],
				Namespace: labels["io.kubernetes.pod.namespace"],
			}
		}

		sandboxes = append(sandboxes, sandbox)
	}

	return &v1.ListPodSandboxResponse{
		Items: sandboxes,
	}, nil
}

// CreateContainer creates a new container in the specified PodSandbox
func (r *RuntimeService) CreateContainer(ctx context.Context, req *v1.CreateContainerRequest) (*v1.CreateContainerResponse, error) {
	logrus.Infof("CreateContainer request: %+v", req)

	if req.PodSandboxId == "" {
		return nil, status.Error(codes.InvalidArgument, "pod sandbox ID is required")
	}
	if req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "config is required")
	}

	// Get sandbox info
	sandboxInfo, err := r.dockerClient.ContainerInspect(ctx, req.PodSandboxId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sandbox not found: %v", err)
	}

	// Create container config
	containerConfig := &container.Config{
		Image: req.Config.Image.Image,
		Labels: map[string]string{
			"io.kubernetes.pod.uid":        sandboxInfo.Config.Labels["io.kubernetes.pod.uid"],
			"io.kubernetes.pod.name":       sandboxInfo.Config.Labels["io.kubernetes.pod.name"],
			"io.kubernetes.pod.namespace":  sandboxInfo.Config.Labels["io.kubernetes.pod.namespace"],
			"io.kubernetes.container.name": req.Config.Metadata.Name,
			"io.kubernetes.container.type": "container",
		},
	}

	// Set command and args
	if req.Config.Command != nil {
		containerConfig.Cmd = req.Config.Command
	}
	if req.Config.Args != nil {
		containerConfig.Cmd = append(containerConfig.Cmd, req.Config.Args...)
	}

	// Set environment variables
	if req.Config.Envs != nil {
		for _, env := range req.Config.Envs {
			containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%s=%s", env.Key, env.Value))
		}
	}

	// Create host config
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(fmt.Sprintf("container:%s", req.PodSandboxId)),
	}

	// Apply resource limits
	if req.Config.Linux != nil && req.Config.Linux.Resources != nil {
		resources := req.Config.Linux.Resources
		if resources.Cpu != nil {
			if resources.Cpu.Shares != nil {
				hostConfig.CPUShares = int64(*resources.Cpu.Shares)
			}
			if resources.Cpu.Quota != nil {
				hostConfig.CPUQuota = *resources.Cpu.Quota
			}
			if resources.Cpu.Period != nil {
				hostConfig.CPUPeriod = int64(*resources.Cpu.Period)
			}
		}
		if resources.Memory != nil {
			if resources.Memory.Limit != nil {
				hostConfig.Memory = int64(*resources.Memory.Limit)
			}
		}
	}

	containerName := fmt.Sprintf("k8s_%s_%s_%s_%s",
		req.Config.Metadata.Name,
		sandboxInfo.Config.Labels["io.kubernetes.pod.name"],
		sandboxInfo.Config.Labels["io.kubernetes.pod.namespace"],
		req.Config.Metadata.Uid)

	// Create the container
	createResp, err := r.dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		"",
		containerName,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create container: %v", err)
	}

	logrus.Infof("Created container: %s", createResp.ID)

	return &v1.CreateContainerResponse{
		ContainerId: createResp.ID,
	}, nil
}

// StartContainer starts the container
func (r *RuntimeService) StartContainer(ctx context.Context, req *v1.StartContainerRequest) (*v1.StartContainerResponse, error) {
	logrus.Infof("StartContainer request: %+v", req)

	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	if err := r.dockerClient.ContainerStart(ctx, req.ContainerId, container.StartOptions{}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start container: %v", err)
	}

	return &v1.StartContainerResponse{}, nil
}

// StopContainer stops a running container
func (r *RuntimeService) StopContainer(ctx context.Context, req *v1.StopContainerRequest) (*v1.StopContainerResponse, error) {
	logrus.Infof("StopContainer request: %+v", req)

	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	timeout := int(time.Second * 10)
	if req.Timeout != 0 {
		timeout = int(time.Duration(req.Timeout) * time.Second)
	}

	if err := r.dockerClient.ContainerStop(ctx, req.ContainerId, container.StopOptions{Timeout: &timeout}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to stop container: %v", err)
	}

	return &v1.StopContainerResponse{}, nil
}

// RemoveContainer removes the container
func (r *RuntimeService) RemoveContainer(ctx context.Context, req *v1.RemoveContainerRequest) (*v1.RemoveContainerResponse, error) {
	logrus.Infof("RemoveContainer request: %+v", req)

	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	if err := r.dockerClient.ContainerRemove(ctx, req.ContainerId, container.RemoveOptions{Force: true}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove container: %v", err)
	}

	return &v1.RemoveContainerResponse{}, nil
}

// ContainerStatus returns the status of the container
func (r *RuntimeService) ContainerStatus(ctx context.Context, req *v1.ContainerStatusRequest) (*v1.ContainerStatusResponse, error) {
	logrus.Debugf("ContainerStatus request: %+v", req)

	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	containerJSON, err := r.dockerClient.ContainerInspect(ctx, req.ContainerId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "container not found: %v", err)
	}

	// Convert Docker state to CRI state
	state := v1.ContainerState_CONTAINER_UNKNOWN
	if containerJSON.State.Running {
		state = v1.ContainerState_CONTAINER_RUNNING
	} else if containerJSON.State.Status == "exited" {
		state = v1.ContainerState_CONTAINER_EXITED
	}

	status := &v1.ContainerStatus{
		Id:          req.ContainerId,
		State:       state,
		CreatedAt:   containerJSON.Created.UnixNano(),
		StartedAt:   containerJSON.State.StartedAt.UnixNano(),
		FinishedAt:  containerJSON.State.FinishedAt.UnixNano(),
		ExitCode:    int32(containerJSON.State.ExitCode),
		Image:       &v1.ImageSpec{Image: containerJSON.Config.Image},
		ImageRef:    containerJSON.Image,
	}

	// Extract metadata from labels
	if containerJSON.Config != nil && containerJSON.Config.Labels != nil {
		labels := containerJSON.Config.Labels
		status.Metadata = &v1.ContainerMetadata{
			Name:    labels["io.kubernetes.container.name"],
			Attempt: 0, // Docker doesn't track attempts
		}
	}

	return &v1.ContainerStatusResponse{
		Status: status,
	}, nil
}

// ListContainers lists all containers
func (r *RuntimeService) ListContainers(ctx context.Context, req *v1.ListContainersRequest) (*v1.ListContainersResponse, error) {
	logrus.Debugf("ListContainers request: %+v", req)

	containers, err := r.dockerClient.ContainerList(ctx, container.ListOptions{
		All: true,
		Filters: map[string][]string{
			"label": {"io.kubernetes.container.type=container"},
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list containers: %v", err)
	}

	var result []*v1.Container
	for _, c := range containers {
		container := &v1.Container{
			Id:       c.ID,
			State:    utils.ConvertContainerState(c.Status),
			CreatedAt: c.Created,
		}

		// Extract metadata from labels
		if labels := c.Labels; labels != nil {
			container.Metadata = &v1.ContainerMetadata{
				Name:      labels["io.kubernetes.container.name"],
				Attempt:   0,
			}
			container.Image = &v1.ImageSpec{Image: c.Image}
			container.ImageRef = c.ImageID
		}

		result = append(result, container)
	}

	return &v1.ListContainersResponse{
		Containers: result,
	}, nil
}

// ContainerStats returns stats of the container
func (r *RuntimeService) ContainerStats(ctx context.Context, req *v1.ContainerStatsRequest) (*v1.ContainerStatsResponse, error) {
	logrus.Debugf("ContainerStats request: %+v", req)

	if req.ContainerId == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID is required")
	}

	stats, err := r.dockerClient.ContainerStats(ctx, req.ContainerId, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get container stats: %v", err)
	}
	defer stats.Body.Close()

	// Parse stats and convert to CRI format
	// This is a simplified version - full implementation would parse the JSON
	containerStats := &v1.ContainerStats{
		Attributes: &v1.ContainerAttributes{
			Id: req.ContainerId,
		},
		// CPU and memory stats would be populated from the Docker stats JSON
	}

	return &v1.ContainerStatsResponse{
		Stats: containerStats,
	}, nil
}

// ListContainerStats lists stats for all containers
func (r *RuntimeService) ListContainerStats(ctx context.Context, req *v1.ListContainerStatsRequest) (*v1.ListContainerStatsResponse, error) {
	logrus.Debugf("ListContainerStats request: %+v", req)

	containers, err := r.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list containers: %v", err)
	}

	var stats []*v1.ContainerStats
	for _, c := range containers {
		stat, err := r.ContainerStats(ctx, &v1.ContainerStatsRequest{ContainerId: c.ID})
		if err != nil {
			logrus.Warnf("Failed to get stats for container %s: %v", c.ID, err)
			continue
		}
		stats = append(stats, stat.Stats)
	}

	return &v1.ListContainerStatsResponse{
		Stats: stats,
	}, nil
}

// UpdateContainerResources updates ContainerConfig of the container
func (r *RuntimeService) UpdateContainerResources(ctx context.Context, req *v1.UpdateContainerResourcesRequest) (*v1.UpdateContainerResourcesResponse, error) {
	logrus.Infof("UpdateContainerResources request: %+v", req)

	// Docker doesn't support updating resources on running containers directly
	// This would require container recreation in a full implementation
	return &v1.UpdateContainerResourcesResponse{}, nil
}

// ReopenContainerLog reopens the container log file
func (r *RuntimeService) ReopenContainerLog(ctx context.Context, req *v1.ReopenContainerLogRequest) (*v1.ReopenContainerLogResponse, error) {
	logrus.Debugf("ReopenContainerLog request: %+v", req)

	// Docker handles log rotation automatically
	return &v1.ReopenContainerLogResponse{}, nil
}

// ExecSync runs a command in a container synchronously
func (r *RuntimeService) ExecSync(ctx context.Context, req *v1.ExecSyncRequest) (*v1.ExecSyncResponse, error) {
	logrus.Infof("ExecSync request: %+v", req)

	// This would require implementing exec functionality
	// For now, return not implemented
	return nil, status.Error(codes.Unimplemented, "ExecSync not yet implemented")
}

// Exec prepares a streaming endpoint to execute a command in the container
func (r *RuntimeService) Exec(req *v1.ExecRequest, stream v1.RuntimeService_ExecServer) error {
	logrus.Infof("Exec request: %+v", req)

	// This would require implementing streaming exec functionality
	return status.Error(codes.Unimplemented, "Exec not yet implemented")
}

// Attach prepares a streaming endpoint to attach to a running container
func (r *RuntimeService) Attach(req *v1.AttachRequest, stream v1.RuntimeService_AttachServer) error {
	logrus.Infof("Attach request: %+v", req)

	// This would require implementing streaming attach functionality
	return status.Error(codes.Unimplemented, "Attach not yet implemented")
}

// PortForward prepares a streaming endpoint to forward ports from a PodSandbox
func (r *RuntimeService) PortForward(req *v1.PortForwardRequest, stream v1.RuntimeService_PortForwardServer) error {
	logrus.Infof("PortForward request: %+v", req)

	// This would require implementing port forwarding functionality
	return status.Error(codes.Unimplemented, "PortForward not yet implemented")
}

// UpdateRuntimeConfig updates the runtime configuration
func (r *RuntimeService) UpdateRuntimeConfig(ctx context.Context, req *v1.UpdateRuntimeConfigRequest) (*v1.UpdateRuntimeConfigResponse, error) {
	logrus.Debugf("UpdateRuntimeConfig request: %+v", req)

	return &v1.UpdateRuntimeConfigResponse{}, nil
}

// Status returns the status of the runtime
func (r *RuntimeService) Status(ctx context.Context, req *v1.StatusRequest) (*v1.StatusResponse, error) {
	logrus.Debugf("Status request: %+v", req)

	// Get container count
	containers, err := r.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list containers: %v", err)
	}

	conditions := []*v1.RuntimeCondition{
		{
			Type:   v1.RuntimeReady,
			Status: true,
		},
		{
			Type:   v1.NetworkReady,
			Status: true,
		},
	}

	return &v1.StatusResponse{
		Status: &v1.RuntimeStatus{
			Conditions: conditions,
		},
		Info: map[string]string{
			"containers": fmt.Sprintf("%d", len(containers)),
		},
	}, nil
}

// CheckpointContainer checkpoints a container
func (r *RuntimeService) CheckpointContainer(ctx context.Context, req *v1.CheckpointContainerRequest) (*v1.CheckpointContainerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CheckpointContainer not implemented")
}

// GetContainerEvents returns stream of container events
func (r *RuntimeService) GetContainerEvents(req *v1.GetContainerEventsRequest, stream v1.RuntimeService_GetContainerEventsServer) error {
	return status.Error(codes.Unimplemented, "GetContainerEvents not implemented")
}

// ListMetricDescriptors lists the descriptors for the metrics
func (r *RuntimeService) ListMetricDescriptors(ctx context.Context, req *v1.ListMetricDescriptorsRequest) (*v1.ListMetricDescriptorsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListMetricDescriptors not implemented")
}

// ListPodSandboxMetrics lists pod sandbox metrics
func (r *RuntimeService) ListPodSandboxMetrics(ctx context.Context, req *v1.ListPodSandboxMetricsRequest) (*v1.ListPodSandboxMetricsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListPodSandboxMetrics not implemented")
}

