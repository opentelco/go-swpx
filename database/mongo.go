package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelco.io/go-swpx/config"
)

func New(conf config.MongoDb, logger hclog.Logger) (*mongo.Client, error) {
	logger.Info("attempting to connect to mongoDb server", "server", conf.Addr)

	reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	opts := options.Client()

	if conf.User != "" {
		opts.SetAuth(options.Credential{
			Username: conf.User,
			Password: conf.Password,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout.AsDuration())

	defer func() {
		cancel()
	}()

	addr := fmt.Sprintf("mongodb://%s:%d", conf.Addr, conf.Port)
	opts.ApplyURI(addr)
	opts.Registry = reg
	client, err := mongo.NewClient(opts)

	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Info("Connecting to Mongo DB",
			"url", addr,
			"username", conf.User,
			"password", hidePassword(conf.Password),
		)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, nil
}

func hidePassword(input string) string {
	chars := len(input)

	if chars < 4 {
		return input
	}

	show := chars / 2
	if show > 4 {
		show = 4
	}

	o := strings.Repeat("*", chars-show) + string(input[chars-show:])
	fmt.Println(o)
	return o
}
