package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderpb "github.com/recktt77/proto-definitions/gen/orders"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	Client orderpb.OrderServiceClient
}

func NewOrderHandler(client orderpb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{Client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id"`
		Status  string `json:"status"`
		Products []struct {
			ProductID string `json:"product_id"`
			Quantity  int32  `json:"quantity"`
		} `json:"products"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderedProducts []*orderpb.OrderedProduct
	for _, p := range req.Products {
		orderedProducts = append(orderedProducts, &orderpb.OrderedProduct{
			ProductId: p.ProductID,
			Quantity:  p.Quantity,
		})
	}

	res, err := h.Client.CreateOrder(c, &orderpb.CreateOrderRequest{
		UserId:  req.UserID,
		Status:  req.Status,
		Products: orderedProducts,
	})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	res, err := h.Client.GetOrder(c, &orderpb.GetOrderRequest{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Client.UpdateOrder(c, &orderpb.UpdateOrderRequest{
		Id:     id,
		Status: req.Status,
	})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	res, err := h.Client.ListOrders(c, &orderpb.ListOrdersRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res.Orders)
}
