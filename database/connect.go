package database

import (
	"context"
	"fmt"
	"time"

	"github.com/D-Bald/fiber-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database settings (insert your own database host, name and user data)
var mongoHost = config.Config("MONGO_HOST")
var dbName = config.Config("DB_NAME")
var dbUser = config.Config("DB_USER")
var dbUserPassword = config.Config("DB_USER_PASSWORD")

func Connect() error {
	var client *mongo.Client
	var err error

	fmt.Println(mongoHost)

	// Distinguish between Atlas and Docker hostet mongo databases
	if mongoHost == "ATLAS" {
		mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@fiber-backend.kooym.mongodb.net/%s?retryWrites=true&w=majority", dbUser, dbUserPassword, dbName)
		client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}
	} else {
		mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", dbUser, dbUserPassword, mongoHost, config.Config("MONGO_PORT"), dbName)
		client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	DB = client.Database(dbName)

	return nil
}
