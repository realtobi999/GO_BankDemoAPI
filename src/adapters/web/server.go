package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type Server struct {
	Addr string
	Router *chi.Mux
	AccountService ports.IAccountService
	CustomerService ports.ICustomerService
	TransactionService ports.ITransactionService
}

func NewServer(addr string, router *chi.Mux) *Server {
	return &Server{
		Addr: addr,
		Router: router,
	}
}

func (s *Server) Run() error {
	log.Println("[EVENT]\tServer running od address: "+s.Addr)
	return http.ListenAndServe(s.Addr, s.Router)
}