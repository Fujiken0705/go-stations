package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodPut:
		h.handleUpdate(w, r)
	case http.MethodGet:
		h.handleRead(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *TODOHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Subject is required
	if req.Subject == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdTodo, err := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.CreateTODOResponse{TODO: createdTodo})
}

func (h *TODOHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.Subject == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedTodo, err := h.svc.UpdateTODO(r.Context(), int64(req.ID), req.Subject, req.Description)
	if err != nil {
		if _, ok := err.(*model.ErrNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.UpdateTODOResponse{TODO: *updatedTodo})
}

func (h *TODOHandler) handleRead(w http.ResponseWriter, r *http.Request) {
	//URLのクエリパラメータを取得しTODORequestに値を代入
	query := r.URL.Query()
	prevID := query.Get("prev_id")
	size := query.Get("size")

	var req model.ReadTODORequest

	if prevID != "" {
		parsedPrevID, err := strconv.ParseInt(prevID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.PrevID = parsedPrevID
	}

	if size != "" {
		parsedSize, err := strconv.Atoi(size)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.Size = parsedSize
	} else {
		req.Size = 10
	}

	// ReadTODO メソッドを呼び出し
	todos, err := h.svc.ReadTODO(r.Context(), req.PrevID, int64(req.Size))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ReadTODOResponse を構築
	response := model.ReadTODOResponse{
		TODOs: []model.TODO{},
	}
	for _, todo := range todos {
		response.TODOs = append(response.TODOs, *todo)
	}

	// JSON Encode を行い HTTP Response を返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *TODOHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var req model.DeleteTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.svc.DeleteTODO(r.Context(), req.IDs)
	if err != nil {
		if _, ok := err.(*model.ErrNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.DeleteTODOResponse{})
}
