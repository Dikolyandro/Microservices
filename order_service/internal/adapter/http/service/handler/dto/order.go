package dto

import (
	"net/http"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderedProduct struct {
	ProductId primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}

type OrderCreateRequest struct {
	UserID   primitive.ObjectID           `json:"user_id"`
	Products []OrderedProduct `json:"products"`
	Status string `json:"status"`
}

type OrderUpdateRequest struct {
	Status *string `json:"status"`
}

type OrderCreateResponse struct {
	ID        primitive.ObjectID           `json:"id"`
	UserID    primitive.ObjectID           `json:"user_id"`
	Products  []OrderedProduct `json:"products"`
	Status    string           `json:"status"`
	CreatedAt string           `json:"created_at"`
}

type OrderUpdateResponse = OrderCreateResponse
type OrderGetResponse = OrderCreateResponse


func FromOrderCreateRequest(ctx *gin.Context) (model.Order, error) {
	var req OrderCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return model.Order{}, err
	}

	var modelProducts []model.OrderedProduct
	for _, p := range req.Products {
		modelProducts = append(modelProducts, model.OrderedProduct{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		})
	}

	now := time.Now()

	return model.Order{
		UserID:    req.UserID,
		Products:  modelProducts,
		Status:    req.Status,
		CreatedAt: now,
		UpdatedAt: now,
		IsDeleted: false,
	}, nil
}


func FromOrderUpdateRequest(ctx *gin.Context) (model.OrderUpdate, error) {
	var req OrderUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.OrderUpdate{}, err
	}

	return model.OrderUpdate{
		Status:    req.Status,
		UpdatedAt: pointerToTime(time.Now()),
	}, nil
}


func ToOrderCreateResponse(order model.Order) OrderCreateResponse {
	var dtoProducts []OrderedProduct
	for _, p := range order.Products {
		dtoProducts = append(dtoProducts, OrderedProduct{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		})
	}

	return OrderCreateResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Products:  dtoProducts,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}
}



func pointerToTime(t time.Time) *time.Time {
	return &t
}
