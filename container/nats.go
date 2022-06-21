package container

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/docker"
)

func NewNats() docker.Container {
	return docker.NewCommonContainer(
		"nats",
		"4222/tcp",
		&dockertest.RunOptions{
			Name:       "nats",
			Repository: "nats",
			Tag:        "2-alpine",
			Cmd:        []string{"nats-server", "-js"},
		},
		docker.ReadinessCheckerFunc(func(ctx context.Context, c docker.Container) (bool, error) {
			conn, err := nats.Connect(c.DSN())
			if err != nil {
				return false, nil
			}

			conn.Close()

			return true, nil
		}),
	)
}
