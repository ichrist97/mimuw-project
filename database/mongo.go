package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client instance
var DB *mongo.Client = ConnectDB()
var Ctx context.Context

func readEnv() (string, string) {
	err := godotenv.Load()
	if err == nil {
		fmt.Println("Loaded .env file")
	}

	mongo_host := os.Getenv("MONGO_HOST")
	if len(mongo_host) == 0 {
		mongo_host = "localhost"
	}

	mongo_port := os.Getenv("MONGO_PORT")
	if len(mongo_port) == 0 {
		mongo_port = "27017"
	}

	return mongo_host, mongo_port
}

func ConnectDB() *mongo.Client {
	// load env variables
	mongo_host, mongo_port := readEnv()
	uri := fmt.Sprintf("mongodb://%s:%s", mongo_host, mongo_port)

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	//serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri) //.SetServerAPIOptions(serverAPI)
	Ctx = context.Background()
	// Create a new client and connect to the server
	client, err := mongo.Connect(Ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	/*
		defer func() {
			if err = client.Disconnect(Ctx); err != nil {
				panic(err)
			}
		}()
	*/

	//ping the database
	err = client.Ping(Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}
