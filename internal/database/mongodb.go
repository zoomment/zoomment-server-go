package database

import (
	"context"
	"log"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect initializes mgm and VERIFIES the MongoDB connection
func Connect(uri string) error {
	// Configure mgm with MongoDB URI
	err := mgm.SetDefaultConfig(
		nil,        // Use default config
		"zoomment", // Database name
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		return err
	}

	// Actually verify the connection works by pinging MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the client and ping it
	_, client, _, err := mgm.DefaultConfigs()
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("❌ Failed to connect to MongoDB: %v", err)
		return err
	}

	log.Println("✅ Connected to MongoDB")
	return nil
}

// GetCollection returns a collection for a model
func GetCollection(model mgm.Model) *mgm.Collection {
	return mgm.Coll(model)
}
