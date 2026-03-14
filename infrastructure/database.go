package infrastructure

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	uri := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	if uri == "" || dbName == "" {
		log.Fatal("❌FATAL !!! : DB_URI or DB_NAME is not set in the .env file")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping failed. Is it running?", err)
	}

	log.Println("✅ Connected to MongoDB 'prepbase_db'")
	database := client.Database(dbName)

	// Create Unique Compound Index for Votes
	// This physically prevents a user from voting on the same question twice
	voteIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "question_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, _ = database.Collection("votes").Indexes().CreateOne(context.TODO(), voteIndex)
	
	return database
}
