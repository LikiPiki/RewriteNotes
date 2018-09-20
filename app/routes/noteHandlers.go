package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/likipiki/RewriteNotes/app/postgres"
)

type NoteHandlers struct {
	Controller postgres.NoteService
}

func (handlers NoteHandlers) NoteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "note_id")
		if noteID == "" {
			next.ServeHTTP(w, r)
			return
		}
		note, err := handlers.Controller.Get(noteID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "note", note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Router - register all handler to chi router
func (handlers NoteHandlers) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(handlers.NoteCtx)

	r.Put("/{note_id}", handlers.Update)
	r.Delete("/{note_id}", handlers.Delete)
	r.Get("/{note_id}", handlers.Get)
	r.Post("/", handlers.Create)

	return r
}

func (handlers NoteHandlers) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get all"))
}

func (handlers NoteHandlers) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create"))
}

func (handlers NoteHandlers) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update"))
}

func (handlers NoteHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete"))
}
