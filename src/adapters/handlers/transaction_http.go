package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type TransactionHandler struct {
	TransactionService ports.ITransactionService
}

func NewTransactionHandler(transactionService ports.ITransactionService) TransactionHandler {
	return TransactionHandler{
		TransactionService: transactionService,
	}
}

func (h *TransactionHandler) Index(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	var accountID uuid.UUID
	accountIDStr := chi.URLParam(r, "account_id")
	if accountIDStr != "" {
		accountID, err = uuid.Parse(accountIDStr)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
			return
		}
	}

	transactions, err := h.TransactionService.Index(accountID, limit, offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	RespondWithJsonAndSerializeList(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	transactionID, err := uuid.Parse(chi.URLParam(r, "transaction_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	transaction, err := h.TransactionService.Get(transactionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJsonAndSerialize(w, http.StatusOK, transaction)
}
