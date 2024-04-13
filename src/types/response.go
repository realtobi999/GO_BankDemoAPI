package types

type SuccessResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	Code         int    `json:"code"`
}