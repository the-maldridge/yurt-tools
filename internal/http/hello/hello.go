package hello

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Hello struct{}

func New() *Hello {
	return &Hello{}
}

func (h *Hello) HTTPEntry() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.index)
	r.Get("/generate_204", h.generate204)

	return r
}

func (h *Hello) index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (h *Hello) generate204(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
