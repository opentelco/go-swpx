package core

import (
	"context"
	"errors"
	"time"

	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionInterfaceCache = "cache_interface"
	collectionResponseCache  = "cache_response"
)

type InterfaceCache interface {
	Pop(ctx context.Context, hostname, port string) (*CachedInterface, error)
	Upsert(ctx context.Context, ne *proto.NetworkElement, interfaces *proto.NetworkElementInterfaces, phys *proto.NetworkElementInterfaces) error
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

func newInterfaceCache(client *mongo.Client, logger hclog.Logger, conf shared.ConfigMongo) (InterfaceCache, error) {
	col := client.Database(conf.Database).Collection(collectionInterfaceCache)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
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

func (c *ifCacheImpl) Upsert(ctx context.Context, ne *proto.NetworkElement, interfaces *proto.NetworkElementInterfaces, phys *proto.NetworkElementInterfaces) error {
	for k, v := range interfaces.Interfaces {
		if physInterface, ok := phys.Interfaces[k]; ok {
			opts := options.Update().SetUpsert(true)
			filter := bson.M{
				"hostname": ne.Hostname,
				"port":     v.Description,
			}
			_, err := c.col.UpdateOne(ctx, filter,
				bson.M{"$set": bson.M{
					"hostname":              ne.Hostname,
					"port":                  v.Description,
					"if_index":              v.Index,
					"physical_entity_index": physInterface.Index,
				}}, opts)

			if err != nil {
				c.logger.Error("error saving interface in cache", "error", err, "hostname", ne.Hostname, "port", v.Description)
				return err
			}

		}
	}
	return nil

}
