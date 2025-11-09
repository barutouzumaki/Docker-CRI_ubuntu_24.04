package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"k8s.io/cri-api/pkg/apis/runtime/v1"
)

// GenerateID generates a random ID
func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ConvertContainerState converts Docker container status to CRI PodSandboxState
func ConvertContainerState(status string) v1.PodSandboxState {
	statusLower := strings.ToLower(status)
	if strings.Contains(statusLower, "up") || strings.Contains(statusLower, "running") {
		return v1.PodSandboxState_SANDBOX_READY
	}
	return v1.PodSandboxState_SANDBOX_NOTREADY
}

