package web

import (
	"github.com/labstack/echo"
	"net/http"
)

// Server serves the user interface over http using Echo.
type Server struct {
	*echo.Echo
}

// New initializes and returns a new http.Server.
func New() (*Server, error) {
	s := Server{
		Echo: echo.New(),
	}

	s.GET("/*", echo.WrapHandler(http.FileServer(HTTP)))

	return &s, nil
}

func (s *Server) Serve() error {
	return s.Start(":8080")
}
