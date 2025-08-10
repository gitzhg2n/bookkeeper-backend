package routes

import (
	"encoding/json"
	"net/http"

	"bookkeeper-backend/middleware"
)

type ErrorResponse struct {
	Error         string `json:"error"`
	Message       string `json:"message"`
	Code          int    `json:"code"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

type SuccessResponse struct {
	Message       string      `json:"message"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationID string      `json:"correlation_id,omitempty"`
}

func writeJSONError(r *http.Request, w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	corrID, _ := middleware.RequestIDFromContext(r.Context())
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:         http.StatusText(code),
		Message:       message,
		Code:          code,
		CorrelationID: corrID,
	})
}

func writeJSONSuccess(r *http.Request, w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	corrID, _ := middleware.RequestIDFromContext(r.Context())
	_ = json.NewEncoder(w).Encode(SuccessResponse{
		Message:       message,
		Data:          data,
		CorrelationID: corrID,
	})
}