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
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)


const (
	CACHE_DATABASE = "test"
	CACHE_COLLECTION = "cache"
)

// CachedInterface is the data object stored in mongo for a cached interface
type CachedInterface struct {
	Hostname string `bson:"hostname"`
	Interface string `bson:"interface"`
	Description string `bson:"description"`
	Alias       string `bson:"alias"`
	// index from the InterfaceTableMIB
	InterfaceIndex       int64  `bson:"if_index"`
	// index from the PhysicalEntityMIB
	PhysicalEntityIndex int64   `bson:"physical_entity_index"`

}

type InterfaceCache interface {
	Pop(hostname, iface string) (*CachedInterface,error)
	Set(ne *proto.NetworkElement, nei *proto.NetworkElementInterface) (*CachedInterface, error)
}

func NewCache(client *mongo.Client, logger  hclog.Logger) (InterfaceCache, error) {
	col := client.Database(CACHE_DATABASE).Collection(CACHE_COLLECTION)
	return &cache{
		client: client,
		col:    col,
		logger: logger,
	}, nil
}

type cache struct {
	client *mongo.Client
	col *mongo.Collection
	logger hclog.Logger
}

func (c *cache) Pop(hostname, iface string) (*CachedInterface, error) {
	res := c.col.FindOne(context.Background(),bson.M{"hostname": hostname, "interface": iface} )
	obj := &CachedInterface{}
	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
	return nil, err
	}
	return obj, nil
}

func (c *cache) Set(ne *proto.NetworkElement, nei *proto.NetworkElementInterface) (*CachedInterface,error) {
	obj := CachedInterface{
		Hostname: ne.Hostname,
		InterfaceIndex:       nei.Index,
		Interface: ne.Interface,
		Description: ne.Hostname,
		Alias:       nei.Alias,
	}

	_, err := c.col.InsertOne(
		context.Background(),
		&obj,
	)
	if err != nil {
		logger.Error("Error saving info in cache: ", err)
		return nil, err
	}
	return &obj, nil

}


// todo we should to set a unique index on host and interface?
func initMongoDB(uri string) (*mongo.Client, error){
	// TODO timeout from config
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var err error
	var mongoClient *mongo.Client
	//TODO parametrize the URI (eg. read from config file)
	if mongoClient, err = mongo.NewClient(options.Client().ApplyURI(uri)); err != nil {
		logger.Error("Error initializing Mongo client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mongoClient.Connect(ctx); err != nil {
		logger.Error("Error connecting Mongo client:", err)
	}

	// TODO double ping ? reuse context?
	// Check the connection
	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		logger.Error("Can't ping mongo server:", err)
		return nil, fmt.Errorf("cant reach mongo server")
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("cant reach mongo server")

	}
	return mongoClient, nil
}

