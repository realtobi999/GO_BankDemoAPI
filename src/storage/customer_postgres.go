package storage

import "github.com/realtobi999/GO_BankDemoApi/src/types"

func (p *Postgres) CreateCustomer(customer types.Customer) (int64, error) {
    query := `INSERT INTO customers (id, first_name, last_name, birthday, email, phone, state, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

    result, err := p.DB.Exec(query, customer.ID.String(), customer.FirstName, customer.LastName, customer.Birthday, customer.Email, customer.Phone, customer.State, customer.Address)
    if err != nil {
        return 0, err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

    return rowsAffected, nil
}
