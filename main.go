package main

import (
	"log"

	"github.com/realtobi999/GO_BankDemoApi/src/api"
	u "github.com/realtobi999/GO_BankDemoApi/src/utils"
)

const port int = 8080;

func main() {
	u.ClearConsole()
	u.PrintASCII()

	server := api.NewServer(port)

	log.Printf("[*] - Server successfully started on port : %v\n", server.Port)
	log.Fatal(server.Run())
}