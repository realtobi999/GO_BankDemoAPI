package web

import (
	"log"
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/adapters/handlers"
	"github.com/urfave/negroni"
)

func (s *Server) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)

		// Catch panics
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERROR]\tPanic recovered: %v", err)
				handlers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(lrw, r)

		log.Printf("[EVENT]\tRequest received: %s %s - Status code: %v", r.Method, r.URL.Path, lrw.Status())
	})
}
