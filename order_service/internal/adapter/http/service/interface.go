package service

import (
	"context"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	GetOrderByID(ctx context.Context, filter model.OrderFilter) (model.Order, error)
	UpdateOrder(ctx context.Context, order model.OrderUpdate) (model.Order, error)
	GetListOfOrders(ctx context.Context, filter model.OrderFilter) ([]model.Order, error)
}
