package api

import (
	"encoding/json"
	"fmt"
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

type multipleErrorResponse struct {
	ErrorMessage string   `json:"message"`
	Code         int      `json:"status"`
	Errors       []string `json:"errors"`
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

func RespondWithJsonAndSerialize(w http.ResponseWriter, code int, payload types.ISerializable) error {
	return RespondWithJson(w, code, payload.ToDTO())
}

func RespondWithJsonAndSerializeList[T types.ISerializable](w http.ResponseWriter, code int, payload []T) error {
	var serializedPayload []types.DTO

	for _, value := range payload {
		serializedPayload = append(serializedPayload, value.ToDTO())
	}

	return RespondWithJson(w, code, serializedPayload)
}

func RespondWithError(w http.ResponseWriter, logger types.ILogger, code int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ErrorResponse{
		Code:         code,
		ErrorMessage: message,
	}

	logger.LogError(fmt.Sprintf("Code: %v Message: %s", code, message))

	return json.NewEncoder(w).Encode(response)
}

func RespondWithValidationErrors(w http.ResponseWriter, logger types.ILogger, code int, message string, errs []error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Convert errors to string slice
	var errorStrings []string
	for _, err := range errs {
		errorStrings = append(errorStrings, err.Error())
	}

	response := multipleErrorResponse{
		ErrorMessage: message,
		Code:         code,
		Errors:       errorStrings,
	}

	logger.LogError(fmt.Sprintf("Code: %v Message: %s", code, message))

	return json.NewEncoder(w).Encode(response)
}
