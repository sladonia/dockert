package container

import (
	"context"
	"errors"
	"net"

	"github.com/nsqio/go-nsq"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/docker"
)

func NewNSQD() docker.Container {
	return docker.NewCommonContainer(
		"",
		"4150/tcp",
		&dockertest.RunOptions{
			Name:       "nsqd",
			Repository: "nsqio/nsq",
			Tag:        "v1.2.1",
			Cmd:        []string{"nsqd"},
		},
		docker.ReadinessCheckerFunc(func(ctx context.Context, c docker.Container) (bool, error) {
			handler := nsq.HandlerFunc(func(message *nsq.Message) error {
				return nil
			})

			config := nsq.NewConfig()

			consumer, err := nsq.NewConsumer("topic", "channel", config)
			if err != nil {
				return false, err
			}

			consumer.SetLogger(&NopLogger{}, nsq.LogLevelMax)
			consumer.AddHandler(handler)

			err = consumer.ConnectToNSQD(c.DSN())
			if err != nil {
				var connectionError *net.OpError

				if errors.As(err, &connectionError) {
					return false, nil
				}

				return false, err
			}

			return true, nil
		}),
	)
}

type NopLogger struct{}

func (n *NopLogger) Output(_ int, _ string) error {
	return nil
}
