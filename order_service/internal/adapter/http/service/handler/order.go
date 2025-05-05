package handler

import (
	_"bytes"
	_"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/recktt77/Microservices-First-/order_service/internal/adapter/http/service/handler/dto"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	uc OrderUsecase
}

func NewOrder(uc OrderUsecase) *Order {
	return &Order{
		uc: uc,
	}
}

// func decreaseStock(productId primitive.ObjectID, qty int) error {
//     req := dto.OrderedProduct{
//         ProductId: productId,
//         Quantity:  qty,
//     }

//     data, _ := json.Marshal(req)
//     resp, err := http.Post(
//         "http://localhost:8081/api/v1/products/",
//         "application/json",
//         bytes.NewBuffer(data),
//     )
//     if err != nil {
//         return err
//     }
//     defer resp.Body.Close()

//     if resp.StatusCode != http.StatusOK {
//         return fmt.Errorf("inventory failed to update stock")
//     }

//     return nil
// }

func (h *Order) CreateOrder(ctx *gin.Context) {
	order, err := dto.FromOrderCreateRequest(ctx)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	newOrder, err := h.uc.CreateOrder(ctx.Request.Context(), order)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToOrderCreateResponse(newOrder))
}


func (h *Order) GetOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Printf("Received GetByID request with ID string: %s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Failed to convert ID string to ObjectID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}
	log.Printf("Successfully converted ID string to ObjectID: %s", objectID.Hex())

	order, err := h.uc.GetOrderByID(ctx.Request.Context(), model.OrderFilter{
		ID: &objectID,
	})

	if err != nil {
		errCtx := dto.FromError(err)
		log.Printf("Error retrieving order: %v", errCtx.Message)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	log.Printf("Successfully retrieved order with ID: %s", order.ID.Hex())
	ctx.JSON(http.StatusOK, dto.ToOrderCreateResponse(order))
}

func (h *Order) UpdateOrder(ctx *gin.Context) {
	idParam := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}

	updateData, err := dto.FromOrderUpdateRequest(ctx)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	updateData.ID = &objectID

	updated, err := h.uc.UpdateOrder(ctx.Request.Context(), updateData)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToOrderCreateResponse(updated))
}


func (h *Order) GetListOfOrders(ctx *gin.Context) {
	orders, err := h.uc.GetListOfOrders(ctx.Request.Context(), model.OrderFilter{})
	fmt.Print(orders)
	if err != nil {
		errCtx := dto.FromError(err)
		ctx.JSON(errCtx.Code, gin.H{"error": errCtx.Message})
		return
	}

	var response []dto.OrderCreateResponse
	for _, o := range orders {
		response = append(response, dto.ToOrderCreateResponse(o))
	}

	ctx.JSON(http.StatusOK, response)
}
