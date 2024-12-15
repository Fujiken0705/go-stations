package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	
	response := &model.HealthzResponse{
		Message: "OK",
	}

	encorder := json.NewEncoder(w)
	if err := encorder.Encode(response); err != nil {
			log.Println(err)
	}
}
