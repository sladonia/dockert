package dockert

import (
	"context"
	"sync"

	"github.com/ory/dockertest/v3"
	"golang.org/x/sync/errgroup"
)

type Registry struct {
	sync.RWMutex

	pool       *dockertest.Pool
	containers map[string]Container
}

func NewRegistry(pool *dockertest.Pool) *Registry {
	return &Registry{
		pool:       pool,
		containers: make(map[string]Container),
	}
}

func (r *Registry) Add(container Container) *Registry {
	r.Lock()
	defer r.Unlock()

	r.containers[container.Name()] = container

	return r
}

func (r *Registry) ByName(name string) (Container, bool) {
	r.RLock()
	defer r.RUnlock()

	c, ok := r.containers[name]

	return c, ok
}

func (r *Registry) StartAndWaitReady(ctx context.Context) error {
	r.RLock()
	defer r.RUnlock()

	g, ctx := errgroup.WithContext(ctx)

	for _, container := range r.containers {
		container := container
		g.Go(func() error {
			err := container.Start(ctx, r.pool)
			if err != nil {
				return err
			}

			return container.WaitReady(ctx)
		})
	}

	return g.Wait()
}

func (r *Registry) Stop() error {
	r.RLock()
	defer r.RUnlock()

	g := &errgroup.Group{}

	for _, container := range r.containers {
		container := container

		g.Go(func() error {
			return container.Stop()
		})
	}

	return g.Wait()
}
