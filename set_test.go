// Copyright 2024 TomTonic
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package Set3

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSet3(t *testing.T) {
	t.Run("strings=0", func(t *testing.T) {
		testSet3(t, genStringData(16, 0))
	})
	t.Run("strings=100", func(t *testing.T) {
		testSet3(t, genStringData(16, 101))
	})
	t.Run("strings=1000", func(t *testing.T) {
		testSet3(t, genStringData(16, 1000))
	})
	t.Run("strings=10_000", func(t *testing.T) {
		testSet3(t, genStringData(16, 10_000))
	})
	t.Run("strings=100_000", func(t *testing.T) {
		testSet3(t, genStringData(16, 100_000))
	})
	t.Run("uint32=0", func(t *testing.T) {
		testSet3(t, genUint32Data(0))
	})
	t.Run("uint32=100", func(t *testing.T) {
		testSet3(t, genUint32Data(100))
	})
	t.Run("uint32=1000", func(t *testing.T) {
		testSet3(t, genUint32Data(1000))
	})
	t.Run("uint32=10_000", func(t *testing.T) {
		testSet3(t, genUint32Data(10_000))
	})
	t.Run("uint32=100_000", func(t *testing.T) {
		testSet3(t, genUint32Data(100_000))
	})
	t.Run("string capacity", func(t *testing.T) {
		testSet3Capacity(t, func(n int) []string {
			return genStringData(16, n)
		})
	})
	t.Run("uint32 capacity", func(t *testing.T) {
		testSet3Capacity(t, genUint32Data)
	})
}

func testSet3[K comparable](t *testing.T, keys []K) {
	// sanity check
	require.Equal(t, len(keys), len(uniq(keys)), keys)
	t.Run("put", func(t *testing.T) {
		testSetPut(t, keys)
	})
	t.Run("has", func(t *testing.T) {
		testSetHas(t, keys)
	})
	t.Run("delete", func(t *testing.T) {
		testSetDelete(t, keys)
	})
	t.Run("clear", func(t *testing.T) {
		testSetClear(t, keys)
	})
	t.Run("iter", func(t *testing.T) {
		testSetIter(t, keys)
	})
	t.Run("grow", func(t *testing.T) {
		testSetGrow(t, keys)
	})
}

func uniq[K comparable](keys []K) []K {
	s := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	u := make([]K, 0, len(keys))
	for k := range s {
		u = append(u, k)
	}
	return u
}

func genStringData(size, count int) (keys []string) {
	src := rand.New(rand.NewSource(int64(size * count)))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := make([]rune, size*count)
	for i := range r {
		r[i] = letters[src.Intn(len(letters))]
	}
	keys = make([]string, count)
	for i := range keys {
		keys[i] = string(r[:size])
		r = r[size:]
	}
	return
}

func genUint32Data(count int) (keys []uint32) {
	keys = make([]uint32, count)
	var x uint32
	for i := range keys {
		x += (rand.Uint32() % 128) + 1
		keys[i] = x
	}
	return
}

func testSetPut[K comparable](t *testing.T, keys []K) {
	m := NewSet3[K](uint32(len(keys)))
	assert.Equal(t, 0, m.Count())
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, len(keys), m.Count())
	// overwrite
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, len(keys), m.Count())
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
	assert.Equal(t, len(keys), int(m.resident))
}

func testSetHas[K comparable](t *testing.T, keys []K) {
	m := NewSet3[K](uint32(len(keys)))
	for _, key := range keys {
		m.Add(key)
	}
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
}

func testSetDelete[K comparable](t *testing.T, keys []K) {
	m := NewSet3[K](uint32(len(keys)))
	assert.Equal(t, 0, m.Count())
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, len(keys), m.Count())
	for _, key := range keys {
		m.Remove(key)
		ok := m.Contains(key)
		assert.False(t, ok)
	}
	assert.Equal(t, 0, m.Count())
	// put keys back after deleting them
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, len(keys), m.Count())
}

func testSetClear[K comparable](t *testing.T, keys []K) {
	m := NewSet3[K](0)
	assert.Equal(t, 0, m.Count())
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, len(keys), m.Count())
	m.Clear()
	assert.Equal(t, 0, m.Count())
	for _, key := range keys {
		ok := m.Contains(key)
		assert.False(t, ok)
	}
	var calls int
	m.Iter(func(k K) (stop bool) {
		calls++
		return
	})
	assert.Equal(t, 0, calls)

	// Assert that the Set was actually cleared...
	var k K
	for _, d := range m.group {
		g := d.slot
		for i := range g {
			assert.Equal(t, k, g[i])
		}
	}
}

func testSetIter[K comparable](t *testing.T, keys []K) {
	m := NewSet3[K](uint32(len(keys)))
	for _, key := range keys {
		m.Add(key)
	}
	visited := make(map[K]uint, len(keys))
	m.Iter(func(k K) (stop bool) {
		visited[k] = 0
		stop = true
		return
	})
	if len(keys) == 0 {
		assert.Equal(t, len(visited), 0)
	} else {
		assert.Equal(t, len(visited), 1)
	}
	for _, k := range keys {
		visited[k] = 0
	}
	m.Iter(func(k K) (stop bool) {
		visited[k]++
		return
	})
	for _, c := range visited {
		assert.Equal(t, c, uint(1))
	}
	// mutate on iter
	m.Iter(func(k K) (stop bool) {
		m.Add(k)
		return
	})
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
}

func testSetGrow[K comparable](t *testing.T, keys []K) {
	n := uint32(len(keys))
	m := NewSet3[K](n / 10)
	for _, key := range keys {
		m.Add(key)
	}
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
}

func testSet3Capacity[K comparable](t *testing.T, gen func(n int) []K) {
	// Capacity() behavior depends on |groupSize|
	// which varies by processor architecture.
	caps := []uint32{
		1 * set3maxAvgGroupLoad,
		2 * set3maxAvgGroupLoad,
		3 * set3maxAvgGroupLoad,
		4 * set3maxAvgGroupLoad,
		5 * set3maxAvgGroupLoad,
		10 * set3maxAvgGroupLoad,
		25 * set3maxAvgGroupLoad,
		50 * set3maxAvgGroupLoad,
		100 * set3maxAvgGroupLoad,
	}
	for _, c := range caps {
		m := NewSet3[K](c)
		assert.Equal(t, int(c), m.Capacity())
		keys := gen(rand.Intn(int(c)))
		for _, k := range keys {
			m.Add(k)
		}
		assert.Equal(t, int(c)-len(keys), m.Capacity())
		assert.Equal(t, int(c), m.Count()+m.Capacity())
	}
}
