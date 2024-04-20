package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func (s *Server) IndexTransactionHandler(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.ParseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	transactions, err := s.Storage.GetAllTransactions(accountID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, s.Logger, http.StatusNotFound, "No transactions found!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to fetch transactions: "+err.Error())
		return
	}

	RespondWithJson(w, http.StatusOK, transactions)
}

func (s *Server) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	
}