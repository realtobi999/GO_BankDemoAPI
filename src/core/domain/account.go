package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                  uuid.UUID
	CustomerID          uuid.UUID
	Balance             float64
	Type                AccountType
	Currency            Currency
	Status              bool
	OpeningDate         time.Time
	LastTransactionDate time.Time
	InterestRate        float64
	CreatedAt           time.Time
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

type AccountType int

var AccountTypes = map[AccountType]string{
	1: "Business",
	2: "Personal",
	3: "Savings",
}

/* ------------------------------------------------------------ */
func (a Account) Validate() *ValidationErrors {
    var errors []string

    if a.ID == uuid.Nil {
        errors = append(errors, "ID cannot be nil")
    }

    if a.CustomerID == uuid.Nil {
		errors = append(errors, "CustomerID cannot be nil")
	}
	
    if a.Balance < 0 {
        errors = append(errors, "Balance cannot be negative")
    }

    if _, ok := AccountTypes[a.Type]; !ok {
        errors = append(errors, "Invalid account type")
    }

    if _, ok := CurrencyLookupMap[a.Currency]; !ok {
		errors = append(errors, "This currency is not supported!")
	}

    if a.InterestRate < 0 {
        errors = append(errors, "InterestRate cannot be negative")
    } else if a.Type != 3 && a.InterestRate != 0 {
		errors = append(errors, "Non-savings account cannot have interest rate")
	}

    if len(errors) > 0 {
        return &ValidationErrors{Errors: errors}
    }
    return nil
}