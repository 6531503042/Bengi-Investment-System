package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PortfolioRepository struct {
	portfolioCollection   *mongo.Collection
	positionCollection    *mongo.Collection
	positionLotCollection *mongo.Collection
}

func NewPortfolioRepository() *PortfolioRepository {
	return &PortfolioRepository{
		portfolioCollection:   database.GetCollection(model.PortfolioCollection),
		positionCollection:    database.GetCollection(model.PositionCollection),
		positionLotCollection: database.GetCollection(model.PositionLotCollection),
	}
}

// ==================== Portfolio Methods ====================

func (r *PortfolioRepository) CreatePortfolio(ctx context.Context, portfolio *model.Portfolio) error {
	portfolio.CreatedAt = time.Now()
	portfolio.UpdatedAt = time.Now()

	result, err := r.portfolioCollection.InsertOne(ctx, portfolio)
	if err != nil {
		return err
	}

	portfolio.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PortfolioRepository) FindPortfolioByID(ctx context.Context, id string) (*model.Portfolio, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var portfolio model.Portfolio
	err = r.portfolioCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&portfolio)
	if err != nil {
		return nil, err
	}
	return &portfolio, nil
}

func (r *PortfolioRepository) FindPortfoliosByUserID(ctx context.Context, userID string) ([]model.Portfolio, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.portfolioCollection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var portfolios []model.Portfolio
	if err := cursor.All(ctx, &portfolios); err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (r *PortfolioRepository) UpdatePortfolio(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.portfolioCollection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *PortfolioRepository) DeletePortfolio(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.portfolioCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *PortfolioRepository) ClearDefaultPortfolios(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.portfolioCollection.UpdateMany(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": bson.M{"isDefault": false, "updatedAt": time.Now()}},
	)
	return err
}

// ==================== Position Methods ====================

func (r *PortfolioRepository) CreatePosition(ctx context.Context, position *model.Position) error {
	position.CreatedAt = time.Now()
	position.UpdatedAt = time.Now()

	result, err := r.positionCollection.InsertOne(ctx, position)
	if err != nil {
		return err
	}

	position.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PortfolioRepository) FindPositionByID(ctx context.Context, id string) (*model.Position, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var position model.Position
	err = r.positionCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&position)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *PortfolioRepository) FindPositionsByPortfolioID(ctx context.Context, portfolioID string) ([]model.Position, error) {
	objectID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.D{{Key: "symbol", Value: 1}})
	cursor, err := r.positionCollection.Find(ctx, bson.M{"portfolioId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var positions []model.Position
	if err := cursor.All(ctx, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *PortfolioRepository) FindPositionByPortfolioAndSymbol(ctx context.Context, portfolioID primitive.ObjectID, symbol string) (*model.Position, error) {
	var position model.Position
	err := r.positionCollection.FindOne(ctx, bson.M{
		"portfolioId": portfolioID,
		"symbol":      symbol,
	}).Decode(&position)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *PortfolioRepository) UpdatePosition(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.positionCollection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *PortfolioRepository) DeletePosition(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.positionCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ==================== Position Lot Methods ====================

func (r *PortfolioRepository) CreatePositionLot(ctx context.Context, lot *model.PositionLot) error {
	lot.CreatedAt = time.Now()

	result, err := r.positionLotCollection.InsertOne(ctx, lot)
	if err != nil {
		return err
	}

	lot.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PortfolioRepository) FindLotsByPositionID(ctx context.Context, positionID primitive.ObjectID) ([]model.PositionLot, error) {
	opts := options.Find().SetSort(bson.D{{Key: "purchasedAt", Value: 1}}) // FIFO order
	cursor, err := r.positionLotCollection.Find(ctx, bson.M{
		"positionId":   positionID,
		"remainingQty": bson.M{"$gt": 0},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var lots []model.PositionLot
	if err := cursor.All(ctx, &lots); err != nil {
		return nil, err
	}
	return lots, nil
}

func (r *PortfolioRepository) UpdatePositionLot(ctx context.Context, id primitive.ObjectID, remainingQty float64) error {
	_, err := r.positionLotCollection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{"remainingQty": remainingQty},
	})
	return err
}
