package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

const PathToMigrations string = "src/storage/migrations/*.sql"

func RunMigrations(path string, db *sql.DB, logger types.ILogger) error {
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {

		// Read migration file
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}
		sqlQuery := string(sqlBytes)

		// Execute migration
		_, err = db.Exec(sqlQuery)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
		logger.LogEvent(fmt.Sprintf("Applied migration: %s", file))
	}

	return nil
}
