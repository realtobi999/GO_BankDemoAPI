package types

import "github.com/google/uuid"

type ILogger interface {
	LogEvent(message any)
	LogError(message any)
	LogDebug(message any)
	LogWarning(message any)
}

type IStorage interface {
	DatabaseHas(table, column string, value any) bool
	ClearAllTables() error

	GetCustomer(id uuid.UUID) (Customer, error)
	GetAllCustomers(limit int, offset int) ([]Customer, error)
	CreateCustomer(customer Customer) (int64, error)
	UpdateCustomer(customer Customer) error
	DeleteCustomer(id uuid.UUID) (int64, error)

	GetAllAccountsFrom(customerID uuid.UUID, limit int, offset int) ([]Account, error)
	GetAccount(accountID uuid.UUID, customerID uuid.UUID) (Account, error)
	CreateAccount(account Account) (int64, error)
	UpdateAccount(account Account) error
}

type ISerializable interface {
	ToDTO() DTO
}

type DTO interface{}
