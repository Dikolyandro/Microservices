package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DiscountHandler struct {
	Client inventorypb.DiscountServiceClient
}

func NewDiscountHandler(client inventorypb.DiscountServiceClient) *DiscountHandler {
	return &DiscountHandler{Client: client}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var req struct {
		Name                 string   `json:"name"`
		Description          string   `json:"description"`
		DiscountPercentage   float64  `json:"discount_percentage"`
		ApplicableProductIDs []string `json:"applicable_product_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount := &inventorypb.CreateDiscountRequest{
		Name:                 req.Name,
		Description:          req.Description,
		DiscountPercentage:   req.DiscountPercentage,
		ApplicableProductIds: req.ApplicableProductIDs,
	}

	res, err := h.Client.CreateDiscount(c, discount)
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *DiscountHandler) GetAllDiscounts(c *gin.Context) {
	res, err := h.Client.GetAllDiscounts(c, &emptypb.Empty{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res.Discounts)
}

func (h *DiscountHandler) GetProductsWithDiscounts(c *gin.Context) {
	res, err := h.Client.GetProductsWithDiscounts(c, &emptypb.Empty{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res.Products)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	id := c.Param("id")
	_, err := h.Client.DeleteDiscount(c, &inventorypb.DiscountID{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
