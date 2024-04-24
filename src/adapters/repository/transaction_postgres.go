package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

func (p *Postgres) GetAllTransactions(limit, offset int) ([]domain.Transaction, error) {
	query := `SELECT * FROM transactions ORDER BY created_at LIMIT $1 OFFSET $2`
	
	rows ,err := p.DB.Query(query, limit, offset) 
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction

	for rows.Next() {
		var transaction domain.Transaction
		var currencyPair string

		if err := rows.Scan(&transaction.ID, &transaction.SenderAccountID, &transaction.ReceiverAccountID, &transaction.Amount, &currencyPair, &transaction.CreatedAt); err != nil {
			return nil ,err		
		}

		// Set the currency pair and skip over if its corrupted (It really shouldn't be)
		transaction.CurrencyPair, err = domain.CurrencyPairParse(currencyPair)
		if err != nil {
			return nil, fmt.Errorf("Bad currency pair format at transaction id: %s", transaction.ID.String())
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, sql.ErrNoRows
	}

	return transactions, nil
}

func (p *Postgres) GetAllTransactionsFromAccount(accountID uuid.UUID, limit int, offset int) ([]domain.Transaction, error) {
	
	query := `SELECT * FROM transactions WHERE sender_account_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`
	
	rows ,err := p.DB.Query(query, accountID, limit, offset) 
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction

	for rows.Next() {
		var transaction domain.Transaction
		var currencyPair string

		if err := rows.Scan(&transaction.ID, &transaction.SenderAccountID, &transaction.ReceiverAccountID, &transaction.Amount, &currencyPair, &transaction.CreatedAt); err != nil {
			return nil ,err		
		}

		// Set the currency pair and skip over if its corrupted (It really shouldn't be)
		transaction.CurrencyPair, err = domain.CurrencyPairParse(currencyPair)
		if err != nil {
			return nil, fmt.Errorf("Bad currency pair format at transaction id: %s", transaction.ID.String())
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, sql.ErrNoRows
	}

	return transactions, nil
}

func (p *Postgres) GetTransaction(transactionID uuid.UUID) (domain.Transaction, error) {
	query := `SELECT * FROM transactions WHERE id = $1 LIMIT 1`

	var transaction domain.Transaction
	var currencyPair string
	
	err := p.DB.QueryRow(query, transactionID).Scan(&transaction.ID, &transaction.SenderAccountID, &transaction.ReceiverAccountID, &transaction.Amount, &currencyPair, &transaction.CreatedAt)
	if err != nil {
		return domain.Transaction{}, err
	}
	
	transaction.CurrencyPair, err = domain.CurrencyPairParse(currencyPair)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("Bad currency pair format at transaction id: %s", transaction.ID.String())
	}

	return transaction, nil
}

func (p *Postgres) CreateTransaction(transaction domain.Transaction) (int64, error) {
	query := `
	INSERT INTO transactions
	(id, sender_account_id, receiver_account_id, amount, currency, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := p.DB.Exec(query, transaction.ID, transaction.SenderAccountID, transaction.ReceiverAccountID, transaction.Amount, transaction.CurrencyPair.String(), transaction.CreatedAt)
	if err != nil {
		return 0, err
	}

	return 1, nil
}