package api

import (
	"net/http"

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

	
	account, err := body.ToAccount();
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
	
}
func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	
}
func (s *Server) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {

}
