package core

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CachedResponse struct {
	Hostname    string      `bson:"hostname"`
	Port        string      `bson:"port"`
	RequestType RequestType `bson:"request_type"`
	Response    interface{} `bson:"response,omitempty"`
	Timestamp   time.Time   `bson:"timestamp,omitempty"`
}

func (c *cache) PopResponse(hostname, port string, requestType RequestType) (*BSONResponse, error) {
	res := c.col.FindOne(context.Background(), bson.M{"hostname": hostname, "port": port, "request_type": requestType})
	obj := &BSONResponse{}
	if err := res.Decode(obj); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return obj, nil
}

func (c *cache) SetResponse(hostname, port string, requestType RequestType, response interface{}) error {
	var _, err = c.col.InsertOne(context.Background(), &CachedResponse{
		hostname,
		port,
		requestType,
		response,
		time.Now(),
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
