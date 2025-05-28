package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var LogCollection *mongo.Collection

func InitMongo() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	collectionName := os.Getenv("MONGO_COLLECTION")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connect error: ", err)
	}

	MongoClient = client
	LogCollection = client.Database(dbName).Collection(collectionName)

	fmt.Println("✅ Connected to MongoDB")
}
