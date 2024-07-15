package ttlstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/ttlstorage"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestTTLMap(t *testing.T) {

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ttl := ttlstorage.NewTTLMap(ctx, 3*time.Second)
	ttl.Store("key1", 1)
	ttl.Store("key2", 2)
	ttl.Store("key3", "hello")

	require.Equal(t, 1, ttl.Load("key1"))
	require.Equal(t, 2, ttl.Load("key2"))
	require.Equal(t, "hello", ttl.Load("key3"))

	ttl.Delete("key1")
	assert.Equal(t, nil, ttl.Load("key1"))
	require.Equal(t, 2, ttl.Load("key2"))
	require.Equal(t, "hello", ttl.Load("key3"))

	time.Sleep(4 * time.Second)

	assert.Equal(t, nil, ttl.Load("key1"))
	assert.Equal(t, nil, ttl.Load("key2"))
	assert.Equal(t, nil, ttl.Load("key3"))
}
