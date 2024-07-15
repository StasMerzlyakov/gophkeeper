package ttlstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/ttlstorage"
	"github.com/stretchr/testify/require"
)

func TestTTLMap(t *testing.T) {

	t.Run("simple_operations_ttl", func(t *testing.T) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		ttl := ttlstorage.NewTTLMap(ctx, 3*time.Second)
		ttl.Store("key1", 1)
		ttl.Store("key2", 2)
		ttl.Store("key3", "hello")

		checkValue(t, ttl, "key1", 1)
		checkValue(t, ttl, "key2", 2)
		checkValue(t, ttl, "key3", "hello")

		ttl.Delete("key1")
		checkValue(t, ttl, "key1", nil)
		checkValue(t, ttl, "key2", 2)
		checkValue(t, ttl, "key3", "hello")

		time.Sleep(4 * time.Second)

		checkValue(t, ttl, "key1", nil)
		checkValue(t, ttl, "key2", nil)
		checkValue(t, ttl, "key3", nil)
	})

	t.Run("sync_operations", func(t *testing.T) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		ttl := ttlstorage.NewTTLMap(ctx, 3*time.Second)

		val, ok := ttl.LoadOrStore("key1", 1)
		require.False(t, ok)
		require.Equal(t, 1, val)
		checkValue(t, ttl, "key1", 1)

		val, ok = ttl.LoadOrStore("key1", 2)
		require.True(t, ok)
		require.Equal(t, 1, val)
		checkValue(t, ttl, "key1", 1)

		val, ok = ttl.LoadAndDelete("key1")
		require.True(t, ok)
		require.Equal(t, 1, val)
		checkValue(t, ttl, "key1", nil)

		val, ok = ttl.LoadAndDelete("key1")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func checkValue(t *testing.T, tMap *ttlstorage.TTLMap, key string, expected any) {
	val, ok := tMap.Load(key)
	if expected == nil {
		require.False(t, ok)
		require.Nil(t, val)
	} else {
		require.True(t, ok)
		require.Equal(t, expected, val)
	}
}
