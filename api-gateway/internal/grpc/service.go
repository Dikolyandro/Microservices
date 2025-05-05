package grpc

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventorypb "github.com/recktt77/proto-definitions/gen/inventory"
	orderpb "github.com/recktt77/proto-definitions/gen/orders"
)

type Clients struct {
	Inventory inventorypb.ProductServiceClient
	Discount  inventorypb.DiscountServiceClient
	Order     orderpb.OrderServiceClient
}

func NewClients(inventoryAddr, orderAddr string) *Clients {
	invConn, err := grpc.NewClient(inventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}

	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}

	return &Clients{
		Inventory: inventorypb.NewProductServiceClient(invConn),
		Discount:  inventorypb.NewDiscountServiceClient(invConn),
		Order:     orderpb.NewOrderServiceClient(orderConn),
	}
}
