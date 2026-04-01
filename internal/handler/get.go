package handler

import (
	"net/http"
)

// Getter defines the interface for retrieving pack sizes.
type Getter interface {
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

}
