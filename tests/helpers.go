package tests

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/api"
	"github.com/realtobi999/GO_BankDemoApi/src/storage"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
	"github.com/realtobi999/GO_BankDemoApi/src/utils/logs"
)

func NewTestServer() api.Server {
	db, err := storage.NewPostgres("localhost", "5432", "postgres", "root", "GoBankApiTesting", "disable")
	if err != nil {
		panic(err)
	}

	logger := logs.NewLogger("./../"+logs.PathToTestLogs)

	if err := storage.RunMigrations("./../"+storage.PathToMigrations, db.DB, logger); err != nil {
		panic(err)
	}

	return api.Server{
		Port: 8080,
		Router: chi.NewMux(),
		Logger: logger,
		Storage: db,
	}
}

func NewTestCustomer() types.Customer {
	customerID, _ := uuid.NewRandom()

	return types.Customer{
		ID:        customerID,
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		Email:     "john@example.com",
		Phone:     "123-456-7890",
		State:     "CA",
		Address:   "123 Main St",
		CreatedAt: time.Now(),
		Token:	   utils.GenerateToken(),
	}
}

func NewTestAccount(customerID uuid.UUID) types.Account {
	return types.Account{
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

func assertDatabaseHas(t *testing.T, table, column string, value any, s types.IStorage) {
	t.Helper()
	if !s.DatabaseHas(table, column, value) {
		t.Errorf("expected %s to have %s = %s", table, column, value)
	}
}

func assertDatabaseMissing(t *testing.T, table, column string, value any, s types.IStorage) {
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
