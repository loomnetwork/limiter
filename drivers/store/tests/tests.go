package tests

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/loomnetwork/limiter"
	"math/rand"
)

// TestStoreSequentialAccess verify that store works as expected with a sequential access.
func TestStoreSequentialAccess(t *testing.T, store limiter.Store) {
	is := require.New(t)
	ctx := context.Background()

	limiter := limiter.New(store, limiter.Rate{
		Limit:  3,
		Period: time.Minute,
	})

	for i := 1; i <= 6; i++ {

		if i <= 3 {

			lctx, err := limiter.Peek(ctx, "foo")
			is.NoError(err)
			is.NotZero(lctx)
			is.Equal(int64(3-(i-1)), lctx.Remaining)

		}

		lctx, err := limiter.Get(ctx, "foo")
		is.NoError(err)
		is.NotZero(lctx)

		if i <= 3 {

			is.Equal(int64(3), lctx.Limit)
			is.Equal(int64(3-i), lctx.Remaining)
			is.True(math.Ceil(time.Since(time.Unix(lctx.Reset, 0)).Seconds()) <= 60)

			lctx, err = limiter.Peek(ctx, "foo")
			is.NoError(err)
			is.Equal(int64(3-i), lctx.Remaining)

		} else {

			is.Equal(int64(3), lctx.Limit)
			is.True(lctx.Remaining == 0)
			is.True(math.Ceil(time.Since(time.Unix(lctx.Reset, 0)).Seconds()) <= 60)

		}
	}
}

// TestStoreConcurrentAccess verify that store works as expected with a concurrent access.
func TestStoreConcurrentAccess(t *testing.T, store limiter.Store) {
	is := require.New(t)
	ctx := context.Background()

	limiter := limiter.New(store, limiter.Rate{
		Limit:  100000,
		Period: 10 * time.Second,
	})

	goroutines := 500
	ops := 500

	wg := &sync.WaitGroup{}
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			for j := 0; j < ops; j++ {
				lctx, err := limiter.Get(ctx, "foo")
				is.NoError(err)
				is.NotZero(lctx)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func getNewLimiter(limit int64, period time.Duration, store limiter.Store) *limiter.Limiter {
	return limiter.New(store, limiter.Rate{
		Limit:  limit,
		Period: period * time.Second,
	})

}

// Test Limiter reset functionality
func TestLimiter_Reset(t *testing.T, store limiter.Store) {
	is := require.New(t)


	scenarios := []struct {
		lim *limiter.Limiter
		ctx context.Context
		key string
		newRate int64
		expected int64
	}{
		{
			lim: getNewLimiter(100, 20, store),
			key: "1",
			newRate: int64(120),
			expected: 120,
		},
		{
			lim: getNewLimiter(3, 1000, store),
			key: "2",
			newRate: int64(761),
			expected: 761,
		},
		{
			lim: getNewLimiter(213, 120, store),
			key: "3",
			newRate: int64(1),
			expected: 1,
		},
		{
			lim: getNewLimiter(10, 4, store),
			key: "4",
			newRate: int64(2101),
			expected: 2101,
		},
	}

	for _, v := range scenarios {
		var lctx limiter.Context
		for i := 0; i < 1 + rand.Intn(int(v.lim.Rate.Limit)); i++ {
			lctx, _ = v.lim.Get(v.ctx, v.key)
		}
		lctx, _ = v.lim.Reset(v.ctx, v.key, v.newRate)
		is.True(v.newRate == lctx.Remaining)
	}
}
