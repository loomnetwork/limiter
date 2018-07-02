package memory_test

import (
	"testing"
	"github.com/loomnetwork/limiter/drivers/store/memory"
	"github.com/loomnetwork/limiter/drivers/store/tests"
)

func TestMemoryStoreReset(t *testing.T) {
	tests.TestLimiter_Reset(t, memory.NewStore())
}
