package main

import (
	"github.com/gin-gonic/gin"

	"github.com/recktt77/Microservices-First-/api-gateway/internal/grpc"
	"github.com/recktt77/Microservices-First-/api-gateway/internal/handler"
)

func main() {
	r := gin.Default()

	clients := grpc.NewClients(
		"localhost:8081", 
		"localhost:8082", 
	)

	productHandler := handler.NewProductHandler(clients.Inventory)
	discountHandler := handler.NewDiscountHandler(clients.Discount)
	orderHandler := handler.NewOrderHandler(clients.Order)

	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetAllProducts)
	r.GET("/products/:id", productHandler.GetProductByID)
	r.PATCH("/products/:id", productHandler.UpdateProduct)
	r.DELETE("/products/:id", productHandler.DeleteProduct)

	r.POST("/discounts", discountHandler.CreateDiscount)
	r.GET("/discounts", discountHandler.GetAllDiscounts)
	r.GET("/discounts/products", discountHandler.GetProductsWithDiscounts)
	r.DELETE("/discounts/:id", discountHandler.DeleteDiscount)

	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders", orderHandler.ListOrders)
	r.GET("/orders/:id", orderHandler.GetOrder)
	r.PATCH("/orders/:id", orderHandler.UpdateOrder)

	r.Run(":8080")
}
