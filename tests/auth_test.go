package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/realtobi999/GO_BankDemoApi/src/constants"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func Test_Middleware_WithToken_Works(t *testing.T) {
	customer := NewTestCustomer()

	server := NewTestServer()
	server.Storage.ClearAllTables()

	server.Storage.CreateCustomer(customer)

	url := fmt.Sprintf("/api/customer/%s", customer.ID.String())
	
	req, err := http.NewRequest("DELETE", url, nil) 
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+utils.GenerateToken()) // Pass in a different generated token
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.With(server.WithToken).Delete("/api/customer/{customer_id}", server.DeleteCustomerHandler)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusUnauthorized, recorder.Code)

	body := struct{
		ErrorMessage string `json:"error_message"`
		Code int `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, body.ErrorMessage, "Not authorized! Bad credentials")
}

func Test_Token_GenerateToken_Works(t *testing.T) {
	token := utils.GenerateToken()

	assertEqual(t, constants.TOKEN_LENGTH, len(token))
}

func Test_Token_GetFromHeader_ValidationWorks(t *testing.T) {
	_, err := utils.GetTokenFromHeader("")
	assertEqual(t, "invalid header", err.Error())
	_, err = utils.GetTokenFromHeader("BEARER_ TOKEN")
	assertEqual(t, "missing Bearer", err.Error())
	_, err = utils.GetTokenFromHeader("Bearer token")
	assertEqual(t, "invalid token", err.Error())
}