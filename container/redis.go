package container

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/dockert"
)

const (
	PortRedis = "6379"
)

func NewRedis() dockert.Container {
	return dockert.NewCommonContainer(
		&dockertest.RunOptions{
			Name:       "redis",
			Repository: "redis",
			Tag:        "7-alpine",
		},
		dockert.ReadinessCheckerFunc(func(ctx context.Context, c dockert.Container) (bool, error) {
			clientOptions := &redis.Options{Addr: RedisDSN(c)}

			client := redis.NewClient(clientOptions)

			_, err := client.Ping(ctx).Result()
			if err != nil {
				fmt.Println(err)
				return false, nil
			}

			err = client.Close()
			if err != nil {
				return false, err
			}

			return true, nil
		}),
	)
}

func RedisDSN(c dockert.Container) string {
	return c.Address("", PortRedis)
}
