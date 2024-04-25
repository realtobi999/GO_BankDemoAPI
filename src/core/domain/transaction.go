package domain

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

const MAX_TRANSFER_AMOUNT = 10000;

type Transaction struct {
	ID	uuid.UUID
	SenderAccountID uuid.UUID
	ReceiverAccountID uuid.UUID
	Amount float64
	CurrencyPair CurrencyPair
	CreatedAt time.Time
}

type TransactionDTO struct {
	ID	uuid.UUID
	SenderAccountID uuid.UUID
	ReceiverAccountID uuid.UUID
	Amount float64
	CurrencyPair string
	CreatedAt time.Time
}

type CreateTransactionRequest struct {
	SenderAccountID uuid.UUID
	ReceiverAccountID uuid.UUID
	Amount float64
	Currency string // The sender preferred currency
}

/* ------------------------------------------------------------ */
func (t Transaction) Validate() *ValidationErrors {
	var errors []string

	if t.ID == uuid.Nil {
		errors = append(errors, "ID cannot be nil")
	}

	if t.SenderAccountID == uuid.Nil || t.ReceiverAccountID == uuid.Nil {
		errors = append(errors, "Both accounts ID's must be set")
	} else if t.SenderAccountID == t.ReceiverAccountID {
		errors = append(errors, "Sender and Receiver account cant have the same ID")
	}
	
	if t.Amount <= 0 {
		errors = append(errors, "Sending amount must be bigger than 0!")
	}
	
	if t.Amount > MAX_TRANSFER_AMOUNT {
		errors = append(errors, "Sending amount must not be bigger than: "+strconv.Itoa(MAX_TRANSFER_AMOUNT))
	}
	
	if _, ok := CurrencyLookupMap[t.CurrencyPair.From]; !ok {
		errors = append(errors, "This currency is not supported!")
	}

	if _, ok := CurrencyLookupMap[t.CurrencyPair.To]; !ok {
		errors = append(errors, "This currency is not supported!")
	}

	if t.CreatedAt.IsZero() {
		errors = append(errors, "CreatedAt must be set")
	}

	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}

	return nil
}