package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var db *mongo.Database

// InitDB initializes the MongoDB database
func InitDB(connectionString string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db = client.Database("ticketoffdb")

	log.Println("Database connected successfully")
	return db, nil
}

// GetDB returns the DB instance for use in other packages
func GetDB() *mongo.Database {
	return db
}
