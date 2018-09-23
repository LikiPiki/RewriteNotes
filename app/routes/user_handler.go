package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/likipiki/RewriteNotes/app/postgres"
)

type UserHandlers struct {
	Controller postgres.UserService
	Router     *chi.Mux
}

func NewUserHandlers(controller postgres.UserService) UserHandlers {
	return UserHandlers{
		Controller: controller,
		Router:     UserHandlers{}.initRouter(),
	}
}

// Router - register all handler to chi router
func (h UserHandlers) initRouter() *chi.Mux {
	r := chi.NewRouter()

	return r
}

func (h UserHandlers) NoteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("id")
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}
		user, err := h.Controller.Get(userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
