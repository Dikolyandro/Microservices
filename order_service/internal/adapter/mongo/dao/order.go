package dao

import (
	"context"
	"errors"
	"fmt"
	"log"
	"github.com/recktt77/Microservices-First-/order_service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderDAO struct {
	col *mongo.Collection
}

func NewOrderDAO(db *mongo.Database) *OrderDAO{
	return &OrderDAO{
		col: db.Collection("orders"),
	}
}

func (d *OrderDAO) CreateOrder(ctx context.Context, order model.Order) (primitive.ObjectID, error){
	order.ID = primitive.NewObjectID()
	order.UserID = primitive.NewObjectID()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.IsDeleted = false

	res, err := d.col.InsertOne(ctx, order)
	if err != nil{
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok{
		return primitive.NilObjectID, errors.New("failed to cast iserted id to objectId")
	}

	log.Printf("created order with id: %s", oid.Hex())
	return oid, nil
}

func (d *OrderDAO) GetOrderByID(ctx context.Context, id primitive.ObjectID) (model.Order, error){
	log.Printf("Attempting to get product with ID: %s", id.Hex())
	var order model.Order
	filter := bson.M{
		"_id": id,
		"isdeleted": false,
	}

	log.Printf("MongoDB query filter: %+v", filter)

	err := d.col.FindOne(ctx, filter).Decode(&order)
	if err != nil{
		if errors.Is(err, mongo.ErrNoDocuments){
			log.Printf("Order not found with Id: %s", id.Hex())
			return model.Order{}, model.ErrOrderNotFound
		}
		log.Printf("Error retrieving order: %v", err)
		return model.Order{}, err
	}

	log.Printf("Successfully retrieved order with ID: %s", id.Hex())
	return order, nil
}

func (d *OrderDAO) UpdateOrder(ctx context.Context, id primitive.ObjectID, update model.OrderUpdate) error{
	updateData := bson.M{}

	if update.Products != nil{
		updateData["products"] = *update.Products
	}
	if update.Status != nil{
		updateData["status"] = *update.Status
	}
	if update.IsDeleted != nil{
		updateData["isdeleted"] = *update.IsDeleted
	}
	updateData["updatedat"] = time.Now()

	_, err:=d.col.UpdateOne(ctx, bson.M{
		"_id": id,
		"isdeleted": false,
	}, bson.M{
		"$set":updateData,
	})

	return err
}

func (d *OrderDAO) GetListOfOrders (ctx context.Context, filter model.OrderFilter) ([]model.Order, error){
	query := bson.M{
		"isdeleted": false,
	}

	if filter.UserID != nil {
		query["userid"] = *filter.UserID
	}
	if filter.Products != nil {
		query["products"] = *filter.Products
	}
	if filter.Status != nil{
		query["status"] = *filter.Status
	}

	cursor, err := d.col.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	fmt.Print(cursor)
	defer cursor.Close(ctx)

	var orders []model.Order
	for cursor.Next(ctx){
		var order model.Order
		if err := cursor.Decode(&order); err ==nil{
			orders = append(orders, order)
		}
	}

	return orders, nil
}