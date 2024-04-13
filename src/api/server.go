package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

type Server struct {
	Port   int
	Router *chi.Mux
	Logger types.ILogger
	Storage types.IStorage
}

func NewServer(port int, storage types.IStorage,logger types.ILogger) Server {
	return Server{
		Port: port,
		Logger:  logger,
		Storage:  storage,
	}
}

func (s Server) Run() error {
	s.setupRouter()
	s.loadRoutes()

	return http.ListenAndServe(fmt.Sprintf(":%v", s.Port), s.Router)
}