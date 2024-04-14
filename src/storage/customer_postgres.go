package storage

import (
	"errors"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
	"github.com/realtobi999/GO_BankDemoApi/src/utils/custom_errors"
)


func (p *Postgres) GetAllCustomers(limit int, offset int) ([]types.Customer, error) {
    query := `SELECT * FROM customers LIMIT $1 OFFSET $2`

    rows, err := p.DB.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var customers []types.Customer

    for rows.Next() {
        var customer types.Customer

        if err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Birthday, &customer.Email, &customer.Phone, &customer.State, &customer.Address); err != nil {
            return nil, err
        }

        customers = append(customers, customer)
    }
    
    if err := rows.Err(); err != nil{
        return nil, err
    }

    if len(customers) == 0 {
        return nil, errors.New(custom_errors.StorageNoResultsFound)
    }

    return customers, nil
}

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
