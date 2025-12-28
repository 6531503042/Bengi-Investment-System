package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InstrumentRepository struct {
	collection *mongo.Collection
}

func NewInstrumentRepository() *InstrumentRepository {
	return &InstrumentRepository{
		collection: database.GetCollection(model.InstrumentCollection),
	}
}

func (r *InstrumentRepository) CreateInstrument(ctx context.Context, instrument *model.Instrument) error {
	instrument.CreatedAt = time.Now()
	instrument.UpdatedAt = time.Now()
	instrument.Status = model.InstrumentStatusActive

	result, err := r.collection.InsertOne(ctx, instrument)
	if err != nil {
		return err
	}

	instrument.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *InstrumentRepository) FindByID(ctx context.Context, id string) (*model.Instrument, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var instrument model.Instrument
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&instrument)
	if err != nil {
		return nil, err
	}

	return &instrument, nil
}

func (r *InstrumentRepository) FindBySymbol(ctx context.Context, symbol string) (*model.Instrument, error) {
	var instrument model.Instrument
	err := r.collection.FindOne(ctx, bson.M{"symbol": symbol}).Decode(&instrument)
	if err != nil {
		return nil, err
	}

	return &instrument, nil
}

func (r *InstrumentRepository) FindAll(ctx context.Context, page, limit int) ([]model.Instrument, int64, error) {
	skil := (page - 1) * limit

	total, err := r.collection.CountDocuments(ctx, bson.M{"status": model.InstrumentStatusActive})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "symbol", Value: 1}}).
		SetSkip(int64(skil)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"status": model.InstrumentStatusActive}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var instruments []model.Instrument
	if err := cursor.All(ctx, &instruments); err != nil {
		return nil, 0, err
	}

	return instruments, total, nil
}

func (r *InstrumentRepository) Search(ctx context.Context, query string, instrumentType string, exchange string, page, limit int) ([]model.Instrument, int64, error) {
	skip := (page - 1) * limit

	filter := bson.M{"status": model.InstrumentStatusActive}

	if query != "" {
		filter["$or"] = []bson.M{
			{"symbol": bson.M{"$regex": query, "$options": "i"}},
			{"name": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	if instrumentType != "" {
		filter["type"] = instrumentType
	}

	if exchange != "" {
		filter["exchange"] = exchange
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "symbol", Value: 1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var instruments []model.Instrument
	if err := cursor.All(ctx, &instruments); err != nil {
		return nil, 0, err
	}

	return instruments, total, nil
}

func (r *InstrumentRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (r *InstrumentRepository) SymbolExists(ctx context.Context, symbol string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"symbol": symbol})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// BulkUpsertInstruments inserts new instruments or updates existing ones (by symbol)
func (r *InstrumentRepository) BulkUpsertInstruments(ctx context.Context, instruments []*model.Instrument) (int64, int64, error) {
	if len(instruments) == 0 {
		return 0, 0, nil
	}

	var inserted, updated int64

	// Use MongoDB's BulkWrite for efficiency
	var operations []mongo.WriteModel

	for _, inst := range instruments {
		filter := bson.M{"symbol": inst.Symbol}
		update := bson.M{
			"$set": bson.M{
				"name":        inst.Name,
				"type":        inst.Type,
				"logoUrl":     inst.LogoURL,
				"description": inst.Description,
				"status":      model.InstrumentStatusActive,
				"updatedAt":   time.Now(),
			},
			"$setOnInsert": bson.M{
				"createdAt": time.Now(),
			},
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, op)
	}

	// Execute in batches of 1000
	batchSize := 1000
	for i := 0; i < len(operations); i += batchSize {
		end := i + batchSize
		if end > len(operations) {
			end = len(operations)
		}

		batch := operations[i:end]
		opts := options.BulkWrite().SetOrdered(false)

		result, err := r.collection.BulkWrite(ctx, batch, opts)
		if err != nil {
			return inserted, updated, err
		}

		inserted += result.UpsertedCount
		updated += result.ModifiedCount
	}

	return inserted, updated, nil
}
