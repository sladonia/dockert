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
	Stop() error
	WaitReady(ctx context.Context) error
	IsReady() bool
	DependsOn(other Container) Container
	Resource() *dockertest.Resource
}

type ReadinessChecker interface {
	IsReady(ctx context.Context, self Container) (bool, error)
}

type ReadinessCheckerFunc func(ctx context.Context, self Container) (bool, error)

func (f ReadinessCheckerFunc) IsReady(ctx context.Context, self Container) (bool, error) {
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
	readinessWaiter ReadinessChecker
}

func NewCommonContainer(
	urlScheme, defaultPort string,
	options *dockertest.RunOptions,
	readinessWaiter ReadinessChecker,
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
	var addr string

	port := c.resource.GetPort(c.defaultPort)

	if IsDarwinOS() && IsRunningInDockerContainer() {
		addr = fmt.Sprintf("127.0.0.1:%s", port)
	} else {
		addr = fmt.Sprintf("%s:%s", c.resource.Container.NetworkSettings.Gateway, port)
	}

	if c.urlScheme == "" {
		return addr
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
	err := c.waitForOtherContainers(ctx)
	if err != nil {
		return err
	}

	t := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			ok, err := c.readinessWaiter.IsReady(ctx, c)
			if err != nil {
				return err
			}

			if !ok {
				t.Reset(100 * time.Millisecond)
				continue
			}

			return nil
		}
	}
}

func (c *commonContainer) waitForOtherContainers(ctx context.Context) error {
	t := time.NewTimer(0)

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

			return nil
		}
	}
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

func (c *commonContainer) Stop() error {
	return c.resource.Close()
}
