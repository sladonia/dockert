package container

import (
	"context"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/dockert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PortMongo = "27017"
)

func NewMongo() dockert.Container {
	return dockert.NewCommonContainer(
		&dockertest.RunOptions{
			Name:       "mongo",
			Repository: "mongo",
			Tag:        "5.0",
		},
		dockert.ReadinessCheckerFunc(func(ctx context.Context, c dockert.Container) (bool, error) {
			clientOptions := options.Client().ApplyURI(MongoDSN(c))
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

func MongoDSN(c dockert.Container) string {
	return c.Address("mongodb", PortMongo)
}
