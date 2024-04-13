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

	// Initiate the logger
	logger := logs.NewLogger(`src\utils\logs\logs.txt`)

	//Initiate the database
	storage, err := storage.NewPostgres("localhost", "5432", "postgres", "root", "GoBank", "disable")
	if err != nil {
		logger.Fatal(err)
	}
	logger.LogEvent("Database is successfully connected")

	server := api.NewServer(port, storage, logger)

	logger.LogEvent(fmt.Sprintf("Server successfully started on port : %v", server.Port))
	log.Fatal(server.Run())
}