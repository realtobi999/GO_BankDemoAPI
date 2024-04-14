package handlers

import (
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

func HealthTestHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	RespondWithJson(w, 200, nil)
}

func ErrorTestHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	RespondWithError(w, 500, "Something went wrong! Oops...",)
}