package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(path string, db *sql.DB) error {
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		// Read migration file
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("[ERROR]\tFailed to read migration file %s: %w", file, err)
		}
		sqlQuery := string(sqlBytes)

		// Execute migration
		_, err = db.Exec(sqlQuery)
		if err != nil {
			return fmt.Errorf("[ERROR]\tFailed to execute migration %s: %w", file, err)
		}
		log.Printf("[EVENT]\tApplied migration: %s", file)
	}

	log.Printf("[EVENT]\tMigration were successfully inserted")

	return nil
}

// DropMigrations drops all tables in the specified database.
func DropMigrations(db *sql.DB) error {
	// Query to select all table names in the public schema
	query := `
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public'
    `

	// Retrieve table names
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query table names: %v", err)
	}
	defer rows.Close()

	// Iterate over each table and drop it
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %v", err)
		}

		dropStmt := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)
		if _, err := db.Exec(dropStmt); err != nil {
			log.Printf("Error dropping table %s: %v", tableName, err)
			// Continue dropping other tables even if one fails
		} else {
			log.Printf("Dropped table %s", tableName)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating table names: %v", err)
	}

	return nil
}