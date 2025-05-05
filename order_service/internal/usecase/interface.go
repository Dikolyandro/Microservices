package usecase

import (
	"context"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrderByID(ctx context.Context, filter model.OrderFilter) (model.Order, error)
	UpdateOrder(ctx context.Context, filter model.OrderFilter, update model.OrderUpdate) error
	GetListOfOrders(ctx context.Context, filter model.OrderFilter) ([]model.Order, error)
}