package ports

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

type IRepository interface {
	DatabaseHas(table, column string, value any) bool
	BeginTransaction() (*sql.Tx, error)
	ClearAllTables() error
}

type IAccountRepository interface {
	GetAllAccounts(limit int, offset int) ([]domain.Account, error)
	GetAllAccountsByCustomer(customerID uuid.UUID, limit int, offset int) ([]domain.Account, error)
	GetAllSavingsAccounts() ([]domain.Account, error) 
	GetAccount(accountID uuid.UUID) (domain.Account, error)
	GetAccountByOwner(customerID, accountID uuid.UUID) (domain.Account, error)
	CreateAccount(account domain.Account) (int64, error)
	UpdateAccount(account domain.Account) (int64, error)
	DeleteAccount(accountID uuid.UUID) (int64, error)	
}

type ICustomerRepository interface {
	GetAllCustomers(limit int, offset int) ([]domain.Customer, error)
	GetCustomer(customerID uuid.UUID) (domain.Customer, error)
	CreateCustomer(customer domain.Customer) (int64, error)
	UpdateCustomer(customer domain.Customer) (int64, error)
	DeleteCustomer(customerID uuid.UUID) (int64, error)
	AuthCustomer(customerID uuid.UUID, token string) (bool, error)
}

type ITransactionRepository interface {
	GetAllTransactions(limit, offset int) ([]domain.Transaction, error)
	GetAllTransactionsFromAccount(accountID uuid.UUID, limit int, offset int) ([]domain.Transaction, error)	
	GetTransaction(transactionID uuid.UUID) (domain.Transaction, error) 	
	CreateTransaction(transaction domain.Transaction) (int64, error)
}