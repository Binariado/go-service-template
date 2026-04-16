package handler

import (
	"backend-service/internal/domain/user"
	"backend-service/internal/transport/rest/dto"
	"backend-service/internal/transport/rest/response"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UserHandler handles REST requests for the user domain.
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler creates a UserHandler with the given service.
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetByID handles GET /api/v1/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	u, err := h.userService.FindByID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "user.FindByID failed", "id", id, "error", err)
		response.InternalError(w)
		return
	}
	if u == nil {
		response.NotFound(w, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, dto.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	})
}
