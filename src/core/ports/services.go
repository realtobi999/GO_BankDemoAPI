package ports

import (
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

type IAccountService interface {
	Index(limit int, offset int) ([]domain.Account, error)
	Get(accountID uuid.UUID) (domain.Account, error)
	Create(customerID uuid.UUID, body domain.CreateAccountRequest) (domain.Account, error)
	Update(accountID uuid.UUID, body domain.UpdateAccountRequest) (int64, error)
	Delete(accountID uuid.UUID) (int64, error)
	IsOwner(customerID, accountID uuid.UUID) (bool, error)
}

type ICustomerService interface {
	Index(limit, offset int) ([]domain.Customer, error)
	Get(customerID uuid.UUID) (domain.Customer, error)
	Create(body domain.CreateCustomerRequest) (domain.Customer, error)
	Update(customerID uuid.UUID, body domain.UpdateCustomerRequest) (int64, error)
	Delete(customerID uuid.UUID) (int64, error)
	Auth(customerID uuid.UUID, token string) (bool, error)
}

type ITransactionService interface {
	Index(accountID uuid.UUID, limit int, offset int) ([]domain.Transaction, error)
	Get(transactionID uuid.UUID) (domain.Transaction, error)	
	Create(body domain.CreateTransactionRequest) (domain.Transaction, error)
}