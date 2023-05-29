package core

import (
	"context"
	"errors"

	"git.liero.se/opentelco/go-swpx/config"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PollInterfaceCache interface {
	Pop(ctx context.Context, hostname, port string) (*CachedInterface, error)
	Upsert(ctx context.Context, hostname string, logicalPortIndex *proto.PortIndex, physicalPortIndex *proto.PortIndex) error
}

// CachedInterface is the data object stored in mongo for a cached interface
type CachedInterface struct {
	Hostname string `bson:"hostname"`
	Port     string `bson:"port"`
	// index from the InterfaceTableMIB
	InterfaceIndex int64 `bson:"if_index"`
	// index from the PhysicalEntityMIB
	PhysicalEntityIndex int64 `bson:"physical_entity_index"`
}

func newInterfaceCache(ctx context.Context, client *mongo.Client, conf *config.MongoCache, logger hclog.Logger) (PollInterfaceCache, error) {
	if conf == nil {
		return nil, errors.New("cannot enable interface cache: no mongo config")
	}

	logger.Info("enabling response cache", "db", conf.Database, "collection", conf.Collection)
	col := client.Database(conf.Database).Collection(conf.Collection)
	ctx, cancel := context.WithTimeout(ctx, defaultCacheConnectionTimeout)
	defer cancel()

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "hostname", Value: -1}, // prop1 asc
				{Key: "port", Value: -1},     // prop2 asc
			},
			Options: options.Index().SetUnique(true).SetName("hostname@port").SetSparse(true),
		},
	})
	if err != nil {
		logger.Error("could not create index", "reason", err)
	}

	c := &ifCacheImpl{
		client: client,
		col:    col,
		logger: logger,
	}

	return c, nil
}

type ifCacheImpl struct {
	client *mongo.Client
	col    *mongo.Collection
	logger hclog.Logger
}

func (c *ifCacheImpl) Pop(ctx context.Context, hostname, iface string) (*CachedInterface, error) {

	opts := &options.FindOneOptions{
		Sort: bson.M{"_id": -1},
	}
	res := c.col.FindOne(ctx, bson.M{"hostname": hostname, "port": iface}, opts)
	obj := &CachedInterface{}
	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *ifCacheImpl) Upsert(ctx context.Context, hostname string, logicalPortIndex *proto.PortIndex, physicalPortIndex *proto.PortIndex) error {
	for k, v := range logicalPortIndex.Ports {

		data := bson.M{
			"hostname": hostname,
			"port":     v.Description,
			"if_index": v.Index,
		}

		// All vendors (CTC Union) does not have implemented/enabled the phys entity MIB
		// This means that we cannot know for sure if the PHYS is set.
		if physicalPortIndex != nil { // check if phys is set
			if physInterface, ok := physicalPortIndex.Ports[k]; ok {
				data["physical_entity_index"] = physInterface.Index
			}
		}

		opts := options.Update().SetUpsert(true)
		filter := bson.M{
			"hostname": hostname,
			"port":     v.Description,
			"if_index": v.Index,
		}
		_, err := c.col.UpdateOne(ctx, filter, bson.M{"$set": data}, opts)

		if err != nil {
			c.logger.Error("error saving interface in cache", "error", err, "hostname", hostname, "port", v.Description)
			return err
		}
	}
	return nil

}
