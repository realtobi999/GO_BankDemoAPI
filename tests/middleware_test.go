package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	customerService "github.com/realtobi999/GO_BankDemoApi/src/core/services/customer"
)

func Test_Middleware_TokenAuth_Works(t *testing.T) {
	customer := NewTestCustomer()

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.ClearAllTables()
	db.CreateCustomer(customer)

	url := fmt.Sprintf("/api/customer/%s", customer.ID.String())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+customerService.GenerateToken()) // Pass in a different generated token
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.With(server.TokenAuth).Delete("/api/customer/{customer_id}", func(w http.ResponseWriter, r *http.Request) {panic("Middleware is not working!")})
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusUnauthorized, recorder.Code)

	body := struct {
		ErrorMessage string `json:"error_message"`
		Code         int    `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "Not authorized! Bad credentials", body.ErrorMessage)
}

func Test_Token_GenerateToken_Works(t *testing.T) {
	token := customerService.GenerateToken()

	assertEqual(t, customerService.TOKEN_LENGTH, len(token))
}

func Test_Token_GetFromHeader_ValidationWorks(t *testing.T) {
	_, err := customerService.GetTokenFromHeader("")
	assertEqual(t, "invalid header", err.Error())
	_, err = customerService.GetTokenFromHeader("BEARER_ TOKEN")
	assertEqual(t, "missing Bearer", err.Error())
	_, err = customerService.GetTokenFromHeader("Bearer token")
	assertEqual(t, "invalid token", err.Error())
}

// the idea is to set the customer_id of the account as customer2.ID but make request using customer.ID
func Test_Middleware_AccountOwnerAuth_Works(t *testing.T) {
	customer := NewTestCustomer()
	customer2 := NewTestCustomer()
	account := NewTestAccount(customer2.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer)
	db.CreateCustomer(customer2)
	db.CreateAccount(account)

	url := fmt.Sprintf("/api/%s/account/%s", customer.ID.String(), account.ID.String())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()	

	router := chi.NewRouter()
	router.With(server.AccountOwnerAuth).Delete("/api/{customer_id}/account/{account_id}", func(w http.ResponseWriter, r *http.Request) {panic("Middleware is not working!")})
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusUnauthorized, recorder.Code)

	body := struct {
		ErrorMessage string `json:"error_message"`
		Code         int    `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "Not authorized!", body.ErrorMessage)
}
