package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type AccountHandler struct {
	AccountService ports.IAccountService
}

func NewAccountHandler(accountService ports.IAccountService) *AccountHandler {
	return &AccountHandler{
		AccountService: accountService,
	}
}

func (h *AccountHandler) Index(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	customerID := uuid.Nil
	if  r.URL.Query().Get("customer_id") != "" {
		customerID, err = uuid.Parse(r.URL.Query().Get("customer_id"))
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
			return
		}
	}

	accounts, err := h.AccountService.Index(customerID, limit, offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJsonAndSerializeList(w, http.StatusOK, accounts)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	account, err := h.AccountService.Get(accountID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJsonAndSerialize(w, http.StatusOK, account)
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decode[domain.CreateAccountRequest](r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}

	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}
	
	account, err := h.AccountService.Create(customerID, body)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			RespondWithValidationErrors(w, http.StatusBadRequest, "Failed to validate request", domain.ExtractValidationErrorsToList(err))
			return
		}
		if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Location", fmt.Sprintf("/api/customer/%s/account/%s", customerID.String(), account.ID.String()))
	RespondWithJson(w, http.StatusCreated, nil)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := decode[domain.UpdateAccountRequest](r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	_, err = h.AccountService.Update(accountID, body)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			RespondWithValidationErrors(w, http.StatusBadRequest, "Failed to validate request", domain.ExtractValidationErrorsToList(err))
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJson(w, http.StatusOK, nil)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	_, err = h.AccountService.Delete(accountID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJson(w, http.StatusOK, nil)
}