package usecase

import (
	"context"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
)

type Order struct {
	repo OrderRepo
}

func NewOrder(repo OrderRepo) *Order {
	return &Order{
		repo: repo,
	}
}

func (o *Order) CreateOrder(ctx context.Context, request model.Order) (model.Order, error) {
	err := o.repo.CreateOrder(ctx, request)
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{
		ID:       request.ID,
		UserID:   request.UserID,
		Products: request.Products,
		Status:   request.Status,
	}, nil
}

func (o *Order) GetOrderByID(ctx context.Context, request model.OrderFilter) (model.Order, error) {
	if request.ID == nil {
		return model.Order{}, model.ErrOrderNotFound
	}
	order, err := o.repo.GetOrderByID(ctx, request)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (o *Order) UpdateOrder(ctx context.Context, request model.OrderUpdate) (model.Order, error) {
	if request.ID == nil {
		return model.Order{}, model.ErrOrderNotFound
	}
	err := o.repo.UpdateOrder(ctx, model.OrderFilter{ID: request.ID}, request)
	if err != nil {
		return model.Order{}, err
	}

	order, err := o.repo.GetOrderByID(ctx, model.OrderFilter{ID: request.ID})
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (o *Order) GetListOfOrders(ctx context.Context, request model.OrderFilter) ([]model.Order, error) {
	orders, err := o.repo.GetListOfOrders(ctx, request)
	if err != nil {
		return []model.Order{}, err
	}

	return orders, nil
}
