package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Client wraps the Docker API client
type Client struct {
	cli *client.Client
}

// NewClient creates a new Docker client
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Test the connection
	ctx := context.Background()
	_, err = cli.Ping(ctx)
	if err != nil {
		return nil, err
	}

	logrus.Info("Successfully connected to Docker daemon")

	return &Client{cli: cli}, nil
}

// Close closes the Docker client
func (c *Client) Close() error {
	return c.cli.Close()
}

// ContainerCreate creates a container
func (c *Client) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *container.NetworkingConfig, platform string, containerName string) (container.CreateResponse, error) {
	return c.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
}

// ContainerStart starts a container
func (c *Client) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	return c.cli.ContainerStart(ctx, containerID, options)
}

// ContainerStop stops a container
func (c *Client) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	return c.cli.ContainerStop(ctx, containerID, options)
}

// ContainerRemove removes a container
func (c *Client) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	return c.cli.ContainerRemove(ctx, containerID, options)
}

// ContainerInspect inspects a container
func (c *Client) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(ctx, containerID)
}

// ContainerList lists containers
func (c *Client) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	return c.cli.ContainerList(ctx, options)
}

// ContainerStats returns container statistics
func (c *Client) ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error) {
	return c.cli.ContainerStats(ctx, containerID, stream)
}

// ContainerLogs returns container logs
func (c *Client) ContainerLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error) {
	return c.cli.ContainerLogs(ctx, containerID, options)
}

// ImagePull pulls an image
func (c *Client) ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error) {
	return c.cli.ImagePull(ctx, refStr, options)
}

// ImageList lists images
func (c *Client) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return c.cli.ImageList(ctx, options)
}

// ImageRemove removes an image
func (c *Client) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	return c.cli.ImageRemove(ctx, imageID, options)
}

// ImageInspect inspects an image
func (c *Client) ImageInspect(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) {
	return c.cli.ImageInspectWithRaw(ctx, imageID)
}

// ImageTag tags an image
func (c *Client) ImageTag(ctx context.Context, source, target string) error {
	return c.cli.ImageTag(ctx, source, target)
}

