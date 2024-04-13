package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Port   int
	Router *chi.Mux
}

func NewServer(port int) Server {
	return Server{
		Port: port,
	}
}

func (s Server) Run() error {
	s.setupRouter()
	s.loadRoutes()

	return http.ListenAndServe(fmt.Sprintf(":%v", s.Port), s.Router)
}