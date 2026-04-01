package handler

import "net/http"

// Submitter defines the interface for updating pack sizes.
type Submitter interface {
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

}