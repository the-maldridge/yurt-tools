package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
	s.fileServer(s.r, "/static", http.Dir("theme/static"))

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

func (s *Server) fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
