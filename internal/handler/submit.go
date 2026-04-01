package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// Submitter defines the interface for updating pack sizes.
type Submitter interface {
	UpdatePackSizes(sizes []int) error
}

// UpdatePacksRequest is the JSON request for updating pack sizes.
type UpdatePacksRequest struct {
	Packs []int `json:"packs"`
}

// Validate checks that pack sizes list is not empty, all values are positive, and no duplicates.
func (r UpdatePacksRequest) Validate() error {
	if len(r.Packs) == 0 {
		return errors.New("pack sizes list is empty")
	}

	seen := make(map[int]bool)
	for _, size := range r.Packs {
		if size <= 0 {
			return errors.New("pack size must be positive")
		}
		if seen[size] {
			return errors.New("duplicate pack size")
		}
		seen[size] = true
	}

	return nil
}

// SubmitHandler handles HTTP requests for updating pack sizes.
type SubmitHandler struct {
	name      string
	submitter Submitter
}

// NewSubmitHandler creates a new SubmitHandler with the given submitter service.
func NewSubmitHandler(packSubmitter Submitter) *SubmitHandler {
	return &SubmitHandler{
		name:      "pack submit handler",
		submitter: packSubmitter,
	}
}

// Handle processes a request to update pack sizes.
func (h *SubmitHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req UpdatePacksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", "handler", h.name, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		slog.Error("validation failed", "handler", h.name, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.submitter.UpdatePackSizes(req.Packs); err != nil {
		slog.Error("failed to update pack sizes", "handler", h.name, "error", err)
		http.Error(w, "failed to update pack sizes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}