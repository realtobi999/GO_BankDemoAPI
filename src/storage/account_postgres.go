package storage

import "github.com/realtobi999/GO_BankDemoApi/src/types"

func (p *Postgres) CreateAccount(account types.Account) (int64, error) {
	query := `
	INSERT INTO accounts
	(id, customer_id, balance, account_type, currency, status, opening_date, last_transaction_date, interest_rate)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	result, err := p.DB.Exec(query, account.ID.String(),account.CustomerID.String(),account.Balance, account.Type, account.Currency, account.Status, account.OpeningDate, account.LastTransactionDate, account.InterestRate)
	if err != nil {
		return 0, err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowAffected, nil
}