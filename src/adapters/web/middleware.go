package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/handlers"
	"github.com/realtobi999/GO_BankDemoApi/src/core/services/customer"
)

func (s *Server) LoadSharedMiddleware() {
	s.Router.Use(s.Logging)
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))
}

func (s *Server) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[EVENT]\tRequest received: %s %s", r.Method, r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERROR]\tPanic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)

	})
}

func (s *Server) WithToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := customer.GetTokenFromHeader(r.Header.Get("Authorization"))
		if err != nil {
			handlers.RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
			return
		}

		customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
			handlers.RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
			return
		}

		authorized, err := s.CustomerService.Auth(customerID, token)
		if err != nil {
			handlers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong: "+err.Error())
			return
		}

		if (!authorized) {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Not authorized! Bad credentials")
			return
		}

		next.ServeHTTP(w, r)
	})
}
