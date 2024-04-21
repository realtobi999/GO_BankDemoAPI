package ports

import (
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

type IRepository interface {
	DatabaseHas(table, column string, value any) bool
	ClearAllTables() error
}

type IAccountRepository interface {
	GetAllAccounts(limit int, offset int) ([]domain.Account, error)
	GetAccount(accountID uuid.UUID) (domain.Account, error)
	CreateAccount(account domain.Account) (int64, error)
	UpdateAccount(account domain.Account) (int64, error)
	DeleteAccount( accountID uuid.UUID) (int64, error)
}

type ICustomerRepository interface {
	GetCustomer(id uuid.UUID) (domain.Customer, error)
	GetAllCustomers(limit int, offset int) ([]domain.Customer, error)
	CreateCustomer(customer domain.Customer) (int64, error)
	UpdateCustomer(customer domain.Customer) (int64, error)
	DeleteCustomer(customerID uuid.UUID) (int64, error)
	AuthCustomer(customerID uuid.UUID, token string) (bool, error)
}