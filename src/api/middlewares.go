package api

import (
	"fmt"
	"net/http"
)

func (s *Server) logging(next http.Handler) http.Handler {
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
