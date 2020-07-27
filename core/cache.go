package core

import (
	"context"
	"errors"
	"fmt"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	CACHE_DATABASE   = "test"
	CACHE_COLLECTION = "cache"
)

// CachedInterface is the data object stored in mongo for a cached interface
type CachedInterface struct {
	Hostname string `bson:"hostname"`
	Port     string `bson:"port"`
	// index from the InterfaceTableMIB
	InterfaceIndex int64 `bson:"if_index"`
	// index from the PhysicalEntityMIB
	PhysicalEntityIndex int64 `bson:"physical_entity_index"`
}

type InterfaceCacher interface {
	Pop(hostname, iface string) (*CachedInterface, error)
	Set(ne *proto.NetworkElement, interfaces *proto.NetworkElementInterfaces, phys *proto.NetworkElementInterfaces) error
}

func NewCache(client *mongo.Client, logger hclog.Logger) (*cache, error) {
	col := client.Database(CACHE_DATABASE).Collection(CACHE_COLLECTION)

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

type cache struct {
	client *mongo.Client
	col    *mongo.Collection
	logger hclog.Logger
}

func (c *cache) Pop(hostname, iface string) (*CachedInterface, error) {
	res := c.col.FindOne(context.Background(), bson.M{"hostname": hostname, "port": iface})
	obj := &CachedInterface{}
	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *cache) Set(ne *proto.NetworkElement, interfaces *proto.NetworkElementInterfaces, phys *proto.NetworkElementInterfaces) error {
	for k, v := range interfaces.Interfaces {
		if physInterface, ok := phys.Interfaces[k]; ok {
			_, err := c.col.InsertOne(context.Background(), bson.M{
				"hostname":              ne.Hostname,
				"port":                  v.Description,
				"if_index":              v.Index,
				"physical_entity_index": physInterface.Index,
			})

			if err != nil {
				logger.Error("Error saving info in cache: ", err.Error())
				return err
			}

		}
	}
	return nil

}

func initMongoDB(uri string) (*mongo.Client, error) {
	// TODO timeout from config
	logger.Info("Attempting to connect to MongoDB...")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var err error
	var mongoClient *mongo.Client
	if mongoClient, err = mongo.NewClient(options.Client().ApplyURI(uri)); err != nil {
		logger.Error("error initializing Mongo client:", err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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