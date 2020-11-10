package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Create the model
	model := mongo.IndexModel{
		Keys:    bson.M{"hostname": -1, "port": -1},
		Options: options.Index().SetUnique(true),
	}
	c := &ifCacheImpl{
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

type ifCacheImpl struct {
	client *mongo.Client
	col    *mongo.Collection
	logger hclog.Logger
}


func (c *ifCacheImpl) Pop(ctx context.Context, hostname, iface string) (*CachedInterface, error) {
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

func (c *ifCacheImpl) Upsert(ctx context.Context, ne *proto.NetworkElement, interfaces *proto.NetworkElementInterfaces, phys *proto.NetworkElementInterfaces) error {
	js, _ := json.MarshalIndent(ne, "", "  ")
	fmt.Println(string(js))
	
	js, _ = json.MarshalIndent(interfaces, "", "  ")
	fmt.Println(string(js))
	
	js, _ = json.MarshalIndent(phys, "", "  ")
	fmt.Println(string(js))
	
	for k, v := range interfaces.Interfaces {
		if physInterface, ok := phys.Interfaces[k]; ok {
			_, err := c.col.InsertOne(context.Background(), bson.M{
				"hostname":              ne.Hostname,
				"port":                  v.Description,
				"if_index":              v.Index,
				"physical_entity_index": physInterface.Index,
			})

			if err != nil {
				logger.Error("error saving interfaceinfo in cache","error", err, "hostname", ne.Hostname, "port", v.Description)
				return err
			}

		}
	}
	return nil

}
