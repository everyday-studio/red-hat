package handlers

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrResponse(err error) ErrorResponse {
	if err == nil {
		return ErrorResponse{Error: "unknown error"}
	}
	return ErrorResponse{Error: errors.Unwrap(err).Error()}
}
