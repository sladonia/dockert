package container

import (
	"context"

	"github.com/ory/dockertest/v3"
	"github.com/sladonia/dockert"
	etcd "go.etcd.io/etcd/client/v3"
)

const (
	PortEtcd = "2379"
)

func NewEtcd() dockert.Container {
	return dockert.NewCommonContainer(
		&dockertest.RunOptions{
			Name:       "etcd",
			Repository: "docker.io/bitnami/etcd",
			Tag:        "3.5",
			Env:        []string{"ALLOW_NONE_AUTHENTICATION=yes"},
		},
		dockert.ReadinessCheckerFunc(func(ctx context.Context, c dockert.Container) (bool, error) {
			key := "test"

			client, err := etcd.New(etcd.Config{
				Endpoints: []string{EtcdDSN(c)},
			})
			if err != nil {
				return false, nil
			}

			_, err = client.Put(ctx, key, "test")
			if err != nil {
				return false, err
			}

			_, err = client.Delete(ctx, key)
			if err != nil {
				return false, err
			}

			return true, nil
		}),
	)
}

func EtcdDSN(c dockert.Container) string {
	return c.Address("", PortEtcd)
}
