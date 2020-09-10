package core

import (
	"context"
	"fmt"
	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type cache struct {
	client *mongo.Client
	col    *mongo.Collection
	logger hclog.Logger
}

func NewCache(client *mongo.Client, logger hclog.Logger, conf shared.ConfigMongo) (*cache, error) {
	col := client.Database(conf.Database).Collection(conf.Collection)

	model := mongo.IndexModel{
		Keys:    bson.M{"hostname": -1, "port": -1},
		Options: options.Index().SetUnique(true),
	}

	if _, err := col.Indexes().CreateOne(context.Background(), model); err != nil {
		logger.Warn("can't create index:", err.Error())
	}

	return &cache{
		client: client,
		col:    col,
		logger: logger,
	}, nil
}

func initMongoDB(conf shared.ConfigMongo) (*mongo.Client, error) {
	logger.Info("Attempting to connect to MongoDB...")

	var err error
	var mongoClient *mongo.Client
	if mongoClient, err = mongo.NewClient(options.Client().ApplyURI(conf.Server)); err != nil {
		logger.Error("error initializing Mongo client:", err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeoutSeconds)*time.Second)
	defer cancel()
	if err = mongoClient.Connect(ctx); err != nil {
		logger.Error("error connecting Mongo client:", err.Error())
		return nil, err
	}

	// Check the connection
	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		logger.Error("Can't ping mongo server:", err.Error())
		return nil, fmt.Errorf("can't reach mongo server")
	}

	logger.Info("Successfully connected to MongoDB")
	return mongoClient, nil
}
