package storage

import (
	"time"

	"github.com/beevik/guid"
)

type Account struct {
	ID                  guid.Guid
	CustomerID          guid.Guid
	Type                string
	Balance             float64
	Currency            string
	Status              bool
	OpeningDate         time.Time
	LastTransactionDate time.Time
	InterestRate        float64
}
