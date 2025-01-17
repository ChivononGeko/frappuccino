package handler

import (
	"encoding/json"
	"errors"
	"frappuchino/internal/apperrors"
	"log/slog"
	"net/http"
	"strings"
)

func isJSONFile(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		slog.Error("Invalid content type: expected application/json")
		writeError(w, "Content type must be 'application/json'", http.StatusBadRequest)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func mapAppErrorToStatus(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrExistConflict):
		return http.StatusConflict
	case errors.Is(err, apperrors.ErrNotExistConflict):
		return http.StatusNotFound
	case errors.Is(err, apperrors.ErrOrderClosed):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
