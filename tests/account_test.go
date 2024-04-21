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
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/handlers"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

func Test_Account_Create_Works(t *testing.T) {
	customer := NewTestCustomer()
	account := NewTestAccount(customer.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.ClearAllTables()
	db.CreateCustomer(customer)
	db.CreateAccount(account)

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
	router.Post("/api/customer/{customer_id}/account", handlers.NewAccountHandler(server.AccountService).Create)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusCreated, recorder.Code)

	idStartIndex := strings.Index(recorder.Header().Get("Location"), url+"/")
	id := recorder.Header().Get("Location")[idStartIndex+len(url+"/"):]

	assertDatabaseHas(t, "accounts", "id", id, db)
}

func Test_Account_Create_ValidationWorks(t *testing.T) {
	db := NewTestDatabase()
	server := NewTestServer(db)

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
	router.Post("/api/customer/{customer_id}/account", handlers.NewAccountHandler(server.AccountService).Create)
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

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.ClearAllTables()
	db.CreateCustomer(customer)
	db.CreateAccount(account1)
	db.CreateAccount(account2)

	url := fmt.Sprintf("/api/customer/%s/account", customer.ID.String())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Get("/api/customer/{customer_id}/account", handlers.NewAccountHandler(server.AccountService).Index)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string              `json:"message"`
		Status  int                 `json:"status"`
		Data    []domain.AccountDTO    `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	// assertEqual(t, account1.ID.String(), body.Data[0].ID.String())
	// assertEqual(t, account2.ID.String(), body.Data[1].ID.String())
}

func Test_Account_Update_Works(t *testing.T) {
	customer := NewTestCustomer()
	account := NewTestAccount(customer.ID)
	account.Status = true
	account.InterestRate = 1

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.ClearAllTables()
	db.CreateCustomer(customer)
	db.CreateAccount(account)

	assertDatabaseMissing(t, "accounts", "status", false, db)
	assertDatabaseMissing(t, "accounts", "interest_rate", 0.025, db)

	body := `
	{
		"Balance": 1000.00,
		"Type": 3,
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
	router.Put("/api/customer/{customer_id}/account/{account_id}", handlers.NewAccountHandler(server.AccountService).Update)
	router.ServeHTTP(recorder, req)
	
	assertEqual(t, http.StatusOK, recorder.Code)

	assertDatabaseHas(t, "accounts", "status", false, db)
	assertDatabaseHas(t, "accounts", "interest_rate", 0.025, db)
}

func Test_Account_Delete_Works(t *testing.T) {
	customer := NewTestCustomer()
	account := NewTestAccount(customer.ID)


	db := NewTestDatabase()
	server := NewTestServer(db)

	db.ClearAllTables()
	db.CreateCustomer(customer)
	db.CreateAccount(account)

	assertDatabaseHas(t, "accounts", "id", account.ID.String(), db)

	url := fmt.Sprintf("/api/customer/%s/account/%s", account.CustomerID.String(), account.ID.String())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Delete("/api/customer/{customer_id}/account/{account_id}", handlers.NewAccountHandler(server.AccountService).Delete)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)
	assertDatabaseMissing(t, "accounts", "id", account.ID.String(), db)
}
