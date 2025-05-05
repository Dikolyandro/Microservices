package dto

import (
	"errors"
	"net/http"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
)

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrInvalidProducts = &HTTPError{
		Code: http.StatusBadRequest,
		Message: "Invalid products",
	}

	ErrInvalidStatus = &HTTPError{
		Code: http.StatusBadRequest,
		Message: "invalid status",
	}

	ErrOrderNotFound = &HTTPError{
		Code: http.StatusBadRequest,
		Message: "order not found",
	}
)

func FromError(err error) *HTTPError{
	switch{
	case errors.Is(err, model.ErrInvalidProducts):
		return ErrInvalidProducts
	case errors.Is(err, model.ErrInvalidStatus):
		return ErrInvalidStatus
	case errors.Is(err, model.ErrOrderNotFound):
		return ErrOrderNotFound
	default:
		return &HTTPError{
			Code: http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
}