package routes

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/likipiki/RewriteNotes/app/postgres"
)

type NoteHandlers struct {
	Controller postgres.UserService
}

// Router - register all handler to chi router
func (handlers NoteHandlers) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.GetAll)

	return r
}

func (handlers NoteHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}
