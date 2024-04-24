package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/handlers"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

func Test_Transaction_GetAll_Works(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()

	senderAcc := NewTestAccount(customer1.ID)
	receiverAcc := NewTestAccount(customer2.ID)

	transaction1 := NewTestTransaction(senderAcc.ID, receiverAcc.ID)
	transaction2 := NewTestTransaction(senderAcc.ID, receiverAcc.ID)
	transaction3 := NewTestTransaction(senderAcc.ID, receiverAcc.ID)
	transaction4 := NewTestTransaction(senderAcc.ID, receiverAcc.ID)
	transaction5 := NewTestTransaction(senderAcc.ID, receiverAcc.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateCustomer(customer2)
	db.CreateAccount(senderAcc)
	db.CreateAccount(receiverAcc)	

	db.CreateTransaction(transaction1)
	db.CreateTransaction(transaction2)
	db.CreateTransaction(transaction3)
	db.CreateTransaction(transaction4)
	db.CreateTransaction(transaction5)

	req, err := http.NewRequest("GET", "/api/transaction", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.NewTransactionHandler(server.TransactionService).Index)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string                  `json:"message"`
		Code    int                     `json:"code"`
		Data    []domain.TransactionDTO `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, transaction1.ID.String(), body.Data[0].ID.String())	
	assertEqual(t, transaction2.ID.String(), body.Data[1].ID.String())	
	assertEqual(t, transaction3.ID.String(), body.Data[2].ID.String())	
	assertEqual(t, transaction4.ID.String(), body.Data[3].ID.String())	
	assertEqual(t, transaction5.ID.String(), body.Data[4].ID.String())	
}

func Test_Transaction_GetAll_GivesErrorWhenNothingFound(t *testing.T) {
	db := NewTestDatabase()
	server := NewTestServer(db)

	req, err := http.NewRequest("GET", "/api/transaction", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.NewTransactionHandler(server.TransactionService).Index)
	handler.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusNotFound, recorder.Code)

	body := struct {
		ErrorMessage string                  `json:"error_message"`
		Code    int                     `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "Error not found: Transactions not found", body.ErrorMessage)
}

func Test_Transaction_Get_Works(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()

	senderAcc := NewTestAccount(customer1.ID)
	receiverAcc := NewTestAccount(customer2.ID)

	transaction := NewTestTransaction(senderAcc.ID, receiverAcc.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateCustomer(customer2)
	db.CreateAccount(senderAcc)
	db.CreateAccount(receiverAcc)
	db.CreateTransaction(transaction)

	url := fmt.Sprintf("/api/transaction/%s", transaction.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}	
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Get("/api/transaction/{transaction_id}", handlers.NewTransactionHandler(server.TransactionService).Get)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusOK, recorder.Code)

	body := struct {
		Message string                  `json:"message"`
		Code    int                     `json:"code"`
		Data    domain.TransactionDTO `json:"data"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, transaction.ID, body.Data.ID)
}

func Test_Transaction_Get_GivesErrorWhenNotExist(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()

	senderAcc := NewTestAccount(customer1.ID)
	receiverAcc := NewTestAccount(customer2.ID)

	transaction := NewTestTransaction(senderAcc.ID, receiverAcc.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateCustomer(customer2)
	db.CreateAccount(senderAcc)
	db.CreateAccount(receiverAcc)

	url := fmt.Sprintf("/api/transaction/%s", transaction.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}	
	recorder := httptest.NewRecorder()

	router := chi.NewMux()
	router.Get("/api/transaction/{transaction_id}", handlers.NewTransactionHandler(server.TransactionService).Get)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusNotFound, recorder.Code)

	body := struct {
		ErrorMessage string                  `json:"error_message"`
		Code    int                     `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "Error not found: Transaction not found", body.ErrorMessage)
}

func Test_Transaction_Create_Works(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()

	senderAcc := NewTestAccount(customer1.ID)
	receiverAcc := NewTestAccount(customer2.ID)

	senderAcc.Balance = 1000
	receiverAcc.Currency = "EUR"

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateCustomer(customer2)
	db.CreateAccount(senderAcc)
	db.CreateAccount(receiverAcc)

	body := fmt.Sprintf(`
	{
  		"ReceiverAccountID": "%s",
 	 	"Amount": 100,
	  	"Currency": "USD"
	}
	`, receiverAcc.ID.String())

	url := fmt.Sprintf("/api/%s/account/%s/transaction", customer1.ID.String(), senderAcc.ID.String())

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/api/{customer_id}/account/{account_id}/transaction", handlers.NewTransactionHandler(server.TransactionService).Create)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusCreated, recorder.Code)

	idStr := recorder.Header().Get("Location")[strings.Index(recorder.Header().Get("Location"), "/api/transaction/")+len("/api/transaction/"):]
	assertDatabaseHas(t, "transactions", "id", idStr, db)

	id, err := uuid.Parse(idStr)
	if err != nil {
		t.Fatal(err)
	}

	transaction, err := server.TransactionService.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, senderAcc.ID.String(), transaction.SenderAccountID.String())
	assertEqual(t, receiverAcc.ID.String(), transaction.ReceiverAccountID.String())

	assertEqual(t, domain.Currency("USD"), transaction.CurrencyPair.From)
	assertEqual(t, receiverAcc.Currency, transaction.CurrencyPair.To)

	assertDatabaseHas(t, "accounts", "balance", domain.ConversionRateMap[transaction.CurrencyPair]*transaction.Amount, db)
	assertDatabaseHas(t, "accounts", "balance", senderAcc.Balance-transaction.Amount, db)
}

func Test_Transaction_Create_GivesErrorWhenSenderDoesntHaveEnoughBalance(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()

	sender := NewTestAccount(customer1.ID)
	receiver := NewTestAccount(customer2.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateCustomer(customer2)
	db.CreateAccount(sender)
	db.CreateAccount(receiver)

	sender.Balance = 0

	body := fmt.Sprintf(`
	{
  		"ReceiverAccountID": "%s",
 	 	"Amount": 1000000000,
	  	"Currency": "USD"
	}
	`, receiver.ID.String())

	url := fmt.Sprintf("/api/%s/account/%s/transaction", customer1.ID.String(), sender.ID.String())

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/api/{customer_id}/account/{account_id}/transaction", handlers.NewTransactionHandler(server.TransactionService).Create)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusBadRequest, recorder.Code)

	rBody := struct {
		ErrorMessage string `json:"error_message"`
		Code         int    `json:"code"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "Error bad request: Sender account doesnt have enough balance", rBody.ErrorMessage)
}

func Test_Transactions_Create_GivesErrorWhenSenderAndReceiverAreTheSame(t *testing.T) {
	customer1 := NewTestCustomer()

	receiver := NewTestAccount(customer1.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateAccount(receiver)

	body := fmt.Sprintf(`
	{
  		"ReceiverAccountID": "%s",
 	 	"Amount": 1000000000,
	  	"Currency": "USD"
	}
	`, receiver.ID.String())

	url := fmt.Sprintf("/api/%s/account/%s/transaction", customer1.ID.String(), receiver.ID.String())

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/api/{customer_id}/account/{account_id}/transaction", handlers.NewTransactionHandler(server.TransactionService).Create)
	router.ServeHTTP(recorder, req)

	assertEqual(t, http.StatusBadRequest, recorder.Code)

	rBody := struct {
		ErrorMessage string   `json:"error_message"`
		Code         int      `json:"code"`
		Errors       []string `json:"errors"`
	}{}
	if err := json.NewDecoder(recorder.Body).Decode(&rBody); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "Sender and Receiver account cant have the same ID", rBody.Errors[0])
}
