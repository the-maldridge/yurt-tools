package http

import (
	"log"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New initializes the server with its default routers.
func New() (*Server, error) {
	s := Server{
	r: chi.NewRouter(),
	n: &http.Server{},
}

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Heartbeat("/healthz"))

	s.r.Get("/-/alive", s.alive)

	return &s, nil
}

// Serve binds, initializes the mux, and serves forever.
func (s *Server) Serve(bind string) error {
	log.Println("HTTP is starting")
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}

func (s *Server) alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// Mount attaches a set of routes to the subpath specified by the path
// argument.
func (s *Server) Mount(path string, router chi.Router) {
	s.r.Mount(path, router)
}
