package usecase

import (
	"context"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	repo ProductRepo
}

func NewProduct(repo ProductRepo) *Product {
	return &Product{
		repo: repo,
	}
}

func (p *Product) Create(ctx context.Context, request model.Product) (model.Product, error) {
	if err := p.repo.Create(ctx, request); err != nil {
		return model.Product{}, err
	}

	return request, nil
}

func (p *Product) GetByID(ctx context.Context, request model.ProductFilter) (model.Product, error) {
	if request.ID == nil {
		return model.Product{}, model.ErrProductNotFound
	}

	product, err := p.repo.GetByID(ctx, request)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (p *Product) Update(ctx context.Context, request model.ProductUpdate) (model.Product, error) {
	if request.ID == nil {
		return model.Product{}, model.ErrProductNotFound
	}

	if err := p.repo.Update(ctx, model.ProductFilter{ID: request.ID}, request); err != nil {
		return model.Product{}, err
	}

	updated, err := p.repo.GetByID(ctx, model.ProductFilter{ID: request.ID})
	if err != nil {
		return model.Product{}, err
	}

	return updated, nil
}

func (p *Product) Delete(ctx context.Context, request model.ProductFilter) error {
	if request.ID == nil {
		return model.ErrProductNotFound
	}
	return p.repo.Delete(ctx, request)
}

func (p *Product) GetAll(ctx context.Context, filter model.ProductFilter) ([]model.Product, error) {
	return p.repo.GetAll(ctx, filter)
}

func (p *Product) GetByIDs(ctx context.Context, ids[] primitive.ObjectID) ([]model.Product, error) {
	return p.repo.GetByIDs(ctx, ids)
}