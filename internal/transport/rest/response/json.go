package response

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ErrorBody is the standard error envelope for all REST error responses.
// Code is an optional stable, machine-readable identifier (e.g. "order.already_paid").
type ErrorBody struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
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

// Unauthorized writes a 401 JSON response.
func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, ErrorBody{Error: "unauthorized", Message: message})
}

// Forbidden writes a 403 JSON response.
func Forbidden(w http.ResponseWriter, message string) {
	JSON(w, http.StatusForbidden, ErrorBody{Error: "forbidden", Message: message})
}

// TooManyRequests writes a 429 JSON response. A retryAfter > 0 (seconds) is
// added to the message and to the Retry-After header.
func TooManyRequests(w http.ResponseWriter, retryAfter int) {
	msg := "rate limit exceeded"
	if retryAfter > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
		msg += ", retry after " + strconv.Itoa(retryAfter) + "s"
	}
	JSON(w, http.StatusTooManyRequests, ErrorBody{Error: "too_many_requests", Message: msg})
}

// InternalError writes a 500 JSON response. The underlying error is never
// exposed to the caller — log it before calling this helper.
func InternalError(w http.ResponseWriter) {
	JSON(w, http.StatusInternalServerError, ErrorBody{Error: "internal_error"})
}
