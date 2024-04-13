package handlers

import (
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
	u "github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func HealthTestHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	u.RespondWithJson(w, 200, nil)
}

func ErrorTestHandler(w http.ResponseWriter, r *http.Request, l types.ILogger, s types.IStorage) {
	u.RespondWithError(w, 500, "Something went wrong! Oops...",)
}