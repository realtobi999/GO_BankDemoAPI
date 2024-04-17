package types

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID uuid.UUID
	CustomerID uuid.UUID
	Balance float64
	Type string
	Currency string
	Status bool
	OpeningDate time.Time
	LastTransactionDate time.Time
	InterestRate float64
}

type AccountDTO struct {
	ID uuid.UUID
	CustomerID uuid.UUID
	Balance float64
	Type string
	Currency string
	Status bool
	OpeningDate time.Time
	LastTransactionDate time.Time
	InterestRate float64
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

    if a.Type == "" {
        validationErrors = append(validationErrors, errors.New("Type cannot be empty"))
    }

    if a.Currency == "" {
        validationErrors = append(validationErrors, errors.New("Currency cannot be empty"))
    }

    if a.OpeningDate.IsZero() {
        validationErrors = append(validationErrors, errors.New("OpeningDate cannot be zero"))
    }

    if a.LastTransactionDate.IsZero() {
        validationErrors = append(validationErrors, errors.New("LastTransactionDate cannot be zero"))
    }

    if a.InterestRate < 0 {
        validationErrors = append(validationErrors, errors.New("InterestRate cannot be negative"))
    }

    return validationErrors
}
