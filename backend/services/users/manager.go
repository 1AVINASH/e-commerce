package user

import (
	"github.com/go-chi/chi/v5"
)

type UserManager struct {
	repo *UserRepository
	apis *UserAPIs
}

func NewUserManager(mux *chi.Mux) *UserManager {
	repo := NewUserRepository()
	apis := NewUserApis(mux, repo)
	return &UserManager{
		repo: repo,
		apis: apis,
	}
}
