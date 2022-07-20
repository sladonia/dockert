package main

import (
	"context"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/dockert/container"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Create a pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Create container
	c := container.NewMongo()
	defer c.Stop()

	// Start container
	err = c.Start(ctx, pool)
	if err != nil {
		panic(err)
	}

	// Wait container ready
	err = c.WaitReady(ctx)
	if err != nil {
		panic(err)
	}

	// Use this data source name to connect to mongodb
	dsn := container.MongoDSN(c)

	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	// Perform tests against mongodb
}
