package storage

import (
	"time"

	"github.com/beevik/guid"
)

type Customer struct {
	ID        guid.Guid
	FirstName string
	LastName  string
	Birthday  time.Time
	Email     string
	Phone     string
	State     string
	Address   string
	Accounts  []Account
}