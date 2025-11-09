package cri

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/cri-api/pkg/apis/runtime/v1"

	"github.com/my-experiments/docker-cri/pkg/docker"
)

// ImageService implements the CRI ImageService
type ImageService struct {
	dockerClient *docker.Client
}

// NewImageService creates a new ImageService
func NewImageService(dockerClient *docker.Client) *ImageService {
	return &ImageService{
		dockerClient: dockerClient,
	}
}

// ListImages lists existing images
func (i *ImageService) ListImages(ctx context.Context, req *v1.ListImagesRequest) (*v1.ListImagesResponse, error) {
	logrus.Debugf("ListImages request: %+v", req)

	images, err := i.dockerClient.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list images: %v", err)
	}

	var result []*v1.Image
	for _, img := range images {
		// Convert Docker image to CRI image
		image := &v1.Image{
			Id:          img.ID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Size_:       uint64(img.Size),
		}

		// Apply filter if specified
		if req.Filter != nil {
			if req.Filter.Image != nil {
				filterMatch := false
				for _, tag := range img.RepoTags {
					if strings.Contains(tag, req.Filter.Image.Image) {
						filterMatch = true
						break
					}
				}
				if !filterMatch {
					continue
				}
			}
		}

		result = append(result, image)
	}

	return &v1.ListImagesResponse{
		Images: result,
	}, nil
}

// ImageStatus returns the status of the image
func (i *ImageService) ImageStatus(ctx context.Context, req *v1.ImageStatusRequest) (*v1.ImageStatusResponse, error) {
	logrus.Debugf("ImageStatus request: %+v", req)

	if req.Image == nil || req.Image.Image == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	imageInspect, _, err := i.dockerClient.ImageInspect(ctx, req.Image.Image)
	if err != nil {
		// Image not found is not an error in CRI - return empty status
		return &v1.ImageStatusResponse{
			Image: nil,
		}, nil
	}

	image := &v1.Image{
		Id:          imageInspect.ID,
		RepoTags:    imageInspect.RepoTags,
		RepoDigests: imageInspect.RepoDigests,
		Size_:       uint64(imageInspect.Size),
	}

	return &v1.ImageStatusResponse{
		Image: image,
	}, nil
}

// PullImage pulls an image with authentication config
func (i *ImageService) PullImage(ctx context.Context, req *v1.PullImageRequest) (*v1.PullImageResponse, error) {
	logrus.Infof("PullImage request: %+v", req)

	if req.Image == nil || req.Image.Image == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	imageRef := req.Image.Image

	// Handle authentication if provided
	options := image.PullOptions{}
	if req.Auth != nil {
		// Docker client will use credentials from config or environment
		// Full implementation would parse req.Auth and configure authentication
		logrus.Debugf("Authentication provided for image pull")
	}

	// Pull the image
	reader, err := i.dockerClient.ImagePull(ctx, imageRef, options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to pull image: %v", err)
	}
	defer reader.Close()

	// Read the pull output (Docker streams progress)
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read pull output: %v", err)
	}

	// Get the pulled image ID
	imageInspect, _, err := i.dockerClient.ImageInspect(ctx, imageRef)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to inspect pulled image: %v", err)
	}

	logrus.Infof("Successfully pulled image: %s (ID: %s)", imageRef, imageInspect.ID)

	return &v1.PullImageResponse{
		ImageRef: imageInspect.ID,
	}, nil
}

// RemoveImage removes the image
func (i *ImageService) RemoveImage(ctx context.Context, req *v1.RemoveImageRequest) (*v1.RemoveImageResponse, error) {
	logrus.Infof("RemoveImage request: %+v", req)

	if req.Image == nil || req.Image.Image == "" {
		return nil, status.Error(codes.InvalidArgument, "image is required")
	}

	_, err := i.dockerClient.ImageRemove(ctx, req.Image.Image, image.RemoveOptions{
		Force:         false,
		PruneChildren: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove image: %v", err)
	}

	return &v1.RemoveImageResponse{}, nil
}

// ImageFsInfo returns information of the filesystem that is used to store images
func (i *ImageService) ImageFsInfo(ctx context.Context, req *v1.ImageFsInfoRequest) (*v1.ImageFsInfoResponse, error) {
	logrus.Debugf("ImageFsInfo request: %+v", req)

	// Get Docker system info to determine filesystem usage
	// This is a simplified version - full implementation would query Docker's disk usage
	images, err := i.dockerClient.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list images: %v", err)
	}

	var totalSize uint64
	for _, img := range images {
		totalSize += uint64(img.Size)
	}

	// Return filesystem info
	// In a full implementation, this would query the actual filesystem
	return &v1.ImageFsInfoResponse{
		ImageFilesystems: []*v1.FilesystemUsage{
			{
				Timestamp: 0, // Would be set to current timestamp
				FsId: &v1.FilesystemIdentifier{
					Mountpoint: "/var/lib/docker",
				},
				UsedBytes: &v1.UInt64Value{
					Value: totalSize,
				},
				InodesUsed: &v1.UInt64Value{
					Value: 0, // Would be calculated from actual filesystem
				},
			},
		},
	}, nil
}

