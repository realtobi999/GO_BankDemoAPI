package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
	"github.com/realtobi999/GO_BankDemoApi/src/utils/custom_errors"
)

func CreateCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	body := types.CreateCustomerRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	    RespondWithError(w, http.StatusBadRequest, "Failed to parse the request: "+err.Error())
		return
	}

	if err := body.Validate(); err != nil {
	    RespondWithValidationErrors(w, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	// Convert the body and create types.Customer struct
	customer := body.ToCustomer()

	_, err := s.CreateCustomer(customer);
	if err != nil {
	    RespondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	w.Header().Set("Location", "/api/customers/"+customer.ID.String())
    RespondWithJsonAndSerialize(w, http.StatusCreated, customer)
}
func IndexCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	// Parse limit and offset parameters
	limit, offset, err := utils.ParseLimitOffsetParams(r)
	if err != nil {
	    RespondWithError(w, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	customers, err := s.GetAllCustomers(limit, offset)
	if err != nil {
		if err.Error() == custom_errors.StorageNoResultsFound {
			RespondWithError(w, http.StatusNotFound, "No customers found!")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch customers: "+err.Error())
		return
	}

	// We do this conversion because we cant
	// safely pass the argument into the function
	var serialized []types.ISerializable
	for _, value := range(customers){
		serialized = append(serialized, value)
	}
	RespondWithJsonAndSerializeList(w, http.StatusOK, serialized)

}
func GetCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}
func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}
func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {

}