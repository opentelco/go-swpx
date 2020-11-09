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
	
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"git.liero.se/opentelco/go-swpx/shared"
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
c := &cache{
	client: client,
	col:    col,
	logger: logger,
}
	err := c.createIndex(col, model, logger)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// create index for cache
func (c *cache) createIndex(col *mongo.Collection, model mongo.IndexModel, logger hclog.Logger) error {
	cursor, err := col.Indexes().List(context.TODO(), options.ListIndexes())
	if err != nil {
		return err
	}
	var indexes []bson.M
	if err = cursor.All(context.TODO(), &indexes); err != nil {
		return err
	}

	if len(indexes) == 0 {
		if _, err := col.Indexes().CreateOne(context.Background(), model); err != nil {
			c.logger.Warn("can't create index","error", err)
			return err
		}
	}
	return nil
}

