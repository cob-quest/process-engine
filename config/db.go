package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	once sync.Once
)

// GetClient returns the singleton instance of the MongoDB client.
func GetClient() *mongo.Client {
	once.Do(func() {
		if clientInstance == nil {
			clientInstance = createDBInstance()
		}
	})
	return clientInstance
}

func createDBInstance() (client *mongo.Client) {

	user := os.Getenv("MONGODB_USERNAME")
	pass := os.Getenv("MONGODB_PASSWORD")
	host := os.Getenv("MONGODB_HOSTNAME")
	port := "27017"

	conn := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)

	log.Printf("Attempting connection with: %s\n", conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("Success!")

	log.Println("Pinging server ...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}
	log.Println("Success!")

	// initialise indexes
	log.Println("Initialising indexes ...")
	InitIndexes(client)
	log.Println("Success!")
	return client
}

func InitIndexes(client *mongo.Client) {

	// cob_process_engine_1 index
	engineCollection := OpenCollection(client, "process_engine")

	engineIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: 1},
			{Key: "corId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	engineIndexCreated, err := engineCollection.Indexes().CreateOne(context.Background(), engineIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created Transaction Index %s\n", engineIndexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("cob").Collection(collectionName)

	return collection
}
