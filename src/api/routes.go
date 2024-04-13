package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	h "github.com/realtobi999/GO_BankDemoApi/src/handlers"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, types.ILogger, types.IStorage)


func (s *Server) loadRoutes() {
	s.Router.Get("/api/health", s.handler(h.HealthTestHandler))
	s.Router.Get("/api/error", s.handler(h.ErrorTestHandler))
}

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}

func (s *Server) handler(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.LogEvent(fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path))
		handlerFunc(w, r, s.Logger, s.Storage)
	}
}
