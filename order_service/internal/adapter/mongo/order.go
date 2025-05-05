package mongo

import (
	"context"
	"github.com/recktt77/Microservices-First-/order_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
)

type orderRepository struct{
	dao *dao.OrderDAO
}
func NewOrderRepository(dao *dao.OrderDAO) *orderRepository{
	return &orderRepository{dao: dao}
}
func (o *orderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	_, err := o.dao.CreateOrder(ctx, order)
	return err
}

func (o *orderRepository) GetOrderByID(ctx context.Context, filter model.OrderFilter) (model.Order, error){
	if filter.ID == nil {
		return model.Order{}, model.ErrOrderNotFound
	}
	return o.dao.GetOrderByID(ctx, *filter.ID)
}

func (o *orderRepository) UpdateOrder(ctx context.Context, filter model.OrderFilter, update model.OrderUpdate) error{
	if filter.ID == nil{
		return model.ErrOrderNotFound
	}
	return o.dao.UpdateOrder(ctx, *filter.ID, update)
}

func (o *orderRepository) GetListOfOrders(ctx context.Context, filter model.OrderFilter) ([]model.Order, error){
	return o.dao.GetListOfOrders(ctx, filter)
}