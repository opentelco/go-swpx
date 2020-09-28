/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package core

import (
	"context"
	"fmt"
	"git.liero.se/opentelco/go-swpx/shared"
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
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
	logger.Info(fmt.Sprintf("Attempting to connect to MongoDB: %s", conf.Server))

	// Register custom codecs for protobuf Timestamp and wrapper types
	reg := codecs.Register(bson.NewRegistryBuilder()).Build()

	var err error
	var mongoClient *mongo.Client
	if mongoClient, err = mongo.NewClient(options.Client().ApplyURI(conf.Server), &options.ClientOptions{Registry: reg}); err != nil {
		logger.Error(fmt.Sprintf("error initializing Mongo client: %s", err.Error()))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeoutSeconds)*time.Second)
	defer cancel()
	if err = mongoClient.Connect(ctx); err != nil {
		logger.Error(fmt.Sprintf("error connecting Mongo client: %s", err.Error()))
		return nil, err
	}

	// Check the connection
	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		logger.Error(fmt.Sprintf("Can't ping mongo server: %s", err.Error()))
		return nil, fmt.Errorf("can't reach mongo server")
	}

	logger.Info("Successfully connected to MongoDB")
	return mongoClient, nil
}
