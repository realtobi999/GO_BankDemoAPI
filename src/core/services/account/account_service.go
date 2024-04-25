package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type AccountService struct {
	AccountRepository ports.IAccountRepository
}

func NewAccountService(accountRepository ports.IAccountRepository) *AccountService {
	return &AccountService{
		AccountRepository: accountRepository,
	}
}

func (ac *AccountService) Index(customerID uuid.UUID, limit int, offset int) ([]domain.Account, error) {
	// Declare variables for accounts and error, because
	// we can then access them in if/else scope
	var accounts []domain.Account
	var err error

	if customerID != uuid.Nil {
		accounts, err = ac.AccountRepository.GetAllAccountsByCustomer(customerID, limit, offset)
	} else {
		accounts, err = ac.AccountRepository.GetAllAccounts(limit, offset)
	}

	// Handle error for both options
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,  domain.NotFoundError("Accounts not found")
		}
		return nil,  domain.InternalFailure(err)
	}

	return accounts, nil
}

func (ac *AccountService) Get(accountID uuid.UUID) (domain.Account, error) {
	account, err := ac.AccountRepository.GetAccount(accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Account{}, domain.NotFoundError("Account not found")
		}
		return domain.Account{}, domain.InternalFailure(err)
	}

	return account, nil
}

func (ac *AccountService) Create(customerID uuid.UUID, body domain.CreateAccountRequest) (domain.Account, error) {
	account := domain.Account{
		ID: uuid.New(),
		CustomerID: customerID,
		Balance: body.Balance,
		Type: body.Type,
		Currency: body.Currency,
		Status: true,
		OpeningDate: time.Now(),
		InterestRate: body.InterestRate,
		CreatedAt: time.Now(),
	}

	if err := account.Validate(); err != nil {
		return domain.Account{}, domain.ValidationError(err)
	}

	_, err := ac.AccountRepository.CreateAccount(account)
	if err != nil {
		return domain.Account{}, domain.InternalFailure(err)
	}

	return account, nil
}

func (ac *AccountService) Update(accountID uuid.UUID, body domain.UpdateAccountRequest) (int64, error) {
	account := domain.Account{
		ID: accountID,
		Balance: body.Balance,
		Type: body.Type,
		Currency: body.Currency,
		Status: body.Status,
		LastTransactionDate: body.LastTransactionDate,
		InterestRate: body.InterestRate,
	}

	affectedRows, err := ac.AccountRepository.UpdateAccount(account)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.NotFoundError("Account not found")
		}
		return 0, domain.InternalFailure(err)
	}

	if affectedRows == 0 {
		return 0, domain.InternalFailure(errors.New("No rows affected"))
	}

	return affectedRows, nil
}
func (ac *AccountService) Delete(accountID uuid.UUID) (int64, error) {
	affectedRows, err := ac.AccountRepository.DeleteAccount(accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.NotFoundError("Account not found")
		}
		return 0, domain.InternalFailure(err)
	}

	if affectedRows == 0 {
		return 0, domain.InternalFailure(errors.New("No rows affected"))
	}

	return affectedRows, nil
}

func (ac *AccountService) IsOwner(customerID, accountID uuid.UUID) (bool, error) {
	_, err := ac.AccountRepository.GetAccountByOwner(customerID, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, domain.InternalFailure(err)
	}

	return true, nil
}