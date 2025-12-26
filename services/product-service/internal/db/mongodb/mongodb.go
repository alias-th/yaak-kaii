package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoClient wraps the MongoDB client
type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect establishes a connection to MongoDB
func Connect(ctx context.Context, uri string, dbName string) (*MongoClient, error) {
	// Create a new client
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify the connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Connected to MongoDB: %s", dbName)

	return &MongoClient{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Disconnect closes the MongoDB connection
func (mc *MongoClient) Disconnect(ctx context.Context) error {
	if err := mc.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	log.Println("Disconnected from MongoDB")
	return nil
}

// GetDatabase returns the MongoDB database instance
func (mc *MongoClient) GetDatabase() *mongo.Database {
	return mc.Database
}
