package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
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
	ErrorMessage string   `json:"error_message"`
	Code         int      `json:"code"`
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

func RespondWithJsonAndSerialize(w http.ResponseWriter, code int, payload ports.ISerializable) error {
	return RespondWithJson(w, code, payload.ToDTO())
}

func RespondWithJsonAndSerializeList[T ports.ISerializable](w http.ResponseWriter, code int, payload []T) error {
	var serializedPayload []domain.DTO

	for _, value := range payload {
		serializedPayload = append(serializedPayload, value.ToDTO())
	}

	return RespondWithJson(w, code, serializedPayload)
}


func RespondWithError(w http.ResponseWriter, code int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ErrorResponse{
		Code:         code,
		ErrorMessage: message,
	}

	log.Printf("[ERROR]\tStatus: %v Message: %s", code, message)

	return json.NewEncoder(w).Encode(response)
}

func RespondWithValidationErrors(w http.ResponseWriter, code int, message string, errs []string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := multipleErrorResponse{
		ErrorMessage: message,
		Code:         code,
		Errors:       errs,
	}

	log.Printf("[ERROR]\tStatus: %v Message: %s", code, message)

	return json.NewEncoder(w).Encode(response)
}
