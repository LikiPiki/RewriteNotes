package routes

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/likipiki/RewriteNotes/app/postgres"
)

type UserHandlers struct {
	Controller postgres.UserService
}

// Router - register all handler to chi router
func (handlers UserHandlers) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.GetAll)

	return r
}

func (handlers UserHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}
