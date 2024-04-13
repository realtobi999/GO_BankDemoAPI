package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(host, port, user, password, dbname, sslmode string) (*Postgres, error) {
	connectionString := fmt.Sprintf(
		"host=%v port=%v user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) DatabaseHas(table, column string, value any) bool{
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1", column, table, column)

	var result sql.NullString
	err := p.db.QueryRow(query, value).Scan(&result)

	return err == nil && result.Valid
}

