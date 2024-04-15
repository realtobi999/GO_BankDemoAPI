package api

import (
	"net/http"
)

func (s *Server) HealthTestHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJson(w, 200, nil)
}

func (s *Server) ErrorTestHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, s.Logger, 500, "Something went wrong! Oops...",)
}