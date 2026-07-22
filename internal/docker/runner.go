package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

// Runner abstracts container execution so the Worker doesn't depend on
// a concrete Docker implementation. This is the Strategy pattern —
// tests can inject a FakeRunner, and future phases could swap in
// containerd or another runtime without touching the Worker.
type Runner interface {
	Run(image string) (containerID string, err error)
}

// ExecRunner implements Runner by shelling out to the Docker CLI.
type ExecRunner struct{}

// NewExecRunner creates a new ExecRunner.
func NewExecRunner() *ExecRunner {
	return &ExecRunner{}
}

// Run executes "docker run -d <image>" and returns the container ID.
func (r *ExecRunner) Run(image string) (string, error) {
	cmd := exec.Command("docker", "run", "-d", image)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker run failed: %s — %w", strings.TrimSpace(string(output)), err)
	}

	containerID := strings.TrimSpace(string(output))
	if containerID == "" {
		return "", fmt.Errorf("docker run returned empty container ID")
	}

	return containerID, nil
}
