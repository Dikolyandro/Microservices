package storage

import (
	"context"
	_ "fmt"
	"time"

	"statistics_service/internal/core/domain"
	"statistics_service/internal/core/ports"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	orderColl     *mongo.Collection
	inventoryColl *mongo.Collection
}

func NewMongoRepo(uri, dbName string) (ports.StatisticsRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoRepo{
		orderColl:     db.Collection("orders"),
		inventoryColl: db.Collection("inventory"),
	}, nil
}

func (m *MongoRepo) SaveOrderStat(stat domain.OrderStat) error {
	_, err := m.orderColl.InsertOne(context.Background(), stat)
	return err
}

func (m *MongoRepo) SaveInventoryStat(stat domain.InventoryStat) error {
	_, err := m.inventoryColl.InsertOne(context.Background(), stat)
	return err
}

func (m *MongoRepo) GetUserOrderStats(userID string) ([]domain.OrderStat, error) {
	var results []domain.OrderStat
	cursor, err := m.orderColl.Find(context.Background(), map[string]interface{}{"user_id": userID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
