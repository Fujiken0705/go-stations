package model

import "time"

type (
	// A TODO expresses a task with its metadata
	TODO struct {
		ID          int64     `json:"id"`
		Subject     string    `json:"subject"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// A CreateTODORequest expresses the request payload for creating a new TODO
	CreateTODORequest struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}

	// A CreateTODOResponse expresses the response payload after creating a TODO
	CreateTODOResponse struct {
		TODO *TODO `json:"todo"`
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct {
		PrevID int64 `json:"prev_id"`
		Size   int   `json:"size"`
	}

	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct {
		TODOs []TODO `json:"todos"`
	}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct {
		//11
		ID          int    `json:"id"`
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}

	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct {
		//11
		TODO TODO `json:"todo"`
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct {
		IDs []int64 `json:"ids"`
	}

	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct {
	}
)
