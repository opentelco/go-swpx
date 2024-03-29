package core

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelco.io/go-swpx/config"
	"go.opentelco.io/go-swpx/proto/go/corepb"
)

const (
	defaultCacheConnectionTimeout = 30 * time.Second
)

type PollResponseCache interface {
	Pop(ctx context.Context, hostname, port, accessId string, rt corepb.PollRequest_Type) (*CachedResponse, error)
	Upsert(ctx context.Context, hostname, port, accessId string, rt corepb.PollRequest_Type, response *corepb.PollResponse) error
	Clear(ctx context.Context, hostname, port, accessId string, rt corepb.PollRequest_Type) error
}

type CachedResponse struct {
	Hostname    string               `bson:"hostname"`
	Port        string               `bson:"port"`
	AccessId    string               `bson:"access_id"`
	RequestType string               `bson:"request_type"`
	Response    *corepb.PollResponse `bson:"response,omitempty"`
	Timestamp   time.Time            `bson:"timestamp" json:"timestamp"`
}

func newResponseCache(ctx context.Context, client *mongo.Client, conf *config.MongoCache, logger hclog.Logger) (PollResponseCache, error) {
	if conf == nil {
		return nil, errors.New("cannot enable response cache: no mongo config")
	}

	logger.Info("enabling response cache", "db", conf.Database, "collection", conf.Collection)

	col := client.Database(conf.Database).Collection(conf.Collection)
	ctx, cancel := context.WithTimeout(ctx, defaultCacheConnectionTimeout)
	defer cancel()
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "hostname", Value: -1}, // prop1 asc
			},
			Options: options.Index().SetUnique(false).SetName("hostname").SetSparse(true),
		},
	})
	if err != nil {
		logger.Error("could not create index", "reason", err)
	}

	c := &respCacheImpl{
		client: client,
		col:    col,
		logger: logger,
	}

	return c, nil
}

type respCacheImpl struct {
	client *mongo.Client
	col    *mongo.Collection
	logger hclog.Logger
}

func (c *respCacheImpl) Pop(ctx context.Context, hostname, port, accessId string, rt corepb.PollRequest_Type) (*CachedResponse, error) {
	var filter bson.M
	if accessId != "" {
		c.logger.Debug("access id is set, filter by that", "access_id", accessId, "rt", rt.String())
		filter = bson.M{
			"request_type": rt.String(),
			"access_id":    accessId,
		}
	} else {
		c.logger.Debug("access id is not set, filter by that", "access_id", accessId, "rt", rt.String())
		filter = bson.M{
			"hostname":     hostname,
			"port":         port,
			"request_type": rt.String(),
		}
	}
	opts := &options.FindOneOptions{
		Sort: bson.M{"_id": -1},
	}
	res := c.col.FindOne(ctx, &filter, opts)
	obj := &CachedResponse{}

	if err := res.Decode(obj); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Warn("err, not found")
			return nil, nil
		}

		return nil, err
	}

	return obj, nil
}

func (c *respCacheImpl) Upsert(ctx context.Context, hostname, port, accessId string, rt corepb.PollRequest_Type, response *corepb.PollResponse) error {

	c.logger.Debug("insert into cache", "hostname", hostname, "port", port, "type", rt.String(), "access_id", accessId)
	_, err := c.col.InsertOne(ctx, &CachedResponse{
		Hostname:    hostname,
		Port:        port,
		AccessId:    accessId,
		RequestType: rt.String(),
		Response:    response,
		Timestamp:   time.Now(),
	})

	if err != nil {
		c.logger.Error("error saving response in cache", "error", err)
		return err
	}

	return nil
}

func (c *respCacheImpl) Clear(ctx context.Context, hostname, port, accessId string, rt corepb.PollRequest_Type) error {
	_, err := c.col.DeleteMany(ctx, bson.M{"hostname": hostname, "port": port, "request_type": rt.String(), "access_id": accessId})
	return err
}
