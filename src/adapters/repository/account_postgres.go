package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)

func (p *Postgres) GetAllAccounts(limit int, offset int) ([]domain.Account, error) {
	query := `SELECT * FROM accounts ORDER BY created_at LIMIT $1 OFFSET $2`

	rows, err := p.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account

	for rows.Next() {
		var account domain.Account

		if err := rows.Scan(&account.ID, &account.CustomerID, &account.Balance, &account.Type, &account.Currency, &account.Status, &account.OpeningDate, &account.LastTransactionDate, &account.InterestRate, &account.CreatedAt); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, sql.ErrNoRows
	}

	return accounts, nil

}

func (p *Postgres) GetAllAccountsByCustomer(customerID uuid.UUID, limit int, offset int) ([]domain.Account, error) {
	query := `SELECT * FROM accounts WHERE customer_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`

	rows, err := p.DB.Query(query, customerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account

	for rows.Next() {
		var account domain.Account

		if err := rows.Scan(&account.ID, &account.CustomerID, &account.Balance, &account.Type, &account.Currency, &account.Status, &account.OpeningDate, &account.LastTransactionDate, &account.InterestRate, &account.CreatedAt); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, sql.ErrNoRows
	}

	return accounts, nil

}

func (p *Postgres) GetAllSavingsAccounts() ([]domain.Account, error) {
	query := `SELECT * FROM accounts WHERE account_type = 3 ORDER BY created_at`

	rows, err := p.DB.Query(query, )
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account

	for rows.Next() {
		var account domain.Account

		if err := rows.Scan(&account.ID, &account.CustomerID, &account.Balance, &account.Type, &account.Currency, &account.Status, &account.OpeningDate, &account.LastTransactionDate, &account.InterestRate, &account.CreatedAt); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, sql.ErrNoRows
	}

	return accounts, nil
}

func (p *Postgres) GetAccount(accountID uuid.UUID) (domain.Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1`

	var account domain.Account

	err := p.DB.QueryRow(query, accountID).Scan(&account.ID, &account.CustomerID, &account.Balance, &account.Type, &account.Currency, &account.Status, &account.OpeningDate, &account.LastTransactionDate, &account.InterestRate, &account.CreatedAt)
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

func (p *Postgres) GetAccountByOwner(customerID, accountID uuid.UUID) (domain.Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1 AND customer_id = $2 LIMIT 1`

	var account domain.Account

	err := p.DB.QueryRow(query, accountID, customerID).Scan(&account.ID, &account.CustomerID, &account.Balance, &account.Type, &account.Currency, &account.Status, &account.OpeningDate, &account.LastTransactionDate, &account.InterestRate, &account.CreatedAt)
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil
}

func (p *Postgres) CreateAccount(account domain.Account) (int64, error) {
	query := `
	INSERT INTO accounts
	(id, customer_id, balance, account_type, currency, status, opening_date, last_transaction_date, interest_rate, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := p.DB.Exec(query, account.ID.String(), account.CustomerID.String(), account.Balance, account.Type, account.Currency, account.Status, account.OpeningDate, account.LastTransactionDate, account.InterestRate, account.CreatedAt)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (p *Postgres) UpdateAccount(account domain.Account) (int64, error) {
	query := `
	UPDATE accounts
	SET balance = $1, account_type = $2, currency = $3, status = $4, last_transaction_date = $5, interest_rate = $6
	WHERE id = $7
	`

	result, err := p.DB.Exec(query, account.Balance, account.Type, account.Currency, account.Status, account.LastTransactionDate, account.InterestRate, account.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (p *Postgres) DeleteAccount( accountID uuid.UUID) (int64, error) {
	query := `DELETE FROM accounts WHERE id = $1`

	result, err := p.DB.Exec(query, accountID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}