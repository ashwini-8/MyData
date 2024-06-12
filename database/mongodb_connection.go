package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI     = "mongodb://mongodb:27017/vehicles_db"
	DatabaseName = "vehicles_db"
)

var mongoClientInstance *mongo.Client
var mongoClientOnce sync.Once

// GetMongoClient returns a singleton MongoDB client instance
func GetMongoClient() (*mongo.Client, error) {
	var err error
	mongoClientOnce.Do(func() {
		/* this function ensures MongoDB client is initialized only once. Even if multiple
		   goroutines call GetMongoClient concurrently
		   only the first call will initialize the client.
		*/
		clientOptions := options.Client().ApplyURI(mongoURI)

		// if operation not completed within 10 sec it will be aborted
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)

		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		mongoClientInstance = client
	})
	return mongoClientInstance, err
}
