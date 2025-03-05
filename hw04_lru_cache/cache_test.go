package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Set("bbb", 200) // [200, 100]
		c.Set("ccc", 300) // [300, 200, 100]
		c.Set("ddd", 400) // [400, 300, 200]

		aaa, ok := c.Get("aaa") // этого элемента уже нет
		assert.False(t, ok)
		assert.Nil(t, aaa)
		bbb, ok := c.Get("bbb") // [200, 400, 300]
		assert.True(t, ok)
		assert.Equal(t, 200, bbb)
		ccc, ok := c.Get("ccc") // [300, 200, 400]
		assert.True(t, ok)
		assert.Equal(t, 300, ccc)
		ddd, ok := c.Get("ddd") // [400, 300, 200]
		assert.True(t, ok)
		assert.Equal(t, 400, ddd)

		c.Set("ddd", 401) // [401, 300, 200]
		c.Set("ddd", 404) // [404, 300, 200]
		c.Set("bbb", 202) // [202, 404, 300]
		c.Set("eee", 500) // [500, 202, 404]

		ccc, ok = c.Get("ccc")
		assert.False(t, ok)
		assert.Nil(t, ccc)
		ddd, ok = c.Get("ddd")
		assert.True(t, ok)
		assert.Equal(t, 404, ddd)
		bbb, ok = c.Get("bbb")
		assert.True(t, ok)
		assert.Equal(t, 202, bbb)
		eee, ok := c.Get("eee")
		assert.True(t, ok)
		assert.Equal(t, 500, eee)

		c.Clear()
		bbb, ok = c.Get("bbb")
		assert.False(t, ok)
		assert.Nil(t, bbb)
		eee, ok = c.Get("eee")
		assert.False(t, ok)
		assert.Nil(t, eee)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
