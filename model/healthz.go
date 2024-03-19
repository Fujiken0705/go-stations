package model

// HealthzResponse represents the response for the health check API.
type HealthzResponse struct {
	Message string `json:"message"`
}