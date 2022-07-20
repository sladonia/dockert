package container

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	docker "github.com/sladonia/dockert"
)

const (
	PortNats = "4222"
)

func NewNats() docker.Container {
	return docker.NewCommonContainer(
		&dockertest.RunOptions{
			Name:       "nats",
			Repository: "nats",
			Tag:        "2-alpine",
			Cmd:        []string{"nats-server", "-js"},
		},
		docker.ReadinessCheckerFunc(func(ctx context.Context, c docker.Container) (bool, error) {
			conn, err := nats.Connect(NatsDSN(c))
			if err != nil {
				return false, nil
			}

			conn.Close()

			return true, nil
		}),
	)
}

func NatsDSN(c docker.Container) string {
	return c.Address("nats", PortNats)
}
