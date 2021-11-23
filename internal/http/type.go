package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server wraps up all the request routers and associated components
// that serve various parts of the nbuild stack.
type Server struct {
	r chi.Router

	n *http.Server
}
