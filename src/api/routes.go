package api

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) loadRoutes() {
	s.Router.Use(s.logging)
	s.Router.Route("/api", func(r chi.Router){
        r.Get("/health", s.HealthTestHandler)
        r.Get("/error", s.ErrorTestHandler)
    
        r.Route("/customer", func(r chi.Router){
            r.Get("/", s.IndexCustomerHandler) // Params: limit, offset
            r.Get("/{id}", s.GetCustomerHandler)
            r.Post("/", s.CreateCustomerHandler) // Body: types.CreateCustomerRequest 
            r.Put("/{id}", s.UpdateCustomerHandler) // Body: types.UpdateCustomerRequest
            r.Delete("/{id}", s.DeleteCustomerHandler)

			r.Route("/{customer_id}/account", func(r chi.Router) {
				r.Get("/", s.IndexAccountHandler)
				r.Get("/{id}", s.GetAccountHandler)
				r.Post("/", s.CreateAccountHandler) // Body: types.CreateAccountRequest
				r.Put("/{id}", s.UpdateAccountHandler)
				r.Delete("/{id}", s.DeleteAccountHandler)
			})
        })

    })
}

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}
