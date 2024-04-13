package main

import (
	"log"

	"github.com/realtobi999/GO_BankDemoApi/src/api"
)

const port int = 8080;

func main() {
	server := api.NewServer(port)

	log.Printf("[*] - Server successfully started on port : %v\n", server.Port)
	log.Fatal(server.Run())
}