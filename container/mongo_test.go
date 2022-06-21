package container

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	container := NewMongo()

	err = container.Run(pool)
	require.NoError(t, err)

	dsn := container.DSN()

	_, err = url.Parse(dsn)
	require.NoError(t, err)

	err = container.WaitReady(ctx)
	require.NoError(t, err)

	err = container.Stop()
	require.NoError(t, err)
}
