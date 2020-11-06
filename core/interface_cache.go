package core

import (
	"context"
	"errors"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (c *cache) PopInterface(hostname, iface string) (*CachedInterface, error) {
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

func (c *cache) SetInterface(ne *proto.NetworkElement, interfaces *proto.NetworkElementInterfaces, phys *proto.NetworkElementInterfaces) error {
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
