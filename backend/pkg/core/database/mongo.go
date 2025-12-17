package database

import (
	"context"
	"log"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnextMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("failed to connect to mongo", err)
	}

	//Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("failed to ping mongo", err)
	}

	DB = client.Database(config.AppConfig.MongoDatabase)
	log.Println("✅ Connected to MongoDB:", config.AppConfig.MongoDatabase)
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
