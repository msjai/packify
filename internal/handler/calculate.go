package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
)

// Calculator defines the interface for pack calculation logic.
type Calculator interface {
	Calculate(order int) (map[int]int, error)
}

// CalculateRequest is the JSON request for pack calculation.
type CalculateRequest struct {
	Order int `json:"order"`
}

// Validate checks that order quantity is positive.
func (r CalculateRequest) Validate() error {
	if r.Order <= 0 {
		return errors.New("order must be positive")
	}
	return nil
}

// PackEntry represents a single pack size and its quantity in the result.
type PackEntry struct {
	Pack     int `json:"pack"`
	Quantity int `json:"quantity"`
}

// CalculateResponse is the JSON response for pack calculation.
type CalculateResponse struct {
	Packs []PackEntry `json:"packs"`
}

// CalculateHandler handles HTTP requests for pack calculation.
type CalculateHandler struct {
	name       string
	calculator Calculator
}

// NewCalculateHandler creates a new CalculateHandler with the given calculator service.
func NewCalculateHandler(packCalculator Calculator) *CalculateHandler {
	return &CalculateHandler{
		name:       "pack calculate handler",
		calculator: packCalculator,
	}
}

// Handle processes a pack calculation request.
func (h *CalculateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := parseRequest(w, r, &req, h.name); err != nil {
		return
	}

	result, err := h.calculator.Calculate(req.Order)
	if err != nil {
		writeErrorResponse(w, h.name, fmt.Errorf("failed to calculate packs: %w", err), http.StatusInternalServerError)
		return
	}

	packs := make([]PackEntry, 0, len(result))
	for size, qty := range result {
		packs = append(packs, PackEntry{Pack: size, Quantity: qty})
	}

	sort.Slice(packs, func(i, j int) bool {
		return packs[i].Pack < packs[j].Pack
	})
	
	writeSuccessResponse(w, CalculateResponse{Packs: packs})
}