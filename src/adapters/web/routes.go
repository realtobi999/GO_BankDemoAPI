package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/handlers"
)

func (s *Server) LoadRoutes() {
	accountHandler := handlers.NewAccountHandler(s.AccountService)
	customerHandler := handlers.NewCustomerHandler(s.CustomerService)
	transactionsHandler := handlers.NewTransactionHandler(s.TransactionService)

	s.Router.Route("/api", func(r chi.Router) {
		r.Route("/customer", func(r chi.Router) {
			r.Get("/", customerHandler.Index) // Params: Limit, Offset
			r.Get("/{customer_id}", customerHandler.Get)
			r.Post("/", customerHandler.Create)
			r.With(s.TokenAuth).Put("/{customer_id}", customerHandler.Update)
			r.With(s.TokenAuth).Delete("/{customer_id}", customerHandler.Delete)

			r.Route("/account", func(r chi.Router) {
				r.Get("/", accountHandler.Index) // Params: Limit, Offset
				r.Get("/{account_id}", accountHandler.Get)
				r.Post("/", accountHandler.Create)
			})

			r.With(s.TokenAuth, s.AccountOwnerAuth).Route("/{customer_id}/account", func(r chi.Router) {
				r.Put("/{account_id}", accountHandler.Update)
				r.Delete("/{account_id}", accountHandler.Delete)

				r.Post("/transaction", transactionsHandler.Create)
			})

		})

		r.Route("/transaction", func(r chi.Router) {
			r.Get("/", transactionsHandler.Index)
			r.Get("/{transaction_id}", transactionsHandler.Get)	
		})	
	})
}