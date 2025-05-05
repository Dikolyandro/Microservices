package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

)

type ProductHandler struct {
	Client inventorypb.ProductServiceClient
}

func NewProductHandler(client inventorypb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{Client: client}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Client.CreateProduct(c, &inventorypb.CreateProductRequest{
		Name:  req.Name,
		Price: float64(req.Price),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.Client.GetProductByID(c, &inventorypb.ProductID{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	res, err := h.Client.GetAllProducts(c, &emptypb.Empty{})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res.Products)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := h.Client.DeleteProduct(c, &inventorypb.ProductID{Id: id})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Client.UpdateProduct(c, &inventorypb.UpdateProductRequest{
		Id:    id,
		Name:  &req.Name,
		Price: &req.Price,
	})
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

