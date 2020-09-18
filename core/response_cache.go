package core

import (
	"context"
	"errors"
	"git.liero.se/opentelco/go-swpx/proto/resource"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type CachedResponse struct {
	Hostname    string                                 `bson:"hostname"`
	Port        string                                 `bson:"port"`
	RequestType RequestType                            `bson:"request_type"`
	Response    *resource.TechnicalInformationResponse `bson:"response,omitempty"`
	Timestamp   *timestamp.Timestamp                 `bson:"timestamp" json:"timestamp"`
}

func (c *cache) PopResponse(hostname, port string, requestType RequestType) (*CachedResponse, error) {
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

func (c *cache) SetResponse(hostname, port string, requestType RequestType, response *resource.TechnicalInformationResponse) error {
	_, err := c.col.InsertOne(context.Background(), &CachedResponse{
		hostname,
		port,
		requestType,
		response,
		ptypes.TimestampNow(),
	})

	if err != nil {
		c.logger.Error("error saving response in cache: ", err.Error())
		return err
	}

	return nil
}

func (c *cache) Clear(hostname, port string, requestType RequestType) error {
	_, err := c.col.DeleteMany(context.Background(), bson.M{"hostname": hostname, "port": port, "request_type": requestType})

	return err
}
