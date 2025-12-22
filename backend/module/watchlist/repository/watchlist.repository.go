package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WatchlistRepository struct {
	collection *mongo.Collection
}

func NewWatchlistRepository() *WatchlistRepository {
	return &WatchlistRepository{
		collection: database.GetCollection(model.WatchlistCollection),
	}
}

func (r *WatchlistRepository) Create(ctx context.Context, watchlist *model.Watchlist) error {
	watchlist.CreatedAt = time.Now()
	watchlist.UpdatedAt = time.Now()

	if watchlist.Symbols == nil {
		watchlist.Symbols = []string{}
	}

	result, err := r.collection.InsertOne(ctx, watchlist)
	if err != nil {
		return err
	}

	watchlist.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *WatchlistRepository) FindByID(ctx context.Context, id string) (*model.Watchlist, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var watchlist model.Watchlist
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&watchlist)
	if err != nil {
		return nil, err
	}
	return &watchlist, nil
}

func (r *WatchlistRepository) FindByUserID(ctx context.Context, userID string) ([]model.Watchlist, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var watchlists []model.Watchlist
	if err := cursor.All(ctx, &watchlists); err != nil {
		return nil, err
	}
	return watchlists, nil
}

func (r *WatchlistRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *WatchlistRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *WatchlistRepository) AddSymbol(ctx context.Context, id primitive.ObjectID, symbol string) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$addToSet": bson.M{"symbols": symbol},
		"$set":      bson.M{"updatedAt": time.Now()},
	})
	return err
}

func (r *WatchlistRepository) RemoveSymbol(ctx context.Context, id primitive.ObjectID, symbol string) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$pull": bson.M{"symbols": symbol},
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	return err
}

func (r *WatchlistRepository) ClearDefaultWatchlists(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": bson.M{"isDefault": false, "updatedAt": time.Now()}},
	)
	return err
}

func (r *WatchlistRepository) CountByUserID(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"userId": userID})
}
