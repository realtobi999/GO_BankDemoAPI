package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

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
		"CustomerID": "%s",
		"Balance": 1000.00,
		"Type": 1,
		"Currency": "USD"
	  }	  
	`, customer.ID.String())

	req, err := http.NewRequest("POST", "/api/account", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler := server.TestHandler(server.CreateAccountHandler)
	handler.ServeHTTP(recorder, req)

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

	body := `{
		"CustomerID": "7f407b99-75d3-4253-b09f-98564538e337",
		"Balance": -1000.00,
		"Type": 134,
		"Currency": "USD"
	}`
	

	req, err := http.NewRequest("POST", "/api/account", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler := server.TestHandler(server.CreateAccountHandler)
	handler.ServeHTTP(recorder, req)

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