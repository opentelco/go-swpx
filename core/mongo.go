package core

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"git.liero.se/opentelco/go-swpx/shared"
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongoDb(conf shared.ConfigMongo, logger hclog.Logger) (*mongo.Client, error) {
	logger.Info("attempting to connect to mongoDb server", "server", conf.Server)

	reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	opts := options.Client()

	if conf.Username != "" {
		opts.SetAuth(options.Credential{
			Username: conf.Username,
			Password: conf.Password,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeoutSeconds)*time.Second)

	defer func() {
		cancel()
	}()

	opts.ApplyURI(conf.Server)
	opts.Registry = reg
	client, err := mongo.NewClient(opts)

	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Info("Connecting to Mongo DB",
			"url", conf.Server,
			"username", conf.Username,
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
