package container

import (
	"context"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo() docker.Container {
	return docker.NewCommonContainer(
		"mongodb",
		"27017/tcp",
		&dockertest.RunOptions{
			Name:       "mongo",
			Repository: "mongo",
			Tag:        "5.0",
		},
		docker.ReadinessCheckerFunc(func(ctx context.Context, c docker.Container) (bool, error) {
			clientOptions := options.Client().ApplyURI(c.DSN())
			client, err := mongo.Connect(ctx, clientOptions)
			if err != nil {
				return false, err
			}

			err = client.Ping(ctx, nil)
			if err != nil {
				return false, nil
			}

			return true, nil
		}),
	)
}
