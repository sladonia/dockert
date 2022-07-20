package dockert

import (
	"context"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
)

type Container interface {
	Name() string
	Address(scheme string, port string) string
	Start(ctx context.Context, pool *dockertest.Pool) error
	Stop() error
	WaitReady(ctx context.Context) error
	IsReady() bool
	DependsOn(other Container) Container
	Resource() *dockertest.Resource
}

type ReadinessChecker interface {
	IsReady(ctx context.Context, self Container) (bool, error)
}

type ReadinessCheckerFunc func(ctx context.Context, c Container) (bool, error)

func (f ReadinessCheckerFunc) IsReady(ctx context.Context, c Container) (bool, error) {
	return f(ctx, c)
}

type commonContainer struct {
	name            string
	options         *dockertest.RunOptions
	resource        *dockertest.Resource
	dependsOn       []Container
	isReady         bool
	readinessWaiter ReadinessChecker
}

func NewCommonContainer(
	options *dockertest.RunOptions,
	readinessWaiter ReadinessChecker,
) Container {
	return &commonContainer{
		name:            options.Name,
		options:         options,
		readinessWaiter: readinessWaiter,
	}
}

func (c *commonContainer) Name() string {
	return c.name
}

func (c *commonContainer) Address(scheme, port string) string {
	var addr string

	containerPort := c.resource.GetPort(fmt.Sprintf("%s/tcp", port))

	if IsDarwinOS() && IsRunningInDockerContainer() {
		addr = fmt.Sprintf("127.0.0.1:%s", containerPort)
	} else {
		addr = fmt.Sprintf("%s:%s", c.resource.Container.NetworkSettings.Gateway, containerPort)
	}

	if scheme == "" {
		return addr
	}

	return fmt.Sprintf("%s://%s", scheme, addr)
}

func (c *commonContainer) Start(ctx context.Context, pool *dockertest.Pool) error {
	err := c.waitForOtherContainers(ctx)
	if err != nil {
		return err
	}

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

			c.isReady = true
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
