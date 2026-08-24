package responses

import "github.com/google/uuid"

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TaskResponse struct {
	Id          uuid.UUID `json:"id"`
	Tag         string    `json:"tag"`
	Description string    `json:"description"`
}
