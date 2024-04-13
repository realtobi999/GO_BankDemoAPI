package utils

import (
	"encoding/json"
	"net/http"
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

type multipleErrorResponse struct {
	Message string  `json:"message"`
	Status  int     `json:"status"`
	Errors  []string `json:"errors"`
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

func RespondWithError(w http.ResponseWriter, code int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ErrorResponse{
		Code:         code,
		ErrorMessage: message,
	}

	return json.NewEncoder(w).Encode(response)
}

func RespondWithValidationErrors(w http.ResponseWriter, code int, message string, errs []error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Convert errors to string slice
	var errorStrings []string
	for _, err := range errs {
		errorStrings = append(errorStrings, err.Error())
	}

	response := multipleErrorResponse{
		Message: message,
		Status:  code,
		Errors:  errorStrings,
	}

	return json.NewEncoder(w).Encode(response)
}
