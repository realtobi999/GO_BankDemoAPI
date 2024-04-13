package handlers

import (
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
	u "github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func HealthTestHandler(w http.ResponseWriter, r *http.Request, l types.ILogger) {
	u.RespondWithJson(w, 200, nil)
}

func ErrorTestHandler(w http.ResponseWriter, r *http.Request, l types.ILogger) {
	u.RespondWithError(w, l, 500, "Something went wrong! Oops...",)
}