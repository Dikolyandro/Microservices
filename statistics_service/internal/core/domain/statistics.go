package domain

import "time"

// OrderStat represents a user's order for statistics.
type OrderStat struct {
	UserID    string    `bson:"user_id"`
	OrderID   string    `bson:"order_id"`
	CreatedAt time.Time `bson:"created_at"`
}

// InventoryStat represents inventory changes for statistics.
type InventoryStat struct {
	ProductID string    `bson:"product_id"`
	Action    string    `bson:"action"` // e.g. created, updated, deleted
	Time      time.Time `bson:"time"`
}
