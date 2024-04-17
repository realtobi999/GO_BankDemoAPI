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
	CreateCustomer(Customer) (int64, error)
	UpdateCustomer(Customer) error
	DeleteCustomer(id uuid.UUID) (int64, error)

	CreateAccount(Account) (int64, error)
}

type ISerializable interface {
	ToDTO() DTO
}

type DTO interface{}
