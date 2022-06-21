package test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/docker/container"
	"github.com/stretchr/testify/suite"
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

	err := c.Run(s.pool)
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

	err := c.Run(s.pool)
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

	err := c.Run(s.pool)
	s.Require().NoError(err)

	dsn := c.DSN()

	_, err = url.Parse(dsn)
	s.Require().NoError(err)

	err = c.WaitReady(ctx)
	s.Require().NoError(err)

	err = c.Stop()
	s.Require().NoError(err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
