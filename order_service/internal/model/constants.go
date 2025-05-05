package model

const (
	OrderStatusPending   = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
	OrderStatusFailed     = "failed"
	OrderStatusReturned    = "returned"

	DefaultLimit = 20
	DefaultOffset = 0
)