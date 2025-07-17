// internal/shared/model/response.go

package model

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

type ErrorResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}
