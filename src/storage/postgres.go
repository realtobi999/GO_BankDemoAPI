package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
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

	return &Postgres{DB: db}, nil
}

func (p *Postgres) DatabaseHas(table, column string, value any) bool{
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1", column, table, column)

	var result sql.NullString
	err := p.DB.QueryRow(query, value).Scan(&result)

	return err == nil && result.Valid
}

func (p *Postgres) ClearAllTables() error {
	// Query to retrieve table names from information_schema.tables
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
	`

	// Execute the query to retrieve table names
	rows, err := p.DB.Query(query)
	if err != nil {
		return fmt.Errorf("failed to retrieve table names: %v", err)
	}
	defer rows.Close()

	// Truncate each table
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return fmt.Errorf("failed to scan table name: %v", err)
		}

		// Truncate the table
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)
		_, err = p.DB.Exec(truncateQuery)
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %v", tableName, err)
		}
	}

	return nil
}

