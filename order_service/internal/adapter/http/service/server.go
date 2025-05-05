package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"github.com/recktt77/Microservices-First-/order_service/config"
	"github.com/recktt77/Microservices-First-/order_service/internal/adapter/http/service/handler"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const serverIPAddress = "0.0.0.0:%d"

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string
	
	orderHandler *handler.Order
}

func (a *API) setupRoutes(){
	v1 := a.server.Group("/api/v1")
	{
		order := v1.Group("/orders")
		{
			order.POST("/", a.orderHandler.CreateOrder)
			order.GET("/:id", a.orderHandler.GetOrderByID)
			order.PATCH("/:id", a.orderHandler.UpdateOrder)
			order.GET("/", a.orderHandler.GetListOfOrders)
		}
	}
	a.server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func New(cfg config.Server, orderUsecase OrderUsecase) *API {
	//setting the Gin mode
	gin.SetMode(cfg.HTTPServer.Mode)
	// creating a new gin engine
	server := gin.New()

	//applying the middleware
	server.Use(gin.Recovery())

	orderHandler := handler.NewOrder(orderUsecase)

	api := &API{
		server:         server,
		cfg:            cfg.HTTPServer,
		addr:           fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
		orderHandler: orderHandler,
	}

	api.setupRoutes()
	return api
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		log.Printf("HTTP server starting on: %v", a.addr)
		if err := a.server.Run(a.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}

func (a *API) Stop() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Shutting signal received: %v", sig.String())

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully")

	log.Println("HTTP server stopped successfully")

	return nil
}