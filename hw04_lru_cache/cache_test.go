package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		checkGet(t, c, "aaa", 100, true)

		checkGet(t, c, "bbb", 200, true)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		checkGet(t, c, "aaa", 300, true)

		checkGet(t, c, "ccc", nil, false)
	})

	t.Run("clear test", func(t *testing.T) {
		c := NewCache(4)
		c.Set("1", 1)
		c.Set("2", 2)

		checkGet(t, c, "1", 1, true)
		checkGet(t, c, "2", 2, true)
		lru := c.(*lruCache)
		require.Equal(t, lru.queue.Len(), 2)
		require.Equal(t, len(lru.items), 2)

		c.Clear()

		checkGet(t, c, "1", nil, false)
		checkGet(t, c, "2", nil, false)
		require.Equal(t, lru.queue.Len(), 0)
		require.Equal(t, len(lru.items), 0)
	})

	t.Run("overload test", func(t *testing.T) {
		c := NewCache(3)
		lru := c.(*lruCache)

		c.Set("1", 1) // [1]
		c.Set("2", 2) // [2,1]
		c.Set("3", 3) // [3,2,1]
		checkCapacity(t, lru, 3)

		c.Set("4", 4) // [4,3,2]
		checkCapacity(t, lru, 3)
		checkGet(t, c, "4", 4, true) // [4,3,2]
		checkGet(t, c, "3", 3, true) // [3,4,2]
		checkGet(t, c, "2", 2, true) // [2,3,4]
		checkGet(t, c, "1", nil, false)

		c.Set("4", 44) // [44,2,3]
		c.Set("3", 33) // [33,44,2]
		c.Set("5", 5)  // [5,33,44]

		checkGet(t, c, "1", nil, false)
		checkGet(t, c, "2", nil, false)
		checkGet(t, c, "3", 33, true)
		checkGet(t, c, "4", 44, true)
		checkGet(t, c, "5", 5, true)
		checkCapacity(t, lru, 3)

		c.Clear()

		checkCapacity(t, lru, 0)
		checkGet(t, c, "3", nil, false)
		checkGet(t, c, "4", nil, false)
		checkGet(t, c, "5", nil, false)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func checkGet(t *testing.T, c Cache, key Key, expectedValue interface{}, exist bool) {
	t.Helper()
	val, ok := c.Get(key)
	require.Equal(t, ok, exist)
	require.Equal(t, expectedValue, val)
}

func checkCapacity(t *testing.T, lru *lruCache, length int) {
	t.Helper()
	require.Equal(t, lru.queue.Len(), length)
	require.Equal(t, len(lru.items), length)
}
