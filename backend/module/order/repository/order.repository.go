package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/order/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		collection: database.GetCollection(model.OrderCollection),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = model.OrderStatusPending

	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order model.Order
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByUserID(ctx context.Context, userID string, filter *model.OrderFilter, page, limit int) ([]model.Order, int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	query := bson.M{"userId": userObjectID}

	if filter != nil {
		if filter.Status != "" {
			query["status"] = filter.Status
		}
		if filter.Side != "" {
			query["side"] = filter.Side
		}
		if filter.Symbol != "" {
			query["symbol"] = filter.Symbol
		}
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var orders []model.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.OrderStatus) error {
	update := bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}

	if status == model.OrderStatusCancelled {
		now := time.Now()
		update["cancelledAt"] = now
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *OrderRepository) UpdateFill(ctx context.Context, id primitive.ObjectID, filledQty, avgPrice float64, status model.OrderStatus) error {
	update := bson.M{
		"filledQty":    filledQty,
		"avgFillPrice": avgPrice,
		"status":       status,
		"updatedAt":    time.Now(),
	}

	if status == model.OrderStatusFilled {
		now := time.Now()
		update["filledAt"] = now
	}

	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *OrderRepository) FindOpenOrders(ctx context.Context, portfolioID primitive.ObjectID, symbol string) ([]model.Order, error) {
	query := bson.M{
		"portfolioId": portfolioID,
		"symbol":      symbol,
		"status": bson.M{"$in": []model.OrderStatus{
			model.OrderStatusPending,
			model.OrderStatusOpen,
			model.OrderStatusPartiallyFilled,
		}},
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []model.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
