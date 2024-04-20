package storage

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

func (p *Postgres) GetAllTransactions(accountID uuid.UUID, limit int, offset int) ([]types.Transaction, error){
	query := `SELECT * FROM accounts WHERE account_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`

	rows, err := p.DB.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []types.Transaction

	for rows.Next() {
		var transaction types.Transaction

		if err := rows.Scan(&transaction.ID, &transaction.SenderAccountID, &transaction.ReceiverAccountID, &transaction.Amount, &transaction.CreatedAt); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil{
        return nil, err
    }

    if len(transactions) == 0 {
        return nil, sql.ErrNoRows
    }

    return transactions, nil
}