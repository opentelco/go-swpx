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
	CACHE_DATABASE             = "test"
	CACHE_COLLECTION           = "cache"
	PHYS_PORT_CACHE_COLLECTION = "physical_port_information"
)

// CachedInterface is the data object stored in mongo for a cached interface
type CachedInterface struct {
	Hostname    string `bson:"hostname"`
	Interface   string `bson:"interface"`
	Description string `bson:"description"`
	Alias       string `bson:"alias"`
	// index from the InterfaceTableMIB
	InterfaceIndex int64 `bson:"if_index"`
	// index from the PhysicalEntityMIB
	PhysicalEntityIndex    string                           `bson:"physical_entity_index"`
	TransceiverInformation *proto.VRPTransceiverInformation `bson:"transceiver_information"`
}

type CachedPhysicalPortInformation struct {
	Provider                string                           `bson:"provider"`
	Driver                  string                           `bson:"driver"`
	PhysicalPortInformation []*proto.PhysicalPortInformation `bson:"physical_port_information"`
}

type InterfaceCacher interface {
	Pop(hostname, iface string) (*CachedInterface, error)
	Set(ne *proto.NetworkElement, nei *proto.NetworkElementInterface, phys *proto.PhysicalPortinformationResponse,
		transceiver *proto.VRPTransceiverInformation) (*CachedInterface, error)
}

type PhysicalPortCacher interface {
	PopPhysical(provider, driver string) (*CachedPhysicalPortInformation, error)
	SetPhysical(provider, driver string, phys *proto.PhysicalPortinformationResponse) (*CachedPhysicalPortInformation, error)
}

func NewCache(client *mongo.Client, logger hclog.Logger) (*cache, error) {
	col := client.Database(CACHE_DATABASE).Collection(CACHE_COLLECTION)

	model := mongo.IndexModel{
		Keys:    bson.M{"hostname": -1, "interface": -1},
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

func NewPhysicalPortCache(client *mongo.Client, logger hclog.Logger) (PhysicalPortCacher, error) {
	col := client.Database(CACHE_DATABASE).Collection(PHYS_PORT_CACHE_COLLECTION)

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
	res := c.col.FindOne(context.Background(), bson.M{"hostname": hostname, "interface": iface})
	obj := &CachedInterface{}
	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *cache) Set(ne *proto.NetworkElement, nei *proto.NetworkElementInterface,
	phys *proto.PhysicalPortinformationResponse, tc *proto.VRPTransceiverInformation) (*CachedInterface, error) {
	obj := CachedInterface{
		Hostname:               ne.Hostname,
		Interface:              ne.Interface,
		Description:            ne.Hostname,
		Alias:                  nei.Alias,
		InterfaceIndex:         nei.Index,
		TransceiverInformation: tc,
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

func (c *cache) PopPhysical(provider, driver string) (*CachedPhysicalPortInformation, error) {
	res := c.col.FindOne(context.Background(), bson.M{"provider": provider, "driver": driver})
	obj := &CachedPhysicalPortInformation{}
	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *cache) SetPhysical(provider, driver string, phys *proto.PhysicalPortinformationResponse) (*CachedPhysicalPortInformation, error) {
	obj := CachedPhysicalPortInformation{
		Provider:                provider,
		Driver:                  driver,
		PhysicalPortInformation: phys.Data,
	}

	_, err := c.col.InsertOne(
		context.Background(),
		&obj,
	)
	if err != nil {
		logger.Error("Error saving physical port info in cache: ", err)
		return nil, err
	}
	return &obj, nil

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
