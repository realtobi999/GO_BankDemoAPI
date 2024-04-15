package storage

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/types"
)


func (p *Postgres) GetCustomer(id uuid.UUID) (types.Customer, error) {
    query := `SELECT * FROM customers WHERE id = $1 LIMIT 1`

    var customer types.Customer

    err := p.DB.QueryRow(query, id).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Birthday, &customer.Email, &customer.Phone, &customer.State, &customer.Address, &customer.CreatedAt)
    if err != nil {
        return types.Customer{}, err
    }

    return customer, nil
}

func (p *Postgres) GetAllCustomers(limit int, offset int) ([]types.Customer, error) {
    query := `SELECT * FROM customers ORDER BY created_at LIMIT $1 OFFSET $2`

    rows, err := p.DB.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var customers []types.Customer

    for rows.Next() {
        var customer types.Customer

        if err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Birthday, &customer.Email, &customer.Phone, &customer.State, &customer.Address, &customer.CreatedAt); err != nil {
            return nil, err
        }

        customers = append(customers, customer)
    }
    
    if err := rows.Err(); err != nil{
        return nil, err
    }

    if len(customers) == 0 {
        return nil, sql.ErrNoRows
    }

    return customers, nil
}

func (p *Postgres) CreateCustomer(customer types.Customer) (int64, error) {
    query := `INSERT INTO customers (id, first_name, last_name, birthday, email, phone, state, address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

    result, err := p.DB.Exec(query, customer.ID.String(), customer.FirstName, customer.LastName, customer.Birthday, customer.Email, customer.Phone, customer.State, customer.Address, customer.CreatedAt)
    if err != nil {
        return 0, err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

    return rowsAffected, nil
}

func (p *Postgres) UpdateCustomer(customer types.Customer) error {
    query := `
    UPDATE customers
    SET first_name = $1, last_name = $2, birthday = $3, email = $4, phone = $5, state = $6, address = $7
    WHERE id = $8`

    result, err := p.DB.Exec(query, customer.FirstName, customer.LastName, customer.Birthday, customer.Email, customer.Phone, customer.State, customer.Address, customer.ID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("no rows affected")
    }

    return nil
}
