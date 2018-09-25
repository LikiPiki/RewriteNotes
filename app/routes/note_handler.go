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
	controller postgres.NoteService
	Router     *chi.Mux
}

func NewNoteHandlers(contoller postgres.NoteService) NoteHandlers {
	return NoteHandlers{
		controller: contoller,
		Router:     NoteHandlers{}.initRouter(),
	}
}

func (h NoteHandlers) NoteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "note_id")
		if noteID == "" {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		note, err := h.controller.Get(noteID)
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
func (h NoteHandlers) initRouter() *chi.Mux {
	r := chi.NewRouter()
	sub := chi.NewRouter()
	r.Mount("/{note_id}", sub)
	sub.Use(h.NoteCtx)

	sub.Put("/", h.Update)
	sub.Delete("/", h.Delete)
	sub.Get("/", h.Get)

	r.Post("/{note_id}", h.Create)

	return r
}

func (h NoteHandlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	note, ok := ctx.Value("note").(app.Note)
	fmt.Println("ctx", ctx)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	json.NewEncoder(w).Encode(note)
}

func (h NoteHandlers) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create"))
}

func (h NoteHandlers) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update"))
}

func (h NoteHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "note_id")
	if noteID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	h.controller.Delete(noteID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}
