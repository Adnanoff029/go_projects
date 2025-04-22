package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the env file")
	}
	MongoDb := os.Getenv("MONGO_URL")
	if MongoDb == "" {
		log.Fatal("MONGO_URL environment variable is not set")
	}
	client, err := mongo.Connect(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}
	_, cancel := context.WithTimeout(context.Background(), 10*(time.Second))
	defer cancel()
	fmt.Println("Connected to MongoDB")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Cluster0").Collection(collectionName)
	return collection
}
