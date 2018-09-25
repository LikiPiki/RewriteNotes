package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	"github.com/likipiki/RewriteNotes/app/crypt"
	"github.com/likipiki/RewriteNotes/app/postgres"
	"github.com/pkg/errors"
)

type UserHandlers struct {
	controller postgres.UserService
	tokenAuth  *jwtauth.JWTAuth
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
		controller: controller,
		tokenAuth:  jwtauth.New("HS256", []byte("mysecret"), nil),
		Router:     UserHandlers{}.initRouter(),
	}
}

// Router - register all handler to chi router
func (h UserHandlers) initRouter() *chi.Mux {
	r := chi.NewRouter()
	// sub := chi.NewRouter()
	// r.Mount("/", sub)

	r.Post("/", h.Create)
	r.Post("/login", h.Login)
	return r
}

func (h UserHandlers) NoteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("id")
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}
		user, err := h.controller.Get(userID)
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
	err = h.controller.Create(postgres.User{
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  false,
	})
	if err != nil {
		log.Println(errors.Wrap(err, "can't create user"))
		http.Error(w, http.StatusText(404), 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func (h UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(errors.Wrap(err, "cant decode username and password from body"))
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if user.Username == "" || user.Password == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": false,
			"msg":    "Empty username or email",
		})
		return
	}
	dbUser, err := h.controller.Get(user.Username)
	if err != nil {
		log.Println(errors.Wrap(err, "can't create user"))
		return
	}
	if crypt.CheckPassword(dbUser.Password, user.Password) {
		_, token, err := h.tokenAuth.Encode(jwtauth.Claims{
			"username": dbUser.Username,
			"isAdmin":  dbUser.IsAdmin,
			"id":       dbUser.Id,
		})
		if err != nil {
			log.Println(errors.Wrap(err, "cant encode jwt token"))
			ErrorHandler(w, err)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": true,
			"token":  token,
		})
	} else {
		ErrorHandler(w, errors.New("Username or password not valid"))
	}
}

func (h UserHandlers) Get(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")
	if noteID == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	user, err := h.controller.GetById(noteID)
	if err != nil {
		log.Println(errors.Wrap(err, "Cant get user by id"))
	}
	result := User{
		Id:       user.Id,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    result,
	})
}
