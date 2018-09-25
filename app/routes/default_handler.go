package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type DefaultHandlers struct {
	Router *chi.Mux
}

func NewDefaultHandlers() DefaultHandlers {
	return DefaultHandlers{
		Router: DefaultHandlers{}.initRouter(),
	}
}

func (h DefaultHandlers) initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", h.hello)
	return r
}

func (h DefaultHandlers) hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API works"))
}

func ErrorHandler(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": false,
		"error":  err.Error(),
	})
}
