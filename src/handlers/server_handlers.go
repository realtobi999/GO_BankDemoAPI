package handlers

import (
	"net/http"

	u "github.com/realtobi999/GO_BankDemoApi/src/utils"
)

func HealthTestHandler(w http.ResponseWriter, r *http.Request) {
	u.RespondWithJson(w, 200, nil)
}

func ErrorTestHandler(w http.ResponseWriter, r *http.Request) {
	u.RespondWithError(w, 500, "Something went wrong! Oops...")
}