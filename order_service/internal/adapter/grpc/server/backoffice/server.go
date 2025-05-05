package backoffice

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/recktt77/Microservices-First-/order_service/config"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
	"github.com/recktt77/Microservices-First-/order_service/internal/usecase"
	orderpb "github.com/recktt77/proto-definitions/gen/orders"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type server struct {
	orderpb.UnimplementedOrderServiceServer

	cfg         *config.Server
	orderUC     usecase.Order
	grpcServer  *grpc.Server
}

func New(cfg *config.Server, orderUC usecase.Order) *server {
	s := &server{
		cfg:        cfg,
		orderUC:    orderUC,
		grpcServer: grpc.NewServer(),
	}

	orderpb.RegisterOrderServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)
	return s
}

func (s *server) Run(errCh chan<- error) {
	addr := fmt.Sprintf(":%d", s.cfg.HTTPServer.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		errCh <- err
		return
	}
	log.Println("gRPC order service started on", addr)
	if err := s.grpcServer.Serve(lis); err != nil {
		errCh <- err
	}
}

func (s *server) Stop() error {
	log.Println("Stopping gRPC order server...")
	s.grpcServer.GracefulStop()
	return nil
}

func (s *server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	var products []model.OrderedProduct
	for _, p := range req.Products {
		id, err := primitive.ObjectIDFromHex(p.ProductId)
		if err != nil {
			return nil, err
		}

		products = append(products, model.OrderedProduct{
			ProductId: id,
			Quantity:  int(p.Quantity),
		})

	}

	userId, err := primitive.ObjectIDFromHex(req.UserId)
		if err != nil {
			return nil, err
		}

	order := model.Order{
		UserID:   userId,
		Status:   req.Status,
		Products: products,
	}

	created, err := s.orderUC.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	var protoProducts []*orderpb.OrderedProduct

	for _, p := range created.Products {
		protoProducts = append(protoProducts, &orderpb.OrderedProduct{
			ProductId: p.ProductId.Hex(),
			Quantity:  int32(p.Quantity),
		})
	}

	return &orderpb.CreateOrderResponse{
		Id:        created.ID.Hex(),
		UserId:    created.UserID.Hex(),
		Products:  protoProducts,
		Status:    created.Status,
		CreatedAt: timestamppb.New(created.CreatedAt),
	}, nil
}

func (s *server) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	order, err := s.orderUC.UpdateOrder(ctx, model.OrderUpdate{
		ID:     &id,
		Status: &req.Status,
	})
	if err != nil {
		return nil, err
	}

	var protoProducts []*orderpb.OrderedProduct
	for _, p := range order.Products {
		protoProducts = append(protoProducts, &orderpb.OrderedProduct{
			ProductId: p.ProductId.Hex(),
			Quantity:  int32(p.Quantity),
		})
	go}

	
	return &orderpb.UpdateOrderResponse{
		Id:        order.ID.Hex(),
		UserId:    order.UserID.Hex(),
		Products:  protoProducts,
		Status:    order.Status,
		CreatedAt: timestamppb.New(order.CreatedAt),
	}, nil
}

func (s *server) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	filter := model.OrderFilter{
		ID: &id,
	}
	order, err := s.orderUC.GetOrderByID(ctx, filter)
	if err != nil {
		return nil, err
	}

	var protoProducts []*orderpb.OrderedProduct
	for _, p := range order.Products {
		protoProducts = append(protoProducts, &orderpb.OrderedProduct{
			ProductId: p.ProductId.Hex(),
			Quantity:  int32(p.Quantity),
		})
	}

	return &orderpb.GetOrderResponse{
		Id:        order.ID.Hex(),
		UserId:    order.UserID.Hex(),
		Products:  protoProducts,
		Status:    order.Status,
		CreatedAt: timestamppb.New(order.CreatedAt),
	}, nil
}

func (s *server) ListOrders(ctx context.Context, _ *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	ordersList, err := s.orderUC.GetListOfOrders(ctx, model.OrderFilter{})
	if err != nil {
		return nil, err
	}

	var protoOrders []*orderpb.GetOrderResponse
	for _, order := range ordersList {
		var protoProducts []*orderpb.OrderedProduct
		for _, p := range order.Products {
			protoProducts = append(protoProducts, &orderpb.OrderedProduct{
				ProductId: p.ProductId.Hex(),
				Quantity:  int32(p.Quantity),
			})
		}
		protoOrders = append(protoOrders, &orderpb.GetOrderResponse{
			Id:        order.ID.Hex(),
			UserId:    order.UserID.Hex(),
			Products:  protoProducts,
			Status:    order.Status,
			CreatedAt: timestamppb.New(order.CreatedAt),
		})
	}

	return &orderpb.ListOrdersResponse{Orders: protoOrders}, nil
}
