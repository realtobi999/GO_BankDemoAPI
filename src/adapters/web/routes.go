package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/handlers"
)

func (s *Server) LoadRoutes() {
	accountHandler := handlers.NewAccountHandler(s.AccountService)
	customerHandler := handlers.NewCustomerHandler(s.CustomerService)

	s.Router.Route("/api", func(r chi.Router) {
		r.Route("/customer", func(r chi.Router) {
			r.Get("/", customerHandler.Index) // Params: Limit, Offset
			r.Get("/{customer_id}", customerHandler.Get)
			r.Post("/", customerHandler.Create)
			r.With(s.WithToken).Put("/{customer_id}", customerHandler.Update)
			r.With(s.WithToken).Delete("/{customer_id}", customerHandler.Delete)

			r.With(s.WithToken).Route("/account", func(r chi.Router) {
				r.Get("/", accountHandler.Index) // Params: Limit, Offset
				r.Get("/{account_id}", accountHandler.Get)
				r.Post("/", accountHandler.Create)
				r.Put("/{account_id}", accountHandler.Update)
				r.Delete("/{account_id}", accountHandler.Delete)
			})
		})
	})
}