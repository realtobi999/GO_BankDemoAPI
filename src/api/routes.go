package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	h "github.com/realtobi999/GO_BankDemoApi/src/handlers"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}

func (s *Server) loadRoutes() {
	s.Router.Get("/api/health", s.handler(h.HealthTestHandler))
	s.Router.Get("/api/error", s.handler(h.ErrorTestHandler))
}

func (s *Server) handler(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO]\tRequest received: %s %s", r.Method, r.URL.Path)
		handlerFunc(w, r)
	}
}
