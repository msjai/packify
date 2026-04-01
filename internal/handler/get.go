package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Getter defines the interface for retrieving pack sizes.
type Getter interface {
	GetPackSizes() ([]int, error)
}

// PacksResponse is the JSON response for pack sizes endpoint.
type PacksResponse struct {
	Packs []int `json:"packs"`
}

// GetHandler handles HTTP requests for getting current pack sizes.
type GetHandler struct {
	name   string
	getter Getter
}

// NewGetHandler creates a new GetHandler with the given getter service.
func NewGetHandler(packGetter Getter) *GetHandler {
	return &GetHandler{
		name:   "packs get handler",
		getter: packGetter,
	}
}

// Handle processes a request to retrieve pack sizes.
func (h *GetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.getter.GetPackSizes()
	if err != nil {
		slog.Error("failed to get pack sizes", "handler", h.name, "error", err)
		http.Error(w, "failed to get pack sizes", http.StatusInternalServerError)
		return
	}

	raw, err := json.Marshal(PacksResponse{Packs: sizes})
	if err != nil {
		slog.Error("failed to marshal pack sizes", "handler", h.name, "error", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}
