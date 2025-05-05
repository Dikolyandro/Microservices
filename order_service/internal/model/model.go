package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID
	UserID    primitive.ObjectID
	Products  []OrderedProduct
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}


type OrderedProduct struct {
	ProductId primitive.ObjectID
	Quantity  int
}

type OrderFilter struct {
	ID        *primitive.ObjectID
	UserID    *primitive.ObjectID
	Products  *[]OrderedProduct
	Status    *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	IsDeleted *bool
}

type OrderUpdate struct {
	ID        *primitive.ObjectID
	UserID    *primitive.ObjectID
	Products  *[]OrderedProduct
	Status    *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	IsDeleted *bool
}

