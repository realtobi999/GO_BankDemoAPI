package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func (s *Server) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.LogEvent(fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path))
		defer func() {
			if err := recover(); err != nil {
				s.Logger.LogError(fmt.Sprintf("Panic recovered: %v", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)

	})
}

func (s *Server) WithToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	
		customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
		if err != nil {
			RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
			return
		}
		
		token , err := utils.GetTokenFromHeader(r.Header.Get("Authorization"))
		if err != nil {
			RespondWithError(w, s.Logger, http.StatusBadRequest, "Invalid authorization header: "+err.Error())
			return
		}

		authorized, err := s.Storage.AuthCustomerWithTokenExists(customerID, token)
		if err != nil {
			RespondWithError(w, s.Logger, http.StatusInternalServerError, "Something went wrong: "+err.Error())
			return
		}

		if (!authorized) {
			RespondWithError(w, s.Logger, http.StatusUnauthorized, "Not authorized! Bad credentials")
			return
		}

		next.ServeHTTP(w, r)
	})
}
