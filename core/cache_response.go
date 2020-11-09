package core

import (
	"context"
	"errors"
	
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

type CachedResponse struct {
	Hostname    string                                 `bson:"hostname"`
	Port        string                                 `bson:"port"`
	RequestType pb_core.Request_Type                   `bson:"request_type"`
	Response    *pb_core.Response                      `bson:"response,omitempty"`
	Timestamp   *timestamp.Timestamp                   `bson:"timestamp" json:"timestamp"`
}

func (c *cache) PopResponse(hostname, port string, requestType pb_core.Request_Type) (*CachedResponse, error) {
	res := c.col.FindOne(context.Background(), bson.M{"hostname": hostname, "port": port, "request_type": requestType})
	obj := &CachedResponse{}

	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return obj, nil
}

func (c *cache) SetResponse(hostname, port string, requestType pb_core.Request_Type, response *pb_core.Response) error {
	_, err := c.col.InsertOne(context.Background(), &CachedResponse{
		hostname,
		port,
		requestType,
		response,
		ptypes.TimestampNow(),
	})

	if err != nil {
		c.logger.Error("error saving response in cache","error",  err)
		return err
	}

	return nil
}

func (c *cache) Clear(hostname, port string, requestType pb_core.Request_Type) error {
	_, err := c.col.DeleteMany(context.Background(), bson.M{"hostname": hostname, "port": port, "request_type": requestType})
	return err
}
