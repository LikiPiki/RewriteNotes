package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/likipiki/RewriteNotes/app/postgres"
	"github.com/pkg/errors"
)

type UserHandlers struct {
	Controller postgres.UserService
	Router     *chi.Mux
}

type User struct {
	Id       uint   `json:"id"`
	Password string `json:"password"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
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

func (h UserHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(errors.Wrap(err, "error decode from request body"))
		return
	}
	err = h.Controller.Create(postgres.User{
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  false,
	})
	if err != nil {
		log.Println(errors.Wrap(err, "can't create user"))
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}
