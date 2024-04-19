package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

func Test_Account_Create_Works(t *testing.T) {
	customer := NewTestCustomer()
	account := NewTestAccount(customer.ID)

	server := NewTestServer()
	server.Storage.ClearAllTables()
	
	server.Storage.CreateCustomer(customer)
	server.Storage.CreateAccount(account)

	body := fmt.Sprintf(`
	{
		"Balance": 1000.00,
		"Type": 1,
		"Currency": "USD"
	  }	  
	`)

	url := fmt.Sprintf("/api/customer/%s/account", customer.ID.String())

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Post("/api/customer/{customer_id}/account", server.TestHandler(server.CreateAccountHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusCreated, recorder.Code)

	rBody := struct {
		Message string              `json:"message"`
		Status  int                 `json:"status"`
		Data    types.AccountDTO    `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}

	assertDatabaseHas(t, "accounts", "id", rBody.Data.ID.String(), server.Storage)
	assertEqual(t, customer.ID.String(), rBody.Data.CustomerID.String())
}

func Test_Account_Create_ValidationWorks(t *testing.T) {
	server := NewTestServer()
	customer := NewTestCustomer()

	body := `{
		"Balance": -1000.00,
		"Type": 134,
		"Currency": "USD"
	}`
	
	url := fmt.Sprintf("/api/customer/%s/account", customer.ID.String())

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Post("/api/customer/{customer_id}/account", server.TestHandler(server.CreateAccountHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusBadRequest, recorder.Code)

	rBody := struct{
		ErrorMessage string   `json:"message"`
		Code         int      `json:"status"`
		Errors       []string `json:"errors"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}


	assertEqual(t, true, slices.Contains(rBody.Errors, "Balance cannot be negative"))
	assertEqual(t, true, slices.Contains(rBody.Errors, "Invalid account type"))
}

func Test_Account_GetAll_Works(t *testing.T) {
	customer := NewTestCustomer()

	account1 := NewTestAccount(customer.ID)
	account2 := NewTestAccount(customer.ID)

	server := NewTestServer()
	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer)
	server.Storage.CreateAccount(account1)
	server.Storage.CreateAccount(account2)

	url := fmt.Sprintf("/api/customer/%s/account", customer.ID.String())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Get("/api/customer/{customer_id}/account", server.TestHandler(server.IndexAccountHandler))
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string              `json:"message"`
		Status  int                 `json:"status"`
		Data    []types.AccountDTO    `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, account1.ID.String(), body.Data[0].ID.String())
	assertEqual(t, account2.ID.String(), body.Data[1].ID.String())
}

func Test_Account_Update_Works(t *testing.T) {
	customer := NewTestCustomer()
	account := NewTestAccount(customer.ID)
	account.Status = true
	account.InterestRate = 1

	server := NewTestServer()
	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer)
	server.Storage.CreateAccount(account)

	assertDatabaseMissing(t, "accounts", "status", false, server.Storage)
	assertDatabaseMissing(t, "accounts", "interest_rate", 0.025, server.Storage)

	body := `
	{
		"Balance": 1000.00,
		"Type": 1,
		"Currency": "USD",
		"Status": false,
		"LastTransactionDate": "2024-04-19T12:00:00Z",
		"InterestRate": 0.025
	}`
		
	url := fmt.Sprintf("/api/customer/%s/account/%s", account.CustomerID.String(), account.ID.String())

	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Put("/api/customer/{customer_id}/account/{id}", server.TestHandler(server.UpdateAccountHandler))
	router.ServeHTTP(recorder, req)
	
	assertEqual(t, http.StatusOK, recorder.Code)

	assertDatabaseHas(t, "accounts", "status", false, server.Storage)
	assertDatabaseHas(t, "accounts", "interest_rate", 0.025, server.Storage)
}