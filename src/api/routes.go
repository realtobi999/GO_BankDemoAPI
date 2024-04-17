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
    
        r.Route("/customer", func(r chi.Router){
            r.Get("/", s.handler(s.IndexCustomerHandler)) // Params: limit, offset
            r.Get("/{id}", s.handler(s.GetCustomerHandler))
            r.Post("/", s.handler(s.CreateCustomerHandler)) // Body: types.CreateCustomerRequest 
            r.Put("/{id}", s.handler(s.UpdateCustomerHandler))
            r.Delete("/{id}", s.handler(s.DeleteCustomerHandler))
        })

		r.Route("/account", func(r chi.Router) {
			r.Get("/", s.handler(s.IndexAccountHandler))
			r.Get("/{id}", s.handler(s.GetAccountHandler))
			r.Post("/", s.handler(s.CreateAccountHandler))
			r.Put("/{id}", s.handler(s.UpdateAccountHandler))
			r.Delete("/{id}", s.handler(s.DeleteAccountHandler))
		})
    })
}

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}

func (s *Server) handler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.LogEvent(fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path))
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				s.Logger.LogError(fmt.Sprintf("Panic recovered: %v", err))
				// Respond with a 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		handlerFunc(w, r)
	}
}

func (s *Server) TestHandler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
}
