package tests

import (
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

func Test_Transactions_Create_Works(t *testing.T) {
	customer1 := NewTestCustomer()
	customer2 := NewTestCustomer()

	senderAcc := NewTestAccount(customer1.ID)
	receiverAcc := NewTestAccount(customer2.ID)

	db := NewTestDatabase()
	server := NewTestServer(db)

	db.CreateCustomer(customer1)
	db.CreateCustomer(customer2)
	db.CreateAccount(senderAcc)
	db.CreateAccount(receiverAcc)

	body := fmt.Sprintf(`
	{
  		"ReceiverAccountID": "%s",
 	 	"Amount": 100.50,
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

	idStr := recorder.Header().Get("Location")[strings.Index(recorder.Header().Get("Location"), "/api/transaction/")+ len("/api/transaction/"):]
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

	//assertDatabaseHas(t, "accounts", "")
}