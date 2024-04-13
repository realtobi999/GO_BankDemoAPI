package types

import (
	"time"

	"github.com/beevik/guid"
)

type Account struct {
	ID guid.Guid
	CustomerID guid.Guid
	Balance float64
	Type string
	Currency string
	Status bool
	OpeningDate time.Time
	LastTransactionDate time.Time
	InterestRate float64
}