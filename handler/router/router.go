package router

import (
	"database/sql"
	"net/http"
	
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	
	// Register health check endpoint
	healthzHandler := handler.NewHealthzHandler()
	mux.Handle("/healthz", healthzHandler)

	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)
	mux.Handle("/todos", todoHandler)

	return mux
}