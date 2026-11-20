package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainuser "github.com/binariado/go-service-template/internal/domain/user"
)

type fakeRepo struct {
	users map[string]*domainuser.User
	err   error
}

func (f fakeRepo) FindByID(_ context.Context, id string) (*domainuser.User, error) {
	return f.users[id], f.err
}

func (f fakeRepo) Create(_ context.Context, name string) (*domainuser.User, error) {
	return &domainuser.User{ID: "new", Name: name}, f.err
}

func TestCreateUser(t *testing.T) {
	uc := NewCreateUserUseCase(domainuser.NewUserService(fakeRepo{}))

	got, err := uc.Execute(context.Background(), "  Ada  ")
	if err != nil || got == nil || got.Name != "Ada" {
		t.Fatalf("got %+v, %v", got, err)
	}

	for _, name := range []string{"", "   ", strings.Repeat("a", 101)} {
		if _, err := uc.Execute(context.Background(), name); !errors.Is(err, domainuser.ErrInvalidName) {
			t.Fatalf("name %q: err = %v, want ErrInvalidName", name, err)
		}
	}
}

func TestGetUserByID(t *testing.T) {
	repo := fakeRepo{users: map[string]*domainuser.User{"1": {ID: "1", Name: "Ada"}}}
	uc := NewGetUserByIDUseCase(domainuser.NewUserService(repo))

	got, err := uc.Execute(context.Background(), "1")
	if err != nil || got == nil || got.Name != "Ada" {
		t.Fatalf("got %+v, %v", got, err)
	}

	got, err = uc.Execute(context.Background(), "missing")
	if err != nil || got != nil {
		t.Fatalf("expected nil, nil for unknown user; got %+v, %v", got, err)
	}
}

func TestGetUserByIDPropagatesErrors(t *testing.T) {
	boom := errors.New("boom")
	uc := NewGetUserByIDUseCase(domainuser.NewUserService(fakeRepo{err: boom}))

	if _, err := uc.Execute(context.Background(), "1"); !errors.Is(err, boom) {
		t.Fatalf("err = %v, want %v", err, boom)
	}
}
