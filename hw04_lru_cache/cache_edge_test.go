package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCacheCapacityEdgeCases(t *testing.T) {
	t.Run("zero capacity", func(t *testing.T) {
		c := NewCache(0)

		require.False(t, c.Set("key", "value"))

		value, ok := c.Get("key")
		require.False(t, ok)
		require.Nil(t, value)
	})

	t.Run("negative capacity", func(t *testing.T) {
		c := NewCache(-1)

		require.False(t, c.Set("key", "value"))
		_, ok := c.Get("key")
		require.False(t, ok)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(2)
		c.Set("a", 1)
		c.Set("b", 2)

		c.Clear()

		require.Equal(t, 0, c.(*lruCache).queue.Len())
		_, ok := c.Get("a")
		require.False(t, ok)
	})
}
