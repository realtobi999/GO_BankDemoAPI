package api

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) loadRoutes() {
	s.Router.Use(s.Logging)
	s.Router.Route("/api", func(r chi.Router){
        r.Get("/health", s.HealthTestHandler)
        r.Get("/error", s.ErrorTestHandler)
    
        r.Route("/customer", func(r chi.Router){
            r.Get("/", s.IndexCustomerHandler) // Params: limit, offset
            r.Get("/{customer_id}", s.GetCustomerHandler)
            r.Post("/", s.CreateCustomerHandler) // Body: types.CreateCustomerRequest 
            r.Put("/{customer_id}", s.UpdateCustomerHandler) // Body: types.UpdateCustomerRequest
            r.Delete("/{customer_id}", s.DeleteCustomerHandler)

			r.Route("/{customer_id}/account", func(r chi.Router) {
				// r.Use(s.WithToken)
				r.Get("/", s.IndexAccountHandler)
				r.Get("/{account_id}", s.GetAccountHandler)
				r.Post("/", s.CreateAccountHandler) // Body: types.CreateAccountRequest
				r.Put("/{account_id}", s.UpdateAccountHandler) // Body: types.UpdateAccountRequest
				r.Delete("/{account_id}", s.DeleteAccountHandler)
			})
        })

    })
}

func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}
