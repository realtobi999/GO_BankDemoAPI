package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	h "github.com/realtobi999/GO_BankDemoApi/src/handlers"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, types.ILogger)


func (s *Server) loadRoutes() {
	s.Router.Get("/api/health", s.handler(h.HealthTestHandler))
	s.Router.Get("/api/error", s.handler(h.ErrorTestHandler))

	s.Router.Get("/api/customer", s.handler(h.IndexCustomerHandler))
	s.Router.Get("/api/customer/{id}", s.handler(h.GetCustomerHandler))
	s.Router.Post("/api/customer", s.handler(h.CreateCustomerHandler))
	s.Router.Put("/api/customer/{id}", s.handler(h.UpdateCustomerHandler))
	s.Router.Delete("/api/customer/{id}", s.handler(h.DeleteCustomerHandler))
}

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}

func (s *Server) handler(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.LogEvent(fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path))
		handlerFunc(w, r, s.Logger)
	}
}
