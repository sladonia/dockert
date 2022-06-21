package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
)

type Container interface {
	Name() string
	DSN() string
	Run(pool *dockertest.Pool) error
	WaitReady(ctx context.Context) error
	IsReady() bool
	DependsOn(other Container) Container
	Resource() *dockertest.Resource
}

type ReadinessWaiter interface {
	Wait(ctx context.Context, self Container) error
}

type ReadinessWaiterFunc func(ctx context.Context, self Container) error

func (f ReadinessWaiterFunc) Wait(ctx context.Context, self Container) error {
	return f(ctx, self)
}

type commonContainer struct {
	urlScheme       string
	defaultPort     string
	name            string
	options         *dockertest.RunOptions
	resource        *dockertest.Resource
	dependsOn       []Container
	isReady         bool
	readinessWaiter ReadinessWaiter
}

func NewCommonContainer(
	urlScheme, defaultPort string,
	options *dockertest.RunOptions,
	readinessWaiter ReadinessWaiter,
) Container {
	return &commonContainer{
		urlScheme:       urlScheme,
		defaultPort:     defaultPort,
		name:            options.Name,
		options:         options,
		readinessWaiter: readinessWaiter,
	}
}

func (c *commonContainer) Name() string {
	return c.name
}

func (c *commonContainer) DSN() string {
	addr := c.resource.GetPort(c.defaultPort)

	if IsDarwinOS() && IsRunningInDockerContainer() {
		addr = fmt.Sprintf("127.0.0.1:%s", addr)
	} else {
		addr = fmt.Sprintf("%s:%s", c.resource.Container.NetworkSettings.Gateway, addr)
	}

	return fmt.Sprintf("%s://%s", c.urlScheme, addr)
}

func (c *commonContainer) Run(pool *dockertest.Pool) error {
	var err error

	resource, ok := pool.ContainerByName(c.name)
	if ok {
		err = pool.Purge(resource)
		if err != nil {
			return err
		}
	}

	resource, err = pool.RunWithOptions(c.options)
	if err != nil {
		return err
	}

	c.resource = resource

	return nil
}

func (c *commonContainer) WaitReady(ctx context.Context) error {
	t := time.NewTimer(0)

	// wait for
timerLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			for _, otherContainer := range c.dependsOn {
				if !otherContainer.IsReady() {
					t.Reset(100 * time.Millisecond)
					continue timerLoop
				}
			}

			break timerLoop
		}
	}

	return c.readinessWaiter.Wait(ctx, c)
}

func (c *commonContainer) DependsOn(other Container) Container {
	c.dependsOn = append(c.dependsOn, other)
	return c
}

func (c *commonContainer) Resource() *dockertest.Resource {
	return c.resource
}

func (c *commonContainer) IsReady() bool {
	return c.isReady
}
