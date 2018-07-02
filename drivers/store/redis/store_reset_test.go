package redis_test

import (
	"testing"
	"github.com/loomnetwork/limiter/drivers/store/redis"
	"github.com/loomnetwork/limiter/drivers/store/tests"
	"github.com/stretchr/testify/require"
	"github.com/loomnetwork/limiter"
)

func TestRedisStoreReset(t *testing.T) {
	is := require.New(t)
	client, err := newRedisClient()
	is.NoError(err)
	is.NotNil(client)

	var store limiter.Store
	store, _ = redis.NewStore(client)

	tests.TestLimiter_Reset(t, store)
}
