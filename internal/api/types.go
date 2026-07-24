package api

type DeployRequest struct {
	Image string `json:"image"`
}

type DeployResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}