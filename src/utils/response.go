package utils

import (
	"encoding/json"
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/types"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	Code         int    `json:"code"`
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := SuccessResponse{
		Code:    code,
		Message: "Success, everything is fine!",
		Data:    payload,
	}

	return json.NewEncoder(w).Encode(response)
}

func RespondWithError(w http.ResponseWriter, logger types.ILogger, code int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ErrorResponse{
		Code:         code,
		ErrorMessage: message,
	}

	logger.LogError(message)

	return json.NewEncoder(w).Encode(response)
}
