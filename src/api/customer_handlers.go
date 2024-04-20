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
		return
	}

	// Convert the body and create types.Customer struct
	// We also create the Id and the Token and CreatedAt
	customer := body.ToCustomer()

	if err := customer.Validate(); err != nil {
		RespondWithValidationErrors(w, s.Logger, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	_, err = s.Storage.CreateCustomer(customer)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	data := struct {
		Customer types.DTO `json:"customer"`
		Token    string            `json:"token"`
	}{
		Customer: customer.ToDTO(),
		Token:    customer.Token,
	}

	w.Header().Set("Location", "/api/customers/"+customer.ID.String())
	RespondWithJson(w, http.StatusCreated, data)
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
		if err == sql.ErrNoRows {
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
	id, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
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
	body, err := utils.Decode[types.UpdateCustomerRequest](r)
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}

	// Convert the body and create types.Customer struct
	customer := body.ToCustomer()

	// Set the ID for the body
	customer.ID, err = uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	// Validate
	if err := customer.Validate(); err != nil {
		RespondWithValidationErrors(w, s.Logger, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	// Update the database
	if err := s.Storage.UpdateCustomer(customer); err != nil {
		RespondWithError(w, s.Logger, http.StatusInternalServerError, "Failed to update the field in the database: "+err.Error())
		return
	}

	RespondWithJson(w, http.StatusOK, nil)
}
func (s *Server) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the UUID id
	id, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, s.Logger, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	// Delete the field from the database
	_, err = s.Storage.DeleteCustomer(id)
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
