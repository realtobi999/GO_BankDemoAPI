package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the body into the CreateAccountRequest struct
	body, err := utils.Decode[types.CreateAccountRequest](r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}
	
	account, err := body.ToAccount(chi.URLParam(r, "customer_id"));
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}
	
	// Validate the Account struct
	if err := account.Validate(); err != nil {
		RespondWithValidationErrors(w, s.Logger, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	// Store in the database
	_, err = s.Storage.CreateAccount(account)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	w.Header().Set("Location", "/api/account/"+account.ID.String())
	RespondWithJsonAndSerialize(w, http.StatusCreated, account)
}
func (s *Server) IndexAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse limit and offset parameters
	limit, offset, err := utils.ParseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	accounts, err := s.Storage.GetAllAccountsFrom(customerID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows{
			RespondWithError(w, s.Logger, http.StatusNotFound, "No accounts found!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to fetch accounts: "+err.Error())
		return
	}

	RespondWithJsonAndSerializeList(w, 200, accounts)
}
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	account, err := s.Storage.GetAccount(accountID, customerID)
	if err != nil {
		if err == sql.ErrNoRows{
			RespondWithError(w, s.Logger, http.StatusNotFound, "No account found!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to fetch accounts: "+err.Error())
		return
	}

	RespondWithJsonAndSerialize(w, http.StatusOK, account)
}
func (s *Server) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Decode[types.UpdateAccountRequest](r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	
	account := body.ToAccount(accountID, customerID)

	// Validate
	if err := account.Validate(); err != nil {
		RespondWithValidationErrors(w, s.Logger, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	if err := s.Storage.UpdateAccount(account); err != nil {
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to update the field in the database: "+err.Error())
		return
	}

	RespondWithJson(w, http.StatusOK, nil)
}
func (s *Server) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the UUID")
		return
	}

	_, err = s.Storage.DeleteAccount(accountID, customerID)
	if err != nil {
		if err.Error() == "no rows affected" {
			RespondWithError(w, s.Logger, http.StatusNotFound, "Customer with that ID doesnt exits!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to delete customer: "+err.Error())
		return
	}

	RespondWithJson(w, http.StatusOK, nil)
}
