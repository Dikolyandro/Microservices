package grpc

import (
	"context"
	"log"
	_ "statistics_service/internal/core/domain"
	"statistics_service/internal/usecase"
	pb "statistics_service/proto"

	"google.golang.org/grpc"
	"net"
	"time"
)

type statisticsServer struct {
	pb.UnimplementedStatisticsServiceServer
	useCase *usecase.StatisticsUseCase
}

func NewServer(useCase *usecase.StatisticsUseCase) {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStatisticsServiceServer(s, &statisticsServer{useCase: useCase})

	log.Println("gRPC Statistics server is running on port 50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *statisticsServer) GetUserOrderStats(ctx context.Context, req *pb.UserStatsRequest) (*pb.UserStatsResponse, error) {
	stats, err := s.useCase.GetUserOrderStats(req.UserId)
	if err != nil {
		return nil, err
	}

	var grpcStats []*pb.OrderStat
	for _, stat := range stats {
		grpcStats = append(grpcStats, &pb.OrderStat{
			OrderId:   stat.OrderID,
			UserId:    stat.UserID,
			CreatedAt: stat.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.UserStatsResponse{Orders: grpcStats}, nil
}
