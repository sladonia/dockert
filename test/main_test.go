package test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/docker"
	"github.com/sladonia/docker/container"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Suite struct {
	suite.Suite

	pool *dockertest.Pool
}

func (s *Suite) SetupSuite() {
	var err error

	s.pool, err = dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
}

func (s *Suite) TearDownSuite() {
}

func (s *Suite) TestMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c := container.NewMongo()

	err := c.Start(s.pool)
	s.Require().NoError(err)

	dsn := c.DSN()

	_, err = url.Parse(dsn)
	s.Require().NoError(err)

	err = c.WaitReady(ctx)
	s.Require().NoError(err)

	err = c.Stop()
	s.Require().NoError(err)
}

func (s *Suite) TestRedis() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c := container.NewRedis()

	err := c.Start(s.pool)
	s.Require().NoError(err)

	err = c.WaitReady(ctx)
	s.Require().NoError(err)

	err = c.Stop()
	s.Require().NoError(err)
}

func (s *Suite) TestNats() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c := container.NewNats()

	err := c.Start(s.pool)
	s.Require().NoError(err)

	dsn := c.DSN()

	_, err = url.Parse(dsn)
	s.Require().NoError(err)

	err = c.WaitReady(ctx)
	s.Require().NoError(err)

	err = c.Stop()
	s.Require().NoError(err)
}

func (s *Suite) TestRegistry() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r := docker.NewRegistry(s.pool).
		Add(container.NewMongo()).
		Add(container.NewRedis()).
		Add(container.NewNats())

	err := r.StartAndWaitReady(ctx)
	s.Require().NoError(err)

	// Connect to mongo
	mongoContainer, ok := r.ByName("mongo")
	s.Require().True(ok)

	clientOptions := options.Client().ApplyURI(mongoContainer.DSN())
	client, err := mongo.Connect(ctx, clientOptions)
	s.Require().NoError(err)
	err = client.Ping(ctx, nil)
	s.Require().NoError(err)

	// Connect to redis
	redisContainer, ok := r.ByName("redis")
	s.Require().NoError(err)

	redisClientOptions := &redis.Options{Addr: redisContainer.DSN()}
	redisClient := redis.NewClient(redisClientOptions)
	_, err = redisClient.Ping(ctx).Result()
	s.Require().NoError(err)

	// Connect to NATS
	natsContainer, ok := r.ByName("nats")
	s.Require().NoError(err)

	_, err = nats.Connect(natsContainer.DSN())
	s.Require().NoError(err)

	err = r.Stop()
	s.Require().NoError(err)
}

func (s *Suite) TestNSQD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c := container.NewNSQD()
	err := c.Start(s.pool)
	s.Require().NoError(err)

	err = c.WaitReady(ctx)
	s.Require().NoError(err)

	err = c.Stop()
	s.Require().NoError(err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
