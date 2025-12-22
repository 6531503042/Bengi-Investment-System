package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/trade/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TradeRepository struct {
	collection *mongo.Collection
}

func NewTradeRepository() *TradeRepository {
	return &TradeRepository{
		collection: database.GetCollection(model.TradeCollection),
	}
}

func (r *TradeRepository) Create(ctx context.Context, trade *model.Trade) error {
	trade.CreatedAt = time.Now()
	if trade.ExecutedAt.IsZero() {
		trade.ExecutedAt = time.Now()
	}

	result, err := r.collection.InsertOne(ctx, trade)
	if err != nil {
		return err
	}

	trade.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *TradeRepository) FindByOrderID(ctx context.Context, orderID string) ([]model.Trade, error) {
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"orderId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var trades []model.Trade
	if err := cursor.All(ctx, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

func (r *TradeRepository) FindByID(ctx context.Context, id string) (*model.Trade, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var trade model.Trade
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&trade)
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *TradeRepository) FindByUserID(ctx context.Context, userID string, filter *model.TradeFilter, page, limit int) ([]model.Trade, int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}
	skip := (page - 1) * limit
	query := bson.M{"userId": userObjectID}
	if filter != nil {
		if filter.Symbol != "" {
			query["symbol"] = filter.Symbol
		}
		if filter.Side != "" {
			query["side"] = filter.Side
		}
	}
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "executedAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var trades []model.Trade
	if err := cursor.All(ctx, &trades); err != nil {
		return nil, 0, err
	}
	return trades, total, nil
}

func (r *TradeRepository) FindByPortfolioID(ctx context.Context, portfolioID string, limit int) ([]model.Trade, error) {
	objectID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return nil, err
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "executedAt", Value: -1}}).
		SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.M{"portfolioId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var trades []model.Trade
	if err := cursor.All(ctx, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

func (r *TradeRepository) GetTradeSummary(ctx context.Context, userID primitive.ObjectID) (*TradeSummary, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"userId": userID}},
		{"$group": bson.M{
			"_id":             nil,
			"totalTrades":     bson.M{"$sum": 1},
			"totalBuyValue":   bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []string{"$side", "BUY"}}, "$total", 0}}},
			"totalSellValue":  bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []string{"$side", "SELL"}}, "$total", 0}}},
			"totalCommission": bson.M{"$sum": "$commission"},
		}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []TradeSummary
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &TradeSummary{}, nil
	}
	result := results[0]
	result.NetProfit = result.TotalSellValue - result.TotalBuyValue - result.TotalCommission
	return &result, nil
}

type TradeSummary struct {
	TotalTrades     int     `bson:"totalTrades"`
	TotalBuyValue   float64 `bson:"totalBuyValue"`
	TotalSellValue  float64 `bson:"totalSellValue"`
	TotalCommission float64 `bson:"totalCommission"`
	NetProfit       float64 `bson:"netProfit"`
}
