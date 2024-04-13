package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/realtobi999/GO_BankDemoApi/src/api"
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

	server := api.NewServer(port)

	log.Printf("[*] - Server successfully started on port : %v\n", server.Port)
	log.Fatal(server.Run())
}