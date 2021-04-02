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
var mongoURI = fmt.Sprintf("mongodb+srv://%s:%s@fiber-backend.kooym.mongodb.net/%s?retryWrites=true&w=majority", config.Config("DB_USER"), config.Config("DB_USER_PASSWORD"), dbName)

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
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
