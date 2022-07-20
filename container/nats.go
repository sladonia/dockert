package container

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/dockert"
)

const (
	PortNats = "4222"
)

func NewNats() dockert.Container {
	return dockert.NewCommonContainer(
		&dockertest.RunOptions{
			Name:       "nats",
			Repository: "nats",
			Tag:        "2-alpine",
			Cmd:        []string{"nats-server", "-js"},
		},
		dockert.ReadinessCheckerFunc(func(ctx context.Context, c dockert.Container) (bool, error) {
			conn, err := nats.Connect(NatsDSN(c))
			if err != nil {
				return false, nil
			}

			conn.Close()

			return true, nil
		}),
	)
}

func NatsDSN(c dockert.Container) string {
	return c.Address("nats", PortNats)
}
