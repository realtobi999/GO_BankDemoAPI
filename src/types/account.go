package types

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountType int

var AccountTypes = map[AccountType]string{
	1: "Business",
	2: "Personal",
	3: "Savings",
}

type Account struct {
	ID uuid.UUID
	CustomerID uuid.UUID
	Balance float64
	Type AccountType
	Currency Currency
	Status bool
	OpeningDate time.Time
	LastTransactionDate time.Time
	InterestRate float64
	CreatedAt time.Time
}

type AccountDTO struct {
	ID uuid.UUID
	CustomerID uuid.UUID
	Balance float64
	Type AccountType
	Currency Currency
	Status bool
	OpeningDate time.Time
	LastTransactionDate time.Time
	InterestRate float64
	CreatedAt time.Time
}

type CreateAccountRequest struct {
	Balance float64
	Type AccountType
	Currency Currency
	InterestRate float64
}

type UpdateAccountRequest struct {
	Balance float64
	Type 	AccountType
	Currency Currency
	Status	bool
	LastTransactionDate time.Time
	InterestRate float64
}

func (r CreateAccountRequest) ToAccount(customerID string) (Account, error) {
	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return Account{}, err
	}
	
	return Account{
		ID:                  uuid.New(),
		CustomerID:          customerUUID,
		Balance:             r.Balance,
		Type:                r.Type,
		Currency:            r.Currency,
		Status:              true,
		OpeningDate:         time.Now(),
		LastTransactionDate: time.Now(),
		InterestRate:        r.InterestRate,
		CreatedAt:  		 time.Now(),		
	}, nil
}

func (r UpdateAccountRequest) ToAccount(accountID, customerID uuid.UUID) Account {
	return Account{
		ID: accountID,
		CustomerID: customerID,
		Balance: r.Balance,
		Type: r.Type,
		Currency: r.Currency,
		Status: r.Status,
		LastTransactionDate: r.LastTransactionDate,
		InterestRate: r.InterestRate,
	}
}

func (a Account) ToDTO() DTO {
	return AccountDTO{
		ID: a.ID,
		CustomerID: a.CustomerID,
		Balance: a.Balance,
		Type: a.Type,
		Currency: a.Currency,
		Status: a.Status,
		OpeningDate: a.OpeningDate,
		LastTransactionDate: a.LastTransactionDate,
		InterestRate: a.InterestRate,
		CreatedAt: a.CreatedAt,
	}
}

func (a Account) Validate() []error {
    var validationErrors []error

    if a.ID == uuid.Nil {
        validationErrors = append(validationErrors, errors.New("ID cannot be nil"))
    }

    if a.CustomerID == uuid.Nil {
        validationErrors = append(validationErrors, errors.New("CustomerID cannot be nil"))
    }

    if a.Balance < 0 {
        validationErrors = append(validationErrors, errors.New("Balance cannot be negative"))
    }

    if _, ok := AccountTypes[a.Type]; !ok {
        validationErrors = append(validationErrors, errors.New("Invalid account type"))
    }

    if _, ok := CurrencyLookupMap[a.Currency]; !ok {
		validationErrors = append(validationErrors, errors.New("This currency is not supported!"))
	}

    if a.LastTransactionDate.IsZero() {
        validationErrors = append(validationErrors, errors.New("LastTransactionDate cannot be zero"))
    }

    if a.InterestRate < 0 {
        validationErrors = append(validationErrors, errors.New("InterestRate cannot be negative"))
    } else if a.Type != 3 && a.InterestRate != 0 {
		validationErrors = append(validationErrors, errors.New("Non-savings account cannot have interest rate"))
	}

    return validationErrors
}

