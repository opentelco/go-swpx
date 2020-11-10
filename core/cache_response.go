package core

import (
	"context"
	"errors"
	
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/shared"
)

type ResponseCache interface {
	Pop(ctx context.Context, hostname, port string, rt pb_core.Request_Type) (*CachedResponse, error)
	Upsert(ctx context.Context, hostname, port string, rt pb_core.Request_Type, response *pb_core.Response) error
	Clear(ctx context.Context, hostname, port string, rt pb_core.Request_Type) error
}

type CachedResponse struct {
	Hostname    string                                 `bson:"hostname"`
	Port        string                                 `bson:"port"`
	RequestType pb_core.Request_Type                   `bson:"request_type"`
	Response    *pb_core.Response                      `bson:"response,omitempty"`
	Timestamp   *timestamp.Timestamp                   `bson:"timestamp" json:"timestamp"`
}


func newResponseCache(client *mongo.Client, logger hclog.Logger, conf shared.ConfigMongo) (ResponseCache, error) {
	col := client.Database(conf.Database).Collection(collectionResponseCache)
	// Create the model
	model := mongo.IndexModel{
		Keys:    bson.M{"hostname": -1, "port": -1},
		Options: options.Index().SetUnique(true),
	}
	c := &respCacheImpl{
		client: client,
		col:    col,
		logger: logger,
	}
	err := createIndex(col, model, logger)
	if err != nil {
		return nil, err
	}
	
	return c, nil
}

type respCacheImpl struct {
	client *mongo.Client
	col    *mongo.Collection
	logger hclog.Logger
}



func (c *respCacheImpl) Pop(ctx context.Context, hostname, port string, rt pb_core.Request_Type) (*CachedResponse, error) {
	res := c.col.FindOne(context.Background(), bson.M{"hostname": hostname, "port": port, "request_type": rt})
	obj := &CachedResponse{}

	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return obj, nil
}

func (c *respCacheImpl) Upsert(ctx context.Context, hostname, port string, rt pb_core.Request_Type, response *pb_core.Response) error {
	_, err := c.col.InsertOne(context.Background(), &CachedResponse{
		hostname,
		port,
		rt,
		response,
		ptypes.TimestampNow(),
	})

	if err != nil {
		c.logger.Error("error saving response in cache","error",  err)
		return err
	}

	return nil
}

func (c *respCacheImpl) Clear(ctx context.Context, hostname, port string, rt pb_core.Request_Type) error {
	_, err := c.col.DeleteMany(context.Background(), bson.M{"hostname": hostname, "port": port, "request_type": rt})
	return err
}
