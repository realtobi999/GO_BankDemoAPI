package types

import (
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