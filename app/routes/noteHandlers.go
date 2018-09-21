package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/likipiki/RewriteNotes/app"

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
			http.Error(w, http.StatusText(404), 404)
			return
		}
		note, err := handlers.Controller.Get(noteID)
		fmt.Println("note", note, err)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "note", note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Router - register all handler to chi router
func (handlers NoteHandlers) Router() *chi.Mux {
	r := chi.NewRouter()
	sub := chi.NewRouter()
	r.Mount("/{note_id}", sub)
	sub.Use(handlers.NoteCtx)

	sub.Put("/", handlers.Update)
	sub.Delete("/", handlers.Delete)
	sub.Get("/", handlers.Get)

	r.Post("/{note_id}", handlers.Create)

	return r
}

func (handlers NoteHandlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	note, ok := ctx.Value("note").(app.Note)
	fmt.Println("ctx", ctx)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	json.NewEncoder(w).Encode(note)
}

func (handlers NoteHandlers) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create"))
}

func (handlers NoteHandlers) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update"))
}

func (handlers NoteHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "note_id")
	if noteID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	handlers.Controller.Delete(noteID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}
