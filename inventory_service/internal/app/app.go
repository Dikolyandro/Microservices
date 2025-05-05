package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/recktt77/Microservices-First-/inventory_service/config"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/grpc/server/backoffice"
	mongorepo "github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/usecase"
	mongocon "github.com/recktt77/Microservices-First-/inventory_service/pkg/mongo"
)

const (
	serviceName     = "product-service"
	shutdownTimeout = 30 * time.Second
)

type App struct {
	grpcServer backoffice.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Println(fmt.Sprintf("starting %v service", serviceName))

	mongoDB, err := mongocon.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	productDAO := dao.NewProductDAO(mongoDB.Conn)
	discountDAO := dao.NewDiscountDAO(mongoDB.Conn)

	productRepo := mongorepo.NewProductRepository(productDAO)
	discountRepo := mongorepo.NewDiscountRepository(discountDAO)

	productUC := usecase.NewProduct(productRepo)
	discountUC := usecase.NewDiscount(discountRepo, productRepo)

	var productInterface usecase.Product = *productUC
	var discountInterface usecase.Discount = *discountUC

	grpcSrv := backoffice.New(&cfg.Server, productInterface, discountInterface)

	return &App{
		grpcServer: grpcSrv,
	}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	go a.grpcServer.Run(errCh)
	log.Println(fmt.Sprintf("service %v started (gRPC)", serviceName))

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("gRPC server error: %w", err)
	case sig := <-shutdownCh:
		log.Println(fmt.Sprintf("received signal: %v. Shutting down...", sig))
		return a.Close()
	}
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.grpcServer.Stop(); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out after %v", shutdownTimeout)
	}

	log.Println("graceful shutdown completed successfully")
	return nil
}