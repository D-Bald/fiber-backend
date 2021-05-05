package database

import (
	"context"
	"fmt"
	"time"

	"github.com/D-Bald/fiber-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database settings (insert your own database name and connection URI)
var dbName = config.Config("DB_NAME")
var mongoAtlasURI = fmt.Sprintf("mongodb+srv://%s:%s@fiber-backend.kooym.mongodb.net/%s?retryWrites=true&w=majority", config.Config("DB_USER"), config.Config("DB_USER_PASSWORD"), dbName)
var mongoDockerURI = fmt.Sprintf("mongodb://%s:%s@mongodb:%s", config.Config("DB_USER"), config.Config("DB_USER_PASSWORD"), config.Config("MONGO_PORT"))

func Connect() error {
	var client *mongo.Client
	var err error
	switch config.Config("HOSTED") {
	case "ATLAS":
		client, err = mongo.NewClient(options.Client().ApplyURI(mongoAtlasURI))
		if err != nil {
			return err
		}
	case "DOCKER":
		client, err = mongo.NewClient(options.Client().ApplyURI(mongoDockerURI))
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
