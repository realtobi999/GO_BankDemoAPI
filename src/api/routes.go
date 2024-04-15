package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) loadRoutes() {
	s.Router.Route("/api", func(r chi.Router){
        r.Get("/health", s.handler(s.HealthTestHandler))
        r.Get("/error", s.handler(s.ErrorTestHandler))
    
        r.Route("/customer", func(r chi.Router){ // Use 'r' here instead of 's.Router'
            r.Get("/", s.handler(s.IndexCustomerHandler)) // Params: limit, offset
            r.Get("/{id}", s.handler(s.GetCustomerHandler))
            r.Post("/", s.handler(s.CreateCustomerHandler)) // Body: types.CreateCustomerRequest 
            r.Put("/{id}", s.handler(s.UpdateCustomerHandler))
            r.Delete("/{id}", s.handler(s.DeleteCustomerHandler))
        })
    })
}

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}

func (s *Server) handler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.LogEvent(fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path))
		handlerFunc(w, r)
	}
}
