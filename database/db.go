// database/db.go
package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ClipboardCollection *mongo.Collection
var client *mongo.Client

func Init() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "clipboard"
	}

	ClipboardCollection = client.Database(dbName).Collection("items")
}

func Disconnect() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from database: %v", err)
		}
	}
}
