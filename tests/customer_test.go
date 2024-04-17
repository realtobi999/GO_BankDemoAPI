package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

func Test_Customer_GetAll_Works(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()
	customer3 := NewTestCustomer()

	server := NewTestServer()

	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer1)
	server.Storage.CreateCustomer(customer2)
	server.Storage.CreateCustomer(customer3)

	req, err := http.NewRequest("GET", "/api/customer", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := server.TestHandler(server.IndexCustomerHandler)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string              `json:"message"`
		Status  int                 `json:"status"`
		Data    []types.CustomerDTO `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 3, len(body.Data))
}

func Test_Customer_GetAll_FailsWhenNoResults(t *testing.T) {
	server := NewTestServer()

	server.Storage.ClearAllTables()

	req, err := http.NewRequest("GET", "/api/customer", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := server.TestHandler(server.IndexCustomerHandler)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusNotFound, recorder.Code)

	body := struct {
		ErrorMessage string `json:"error_message"`
		Code         int    `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, body.ErrorMessage, "No customers found!")
}

func Test_Customer_GetAll_TestLimitAndOffset(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()
	customer3 := NewTestCustomer()
	customer4 := NewTestCustomer()

	server := NewTestServer()

	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer1)
	server.Storage.CreateCustomer(customer2)
	server.Storage.CreateCustomer(customer3)
	server.Storage.CreateCustomer(customer4)

	offset := 1
	limit := 2

	url := fmt.Sprintf("/api/customers?limit=%v&offset=%v", limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := server.TestHandler(server.IndexCustomerHandler)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string              `json:"message"`
		Status  int                 `json:"status"`
		Data    []types.CustomerDTO `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, limit, len(body.Data))

	// Test for the offset if the offset is 1 then the customer2
	// would match the first customer in the response
	assertEqual(t, customer2.ID.String(), body.Data[0].ID)
}

func Test_Customer_GetSpecific_Works(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()
	customer3 := NewTestCustomer()
	customer4 := NewTestCustomer()

	server := NewTestServer()

	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer1)
	server.Storage.CreateCustomer(customer2)
	server.Storage.CreateCustomer(customer3)
	server.Storage.CreateCustomer(customer4)

	url := fmt.Sprintf("/api/customers/%s", customer1.ID.String())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Get("/api/customers/{id}", server.TestHandler(server.GetCustomerHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string            `json:"message"`
		Status  int               `json:"status"`
		Data    types.CustomerDTO `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, customer1.ID.String(), body.Data.ID)
}

func Test_Customer_GetSpecific_FailsWhenNotFound(t *testing.T) {
	customer := NewTestCustomer()

	server := NewTestServer()

	url := fmt.Sprintf("/api/customers/%s", customer.ID.String())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Get("/api/customers/{id}", server.TestHandler(server.GetCustomerHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusNotFound, recorder.Code)

	body := struct {
		ErrorMessage string `json:"error_message"`
		Code         int    `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, http.StatusNotFound, body.Code)
}

func Test_Customer_Create_Works(t *testing.T) {
	customer := NewTestCustomer()
	server := NewTestServer()

	server.Storage.ClearAllTables()

	body := fmt.Sprintf(`
	{
		"first_name": "%s",
		"last_name": "%s",
		"birthday": "1990-01-01T00:00:00Z",
		"email": "%s",
		"phone": "%s",
		"state": "%s",
		"address": "%s"
	}
	`, customer.FirstName, customer.LastName, customer.Email, customer.Phone, customer.State, customer.Address)

	req, err := http.NewRequest("POST", "/api/customers", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler := server.TestHandler(server.CreateCustomerHandler)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusCreated, recorder.Code)

	rBody := struct {
		Message string            `json:"message"`
		Status  int               `json:"status"`
		Data    types.CustomerDTO `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}

	assertDatabaseHas(t, "customers", "id", rBody.Data.ID, server.Storage)
}

func Test_Customer_Create_ValidationWorks(t *testing.T) {
	customer := NewTestCustomer()
	server := NewTestServer()

	server.Storage.ClearAllTables()

	// missing fist_name and state
	body := fmt.Sprintf(`
	{
		"last_name": "%s",
		"birthday": "1990-01-01T00:00:00Z",
		"email": "%s",
		"phone": "%s",
		"address": "%s"
	}
	`, customer.LastName, customer.Email, customer.Phone, customer.Address)

	req, err := http.NewRequest("POST", "/api/customers", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler := server.TestHandler(server.CreateCustomerHandler)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusBadRequest, recorder.Code)

	rBody := struct {
		Message string   `json:"message"`
		Code    int      `json:"status"`
		Errors  []string `json:"errors"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "first name is required", rBody.Errors[0])
	assertEqual(t, "state is required", rBody.Errors[1])
}

func Test_Customer_Update_Works(t *testing.T) {
	customer1 := NewTestCustomer()

	server := NewTestServer()

	server.Storage.ClearAllTables()
	
	server.Storage.CreateCustomer(customer1)

	newFirstName := "Tobiáš"
	newLastName := "Filgas"
	newState := "Vsetín"

	body := fmt.Sprintf(`
	{
		"first_name": "%s",
		"last_name": "%s",
		"birthday": "1990-01-01T00:00:00Z",
		"email": "john.doe@example.com",
		"phone": "+420605401050",
		"state": "%s",
		"address": "123 Main St"
	}
	`, newFirstName, newLastName, newState)

	url := fmt.Sprintf("/api/customer/%s", customer1.ID.String())

	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Put("/api/customer/{id}", server.TestHandler(server.UpdateCustomerHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	updatedCustomer, err := server.Storage.GetCustomer(customer1.ID)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, newFirstName, updatedCustomer.FirstName)
	assertEqual(t, newLastName, updatedCustomer.LastName)
	assertEqual(t, newState, updatedCustomer.State)
}

func Test_Customer_Update_ValidationWorks(t *testing.T) {
	customer := NewTestCustomer()
	server := NewTestServer()

	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer)

	// missing fist_name and state
	body := fmt.Sprintf(`
	{
		"last_name": "%s",
		"birthday": "1990-01-01T00:00:00Z",
		"email": "%s",
		"phone": "%s",
		"address": "%s"
	}
	`, customer.LastName, customer.Email, customer.Phone, customer.Address)

	url := fmt.Sprintf("/api/customer/%s", customer.ID.String())

	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Put("/api/customer/{id}", server.TestHandler(server.UpdateCustomerHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusBadRequest, recorder.Code)

	rBody := struct {
		Message string   `json:"message"`
		Code    int      `json:"status"`
		Errors  []string `json:"errors"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "first name is required", rBody.Errors[0])
	assertEqual(t, "state is required", rBody.Errors[1])
}

func Test_Customer_Delete_Works(t *testing.T) {
	customer := NewTestCustomer()
	server := NewTestServer()

	server.Storage.ClearAllTables()
	server.Storage.CreateCustomer(customer)
		
	url := fmt.Sprintf("/api/customer/%s", customer.ID.String())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Delete("/api/customer/{id}", server.TestHandler(server.DeleteCustomerHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)
	assertDatabaseMissing(t, "customers", "id", customer.ID, server.Storage);
}
