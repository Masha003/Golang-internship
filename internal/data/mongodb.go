package data

import (
	"context"
	"log"
	"time"

	"github.com/Masha003/Golang-internship/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg config.Config) (*mongo.Database, error) {
	log.Print("Connecting mongoDB")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoUrl))
	if err != nil {
		return nil, err
	}

	return client.Database("mongo"), nil
}
