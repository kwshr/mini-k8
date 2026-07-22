package models

// DeployRequest represents a user's request to deploy a container.
// It acts as a Command object — encapsulating the deployment intent
// as a serializable message that flows across process boundaries.
type DeployRequest struct {
	Image string `json:"image"`
}

// DeployResponse carries the result of a deployment attempt back to the caller.
type DeployResponse struct {
	Status      string `json:"status"`
	ContainerID string `json:"container_id,omitempty"`
	Error       string `json:"error,omitempty"`
}
