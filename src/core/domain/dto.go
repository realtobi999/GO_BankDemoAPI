package domain

import (
	"time"

	"github.com/google/uuid"
)

type DTO any

/* ------------------------------------------------------------ */
type AccountDTO struct {
	ID                  uuid.UUID
	CustomerID          uuid.UUID
	Balance             float64
	Type                string
	Currency            string
	Status              bool
	OpeningDate         time.Time
	LastTransactionDate time.Time
	InterestRate        float64
	CreatedAt           time.Time
}

func (a Account) ToDTO() DTO {
	return AccountDTO{
		ID: a.ID,
		CustomerID: a.CustomerID,
		Balance: a.Balance,
		Type: AccountLookupMap[a.Type],
		Currency: CurrencyLookupMap[a.Currency],
		Status: a.Status,
		OpeningDate: a.OpeningDate,
		LastTransactionDate: a.LastTransactionDate,
		InterestRate: a.InterestRate,
		CreatedAt: a.CreatedAt,
	}
}
/* ------------------------------------------------------------ */
type CustomerDTO struct {
	ID        string    
	FirstName string   
	LastName  string    
	Birthday  time.Time 
	Email     string    
	Phone     string    
	State     string    
	Address   string    
	CreatedAt time.Time 
}

func (c Customer) ToDTO() DTO {
	return CustomerDTO{
		ID:        c.ID.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Birthday:  c.Birthday,
		Email:     c.Email,
		Phone:     c.Phone,
		State:     c.State,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
	}
}
/* ------------------------------------------------------------ */
func (c Transaction) ToDTO() DTO {
	return TransactionDTO{
		ID:  c.ID,
		SenderAccountID: c.SenderAccountID,
		ReceiverAccountID: c.ReceiverAccountID,
		Amount: c.Amount,
		CurrencyPair: c.CurrencyPair.String(),
		CreatedAt: c.CreatedAt,
	}
}