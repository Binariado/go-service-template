package userhandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/binariado/go-service-template/internal/application"
	domainuser "github.com/binariado/go-service-template/internal/domain/user"
	"github.com/binariado/go-service-template/internal/transport/rest/response"
	userdto "github.com/binariado/go-service-template/internal/transport/rest/user/dto"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// UserHandler handles REST requests for the user module.
type UserHandler struct {
	app *application.Application
}

// NewUserHandler creates a UserHandler with the given application.
func NewUserHandler(app *application.Application) *UserHandler {
	return &UserHandler{app: app}
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req userdto.CreateUserRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}

	u, err := h.app.User.CreateUser.Execute(r.Context(), req.Name)
	if errors.Is(err, domainuser.ErrInvalidName) {
		response.BadRequest(w, err.Error())
		return
	}
	if err != nil {
		slog.ErrorContext(r.Context(), "user.CreateUser failed", "error", err)
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, userdto.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	})
}

// GetByID handles GET /api/v1/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		response.BadRequest(w, "id must be a valid UUID")
		return
	}

	u, err := h.app.User.GetUserByID.Execute(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "user.GetUserByID failed", "id", id, "error", err)
		response.InternalError(w)
		return
	}
	if u == nil {
		response.NotFound(w, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, userdto.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	})
}
