package container

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/docker"
)

func NewRedis() docker.Container {
	return docker.NewCommonContainer(
		"",
		"6379/tcp",
		&dockertest.RunOptions{
			Name:       "redis",
			Repository: "redis",
			Tag:        "7-alpine",
		},
		docker.ReadinessCheckerFunc(func(ctx context.Context, c docker.Container) (bool, error) {
			clientOptions := &redis.Options{Addr: c.DSN()}

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
