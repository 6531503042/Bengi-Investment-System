package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/account/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountRepository struct {
	accountCollection     *mongo.Collection
	transactionCollection *mongo.Collection
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accountCollection:     database.GetCollection(model.AccountCollection),
		transactionCollection: database.GetCollection(model.TransactionCollection),
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *model.Account) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	account.Status = model.AccountStatusActive

	result, err := r.accountCollection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	account.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
func (r *AccountRepository) FindByID(ctx context.Context, id string) (*model.Account, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var account model.Account
	err = r.accountCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByUserID(ctx context.Context, userID string) ([]model.Account, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.accountCollection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var accounts []model.Account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID primitive.ObjectID, newBalance float64) error {
	_, err := r.accountCollection.UpdateByID(ctx, accountID, bson.M{
		"$set": bson.M{
			"balance":   newBalance,
			"updatedAt": time.Now(),
		},
	})
	return err
}
func (r *AccountRepository) UpdateStatus(ctx context.Context, accountID primitive.ObjectID, status model.AccountStatus) error {
	_, err := r.accountCollection.UpdateByID(ctx, accountID, bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	})
	return err
}

func (r *AccountRepository) ExistsByUserAndCurrency(ctx context.Context, userID primitive.ObjectID, currency string) (bool, error) {
	count, err := r.accountCollection.CountDocuments(ctx, bson.M{
		"userId":   userID,
		"currency": currency,
	})
	return count > 0, err
}

func (r *AccountRepository) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	tx.CreatedAt = time.Now()
	result, err := r.transactionCollection.InsertOne(ctx, tx)
	if err != nil {
		return err
	}
	tx.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
func (r *AccountRepository) GetTransactionsByAccountID(ctx context.Context, accountID string, limit, offset int) ([]model.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(accountID)
	if err != nil {
		return nil, err
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.transactionCollection.Find(ctx, bson.M{"accountId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var transactions []model.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

// UpdateField updates a single field on an account
func (r *AccountRepository) UpdateField(ctx context.Context, accountID primitive.ObjectID, field string, value interface{}) error {
	_, err := r.accountCollection.UpdateByID(ctx, accountID, bson.M{
		"$set": bson.M{
			field:       value,
			"updatedAt": time.Now(),
		},
	})
	return err
}

// FindDemoByUserID finds demo account for a user
func (r *AccountRepository) FindDemoByUserID(ctx context.Context, userID string) (*model.Account, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var account model.Account
	err = r.accountCollection.FindOne(ctx, bson.M{
		"userId": objectID,
		"type":   model.AccountTypeDemo,
	}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
