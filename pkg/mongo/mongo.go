package mongo

import (
	"context"

	logx "github.com/chai-rs/sevenhunter/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI string
}

func (conf *Config) New(ctx context.Context) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(conf.URI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func (conf *Config) MustNew(ctx context.Context) *mongo.Client {
	client, err := conf.New(ctx)
	if err != nil {
		logx.Panic().Err(err).Msg("failed to initialize mongo client")
	}
	return client
}
