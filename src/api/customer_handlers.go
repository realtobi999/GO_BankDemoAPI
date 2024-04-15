package api

import (
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
	"github.com/realtobi999/GO_BankDemoApi/src/utils/custom_errors"
)

func (s *Server) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Decode[types.CreateCustomerRequest](r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
	}

	if err := body.Validate(); err != nil {
		RespondWithValidationErrors(w, s.Logger, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	// Convert the body and create types.Customer struct
	customer := body.ToCustomer()

	_, err = s.Storage.CreateCustomer(customer)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	w.Header().Set("Location", "/api/customers/"+customer.ID.String())
	RespondWithJsonAndSerialize(w, http.StatusCreated, customer)
}
func (s *Server) IndexCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse limit and offset parameters
	limit, offset, err := utils.ParseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	customers, err := s.Storage.GetAllCustomers(limit, offset)
	if err != nil {
		if err.Error() == custom_errors.StorageNoResultsFound {
			RespondWithError(w, s.Logger, http.StatusNotFound, "No customers found!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to fetch customers: "+err.Error())
		return
	}

	RespondWithJsonAndSerializeList(w, http.StatusOK, customers)

}
func (s *Server) GetCustomerHandler(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {

}
