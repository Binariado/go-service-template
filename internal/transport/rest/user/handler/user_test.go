package userhandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/binariado/go-service-template/internal/application"
	appuser "github.com/binariado/go-service-template/internal/application/user"
	domainuser "github.com/binariado/go-service-template/internal/domain/user"

	"github.com/go-chi/chi/v5"
)

const knownID = "6f1c1d7e-0b0a-4c63-9f64-0d3a2f1b7a10"

type fakeRepo struct{}

func (fakeRepo) FindByID(_ context.Context, id string) (*domainuser.User, error) {
	if id == knownID {
		return &domainuser.User{ID: id, Name: "Ada"}, nil
	}
	return nil, nil
}

func (fakeRepo) Create(_ context.Context, name string) (*domainuser.User, error) {
	return &domainuser.User{ID: knownID, Name: name}, nil
}

func TestCreate(t *testing.T) {
	app := &application.Application{User: appuser.NewApp(domainuser.NewUserService(fakeRepo{}))}
	r := chi.NewRouter()
	r.Post("/users", NewUserHandler(app).Create)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"created", `{"name":"Ada"}`, http.StatusCreated},
		{"empty name", `{"name":" "}`, http.StatusBadRequest},
		{"invalid json", `{`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body)))
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d (%s)", rec.Code, tt.want, rec.Body)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	app := &application.Application{User: appuser.NewApp(domainuser.NewUserService(fakeRepo{}))}
	r := chi.NewRouter()
	r.Get("/users/{id}", NewUserHandler(app).GetByID)

	tests := []struct {
		name string
		id   string
		want int
	}{
		{"found", knownID, http.StatusOK},
		{"not found", "11111111-1111-1111-1111-111111111111", http.StatusNotFound},
		{"invalid id", "abc", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/"+tt.id, nil))
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d (%s)", rec.Code, tt.want, rec.Body)
			}
		})
	}
}
