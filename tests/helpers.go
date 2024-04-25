package tests

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/repository"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/repository/migrations"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/web"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
	"github.com/realtobi999/GO_BankDemoApi/src/core/services/account"
	"github.com/realtobi999/GO_BankDemoApi/src/core/services/customer"
	"github.com/realtobi999/GO_BankDemoApi/src/core/services/transactions"
)

func NewTestServer(db *repository.Postgres) *web.Server {
	server := web.NewServer(":8080", chi.NewMux())
	server.CustomerService = customer.NewCustomerService(db)
	server.AccountService = account.NewAccountService(db)
	server.TransactionService = transactions.NewTransactionService(db, db)

	if err := migrations.DropMigrations(db.DB); err != nil {
		panic(err)
	}

	if err := migrations.RunMigrations("./../src/adapters/repository/migrations/*.sql", db.DB); err != nil {
		panic(err)
	}

	return server
}

func NewTestDatabase() *repository.Postgres {
	db, err := repository.NewPostgres("localhost", "5432", "postgres", "root", "GoBankApiTesting", "disable")
	if err != nil {
		panic(err)
	}

	return db
} 

func NewTestCustomer() domain.Customer {
	customerID, _ := uuid.NewRandom()

	return domain.Customer{
		ID:        customerID,
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		Email:     "john@example.com",
		Phone:     "123-456-7890",
		State:     "CA",
		Address:   "123 Main St",
		CreatedAt: time.Now(),
		Token:	   customer.GenerateToken(),
	}
}

func NewTestAccount(customerID uuid.UUID) domain.Account {
	return domain.Account{
		ID:                  uuid.New(),
		CustomerID:          customerID,
		Balance:             0.0,
		Type:                1,
		Currency:            "USD",
		Status:              true,
		OpeningDate:         time.Now(),
		LastTransactionDate: time.Now(),
		InterestRate:        0.0,
	}
}

func NewTestTransaction(senderID uuid.UUID, receiver uuid.UUID) domain.Transaction {
	return domain.Transaction{
		ID: uuid.New(),
		SenderAccountID: senderID,
		ReceiverAccountID: receiver,
		Amount: 0,
		CurrencyPair: domain.NewCurrencyPair("USD", "EUR"),
		CreatedAt: time.Now(),
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func assertNotEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func assertDatabaseHas(t *testing.T, table, column string, value any, s ports.IRepository) {
	t.Helper()
	if !s.DatabaseHas(table, column, value) {
		t.Errorf("expected %s to have %s = %s", table, column, value)
	}
}

func assertDatabaseMissing(t *testing.T, table, column string, value any, s ports.IRepository) {
	t.Helper()
	if s.DatabaseHas(table, column, value) {
		t.Errorf("expected %s to have %s = %s", table, column, value)
	}
}

func generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
