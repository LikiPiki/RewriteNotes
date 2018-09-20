package routes

import (
	"context"
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

	return r
}

func (handlers UserHandlers) NoteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}
		user, err := handlers.Controller.Get(userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
