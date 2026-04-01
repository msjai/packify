package handler

import (
	"fmt"
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
		writeErrorResponse(w, h.name, fmt.Errorf("failed to get pack sizes: %w", err), http.StatusInternalServerError)
		return
	}

	writeSuccessResponse(w, PacksResponse{Packs: sizes})
}
