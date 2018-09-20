package routes

import (
	"github.com/go-chi/chi"
	"github.com/likipiki/RewriteNotes/app/postgres"
)

type UserHandlers struct {
	Controller postgres.UserService
}

// Router - register all handler to chi router
func (handlers UserHandlers) Router() *chi.Mux {
	r := chi.NewRouter()

	return r
}
