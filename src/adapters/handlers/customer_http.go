package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type CustomerHandler struct {
	CustomerService ports.ICustomerService
}

func NewCustomerHandler(customerService ports.ICustomerService) *CustomerHandler {
	return &CustomerHandler{
		CustomerService: customerService,
	}
}

func (h *CustomerHandler) Index(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	accounts, err := h.CustomerService.Index(limit, offset)
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

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	customer, err := h.CustomerService.Get(customerID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJsonAndSerialize(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := decode[domain.CreateCustomerRequest](r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}

	customer, err := h.CustomerService.Create(body)
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

	response := struct {
		Token string `json:"token"`
	}{
		Token: customer.Token,
	}

	w.Header().Set("Location", "/api/customer/"+customer.ID.String())
	RespondWithJson(w, http.StatusCreated, response)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	body, err := decode[domain.UpdateCustomerRequest](r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}

	_, err = h.CustomerService.Update(customerID, body)
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

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	customerID, err := uuid.Parse(chi.URLParam(r, "customer_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	_, err = h.CustomerService.Delete(customerID)
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
}
