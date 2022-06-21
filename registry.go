package docker

import "context"

type Registry interface {
	Add(container Container)
	StartAndWaitReady(ctx context.Context) error
	Stop() error
}
