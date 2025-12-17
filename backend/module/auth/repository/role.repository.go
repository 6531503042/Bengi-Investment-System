package repository

import (
	"context"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/auth/model"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleRepository struct {
	collection *mongo.Collection
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		collection: database.GetCollection(model.RoleCollection),
	}
}

// Create inserts a new role
func (r *RoleRepository) Create(ctx context.Context, role *model.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return err
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a role by ID
func (r *RoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var role model.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName finds a role by name
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Exists checks if a role exists by name
func (r *RoleRepository) Exists(ctx context.Context, name string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": name})
	return count > 0, err
}

// FindAll returns all roles
func (r *RoleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []model.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}
