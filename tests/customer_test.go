package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
		Message string       `json:"message"`
		Status  int          `json:"status"`
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
		Message string       `json:"message"`
		Status  int          `json:"status"`
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