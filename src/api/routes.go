package api

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) loadRoutes() {
    s.Router.Use(s.Logging)
    s.Router.Route("/api", func(r chi.Router) {
        r.Get("/health", s.HealthTestHandler)
        r.Get("/error", s.ErrorTestHandler)

        r.Route("/customer", func(r chi.Router) {
            r.Get("/", s.IndexCustomerHandler)                                 // Params: limit, offset
            r.Get("/{customer_id}", s.GetCustomerHandler)
            r.Post("/", s.CreateCustomerHandler)                               // Body: types.CreateCustomerRequest
            r.With(s.WithToken).Put("/{customer_id}", s.UpdateCustomerHandler) // Body: types.UpdateCustomerRequest
            r.With(s.WithToken).Delete("/{customer_id}", s.DeleteCustomerHandler)

            r.With(s.WithToken).Route("/{customer_id}/account", func(r chi.Router) {
                r.Get("/", s.IndexAccountHandler)                     // Params: limit, offset
                r.Get("/{account_id}", s.GetAccountHandler)
                r.Post("/", s.CreateAccountHandler)                   // Body: types.CreateAccountRequest
                r.Put("/{account_id}", s.UpdateAccountHandler)        // Body: types.UpdateAccountRequest
                r.Delete("/{account_id}", s.DeleteAccountHandler)

                r.Route("/{account_id}/transaction", func(r chi.Router) {
                    r.Get("/", s.IndexTransactionHandler)  
                    // r.Get("/{transaction_id}", s.GetTransactionHandler)
                    // r.Get("/{receiver_id}/{transaction_id}", s.IndexTransactionWithReceiverHandler)
                    r.Post("/{receiver_id}", s.CreateTransactionHandler)
                })
            })
        })
    })
}


func (s *Server) setupRouter() {
	s.Router = chi.NewRouter()
}
