package main

import (
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/repository"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/repository/migrations"
	"github.com/realtobi999/GO_BankDemoApi/src/adapters/web"
	"github.com/realtobi999/GO_BankDemoApi/src/core/services/account"
	"github.com/realtobi999/GO_BankDemoApi/src/core/services/customer"
)

func main() {
	clearConsole()
	ASCII()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("[Error] - Error loading .env file")
	}

	// Get database configuration
	dbConfig := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"username": os.Getenv("DB_USERNAME"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbName":   os.Getenv("DB_NAME"),
	}

	database, err := repository.NewPostgres(dbConfig["host"], dbConfig["port"], dbConfig["username"], dbConfig["password"], dbConfig["dbName"], "disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := migrations.RunMigrations("src/adapters/repository/migrations/*.sql", database.DB); err != nil {
		log.Fatal(err)
	}

	server := web.NewServer(":8080",  chi.NewMux())
	server.AccountService = account.NewAccountService(database)
	server.CustomerService = customer.NewCustomerService(database)

	server.LoadSharedMiddleware()
	server.LoadRoutes()
	
	defer log.Printf("[EVENT]\tShuting down...")
	log.Fatal(server.Run())
}
