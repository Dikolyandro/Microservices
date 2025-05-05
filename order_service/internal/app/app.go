package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/recktt77/Microservices-First-/order_service/config"
	"github.com/recktt77/Microservices-First-/order_service/internal/adapter/grpc/server/backoffice"
	mongorepo "github.com/recktt77/Microservices-First-/order_service/internal/adapter/mongo"
	"github.com/recktt77/Microservices-First-/order_service/internal/adapter/mongo/dao"
	"github.com/recktt77/Microservices-First-/order_service/internal/usecase"
	mongocon "github.com/recktt77/Microservices-First-/order_service/pkg/mongo"
)

const (
	serviceName     = "order-service"
	shutdownTimeout = 30 * time.Second
)

type App struct {
	grpcServer backoffice.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Println(fmt.Sprintf("starting %v service", serviceName))

	log.Println("connecting to mongo", "database", cfg.Mongo.Database)
	mongoDB, err := mongocon.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	// DAO
	orderDAO := dao.NewOrderDAO(mongoDB.Conn)

	// Repository
	orderRepo := mongorepo.NewOrderRepository(orderDAO)

	// UseCase
	orderUsecase := usecase.NewOrder(orderRepo)

	var orderInterface usecase.Order = *orderUsecase

	// gRPC Service
	grpcServer := backoffice.New(&cfg.Server, orderInterface)

	app := &App{
		grpcServer: grpcServer,
	}

	return app, nil
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.grpcServer.Stop(); err != nil {
		return fmt.Errorf("failed to shutdown gRPC service: %w", err)
	}

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out after %v", shutdownTimeout)
	}

	log.Println("graceful shutdown completed successfully")
	return nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	a.grpcServer.Run(errCh)
	log.Println(fmt.Sprintf("service %v started", serviceName))

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return fmt.Errorf("service error: %w", errRun)

	case s := <-shutdownCh:
		log.Println(fmt.Sprintf("received signal: %v. Running graceful shutdown...", s))
		if err := a.Close(); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
	}

	return nil
}
