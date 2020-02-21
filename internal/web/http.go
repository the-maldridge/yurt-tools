package web

import (
	"github.com/labstack/echo"
)

// Server serves the user interface over http using Echo.
type Server struct {
	*echo.Echo

	staticDir string
}

// New initializes and returns a new http.Server.
func New(static, tmpl string) (*Server, error) {
	s := Server{
		Echo:      echo.New(),
		staticDir: static,
	}

	r, err := newRenderer(tmpl)
	if err != nil {
		return nil, err
	}
	s.Renderer = r
	s.Static("/static", s.staticDir)
	return &s, nil
}

func (s *Server) Serve() error {
	return s.Start(":8080")
}
