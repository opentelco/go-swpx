package core

import (
	"fmt"
	"time"
	
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"git.liero.se/opentelco/go-swpx/shared"
)

func initMongoDb(conf shared.ConfigMongo, logger hclog.Logger) (*mongo.Client, error) {
	logger.Info("attempting to connect to mongoDb server", "server", conf.Server)
	
	// Register custom codecs for protobuf Timestamp and wrapper types
	reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	
	var err error
	var mongoClient *mongo.Client
	if mongoClient, err = mongo.NewClient(options.Client().ApplyURI(conf.Server), &options.ClientOptions{Registry: reg}); err != nil {
		logger.Error("error initializing mongo client","error", err, "server", conf.Server)
		return nil, err
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeoutSeconds)*time.Second)
	defer cancel()
	if err = mongoClient.Connect(ctx); err != nil {
		logger.Error("error connecting to server","error", err, "server", conf.Server)
		return nil, err
	}
	
	// Check the connection
	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		logger.Error("can't ping mongo server", "error", err, "server", conf.Server)
		return nil, fmt.Errorf("can't reach mongo server")
	}
	
	logger.Info("successfully connected to MongoDB", "server", conf.Server)
	return mongoClient, nil
}
