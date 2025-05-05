package model

import (
	"errors"
)
var(
	ErrInactiveOrder = errors.New("order is inactive")
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderAlreadyPaid = errors.New("order already paid")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrOrderAlreadyReturned = errors.New("order already returned")
	ErrOrderAlreadyCompleted = errors.New("order already completed")
	ErrOrderAlreadyProcessing = errors.New("order already processing")
	ErrOrderAlreadyFailed = errors.New("order already failed")
	ErrOrderAlreadyShipped = errors.New("order already shipped")
	ErrInvalidProducts = errors.New("Invalid products")
	ErrInvalidStatus = errors.New("Invalid status")
)