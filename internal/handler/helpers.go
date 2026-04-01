package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const applicationJSON = "application/json"

// RequestParser defines the interface for request validation.
type RequestParser interface {
	Validate() error
}

// writeErrorResponse logs the error and sends a plain text error response to the client.
func writeErrorResponse(w http.ResponseWriter, handlerName string, err error, statusCode int) {
	slog.Error(err.Error(), "handler", handlerName)
	http.Error(w, err.Error(), statusCode)
}

// writeSuccessResponse marshals data to JSON and writes it to the response.
func writeSuccessResponse(w http.ResponseWriter, data any) {
	raw, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", applicationJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}

// parseRequest decodes JSON body and validates the request.
func parseRequest(w http.ResponseWriter, r *http.Request, req RequestParser, handlerName string) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, handlerName, fmt.Errorf("decode request: %w", err), http.StatusBadRequest)
		return err
	}

	if err := req.Validate(); err != nil {
		writeErrorResponse(w, handlerName, err, http.StatusBadRequest)
		return err
	}

	return nil
}