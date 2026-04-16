package response

import (
	"encoding/json"
	"net/http"
)

// ErrorBody is the standard error envelope for all REST error responses.
type ErrorBody struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// JSON writes a JSON response with the given status code and payload.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// NotFound writes a 404 JSON response.
func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, ErrorBody{Error: "not_found", Message: message})
}

// BadRequest writes a 400 JSON response.
func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, ErrorBody{Error: "bad_request", Message: message})
}

// InternalError writes a 500 JSON response. The underlying error is never
// exposed to the caller — log it before calling this helper.
func InternalError(w http.ResponseWriter) {
	JSON(w, http.StatusInternalServerError, ErrorBody{Error: "internal_error"})
}
