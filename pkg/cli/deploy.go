package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Deploy sends a deployment request to the control plane.
func Deploy(image, controlPlaneURL string) error {
	if strings.TrimSpace(image) == "" {
		return fmt.Errorf("image name is required")
	}

	// Ensure URL doesn't have trailing slash
	controlPlaneURL = strings.TrimSuffix(controlPlaneURL, "/")

	payload := map[string]string{"image": image}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		controlPlaneURL+"/deploy",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("failed to reach control plane: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		// If response is not JSON, use the raw body
		result = map[string]interface{}{"message": string(respBody)}
	}

	if resp.StatusCode != http.StatusOK {
		msg, ok := result["message"].(string)
		if !ok {
			msg = fmt.Sprintf("Status %d", resp.StatusCode)
		}
		return fmt.Errorf("deployment failed: %s", msg)
	}

	containerID, _ := result["containerId"].(string)
	if containerID == "" {
		containerID = image
	}

	fmt.Printf("Deployment successful: %s\n", containerID)
	return nil
}
