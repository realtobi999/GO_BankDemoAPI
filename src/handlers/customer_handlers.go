package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
	u "github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func CreateCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	body := types.CreateCustomerRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Failed to parse the request: "+err.Error())
		return
	}

	if err := body.Validate(); err != nil {
		u.RespondWithValidationErrors(w, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	customer := body.ToCustomer()

	_, err := s.CreateCustomer(customer);
	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	w.Header().Set("Location", "/api/customers/"+customer.ID.String())
	u.RespondWithJson(w, 201, customer.ToDTO())
}
func IndexCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}
func GetCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}
func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}
func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}