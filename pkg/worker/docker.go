package worker

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunContainer executes `docker run -d <image>` and returns the container ID.
func RunContainer(image string) (string, error) {
	if strings.TrimSpace(image) == "" {
		return "", fmt.Errorf("image name is required")
	}

	cmd := exec.Command("docker", "run", "-d", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker run failed: %s", strings.TrimSpace(string(output)))
	}

	containerID := strings.TrimSpace(string(output))
	return containerID, nil
}
