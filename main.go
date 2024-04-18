package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/realtobi999/GO_BankDemoApi/src/api"
	"github.com/realtobi999/GO_BankDemoApi/src/storage"
	"github.com/realtobi999/GO_BankDemoApi/src/utils/logs"
)
func main() {
	clearConsole()
	printASCII()

	
	// Load the .env file containing Environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[Error] - Error loading .env file")
	}
	
	// Get the server port
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("[Error] - Error parsing port from the .env file")
	}
	
	// Get database configuration
	dbConfig := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"username": os.Getenv("DB_USERNAME"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbName":   os.Getenv("DB_NAME"),
	}
	
	// Initiate the logger
	logger := logs.NewLogger(`src\utils\logs\logs.txt`)
	
	// Initiate the database
	database, err := storage.NewPostgres(dbConfig["host"], dbConfig["port"], dbConfig["username"], dbConfig["password"], dbConfig["dbName"], "disable")
	if err != nil {
		logger.Fatal(err)
	}
	logger.LogEvent("Database is successfully connected!")
	
	// Run migrations
	if err := storage.RunMigrations(storage.PathToMigrations,database.DB, logger); err != nil {
		logger.Fatal(err)
	}
	logger.LogEvent("Migration are successfully inserted!")

	server := api.NewServer(port, database, logger)

	logger.LogEvent(fmt.Sprintf("Server successfully started on port : %v", server.Port))
	log.Fatal(server.Run())
}