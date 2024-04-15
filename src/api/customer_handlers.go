package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func (s *Server) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the body into the given struct
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
		if err == sql.ErrNoRows{
			RespondWithError(w, s.Logger, http.StatusNotFound, "No customers found!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to fetch customers: "+err.Error())
		return
	}

	RespondWithJsonAndSerializeList(w, http.StatusOK, customers)

}
func (s *Server) GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the UUID id
	idStr := chi.URLParam(r, "id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
	}

	customer, err := s.Storage.GetCustomer(id)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, s.Logger, http.StatusNotFound, "No customer found!")
			return
		}
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to fetch customer: "+err.Error())
		return
	}

	RespondWithJsonAndSerialize(w, http.StatusOK, customer)
}
func (s *Server) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {

}
