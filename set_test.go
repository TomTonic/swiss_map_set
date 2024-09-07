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

package swiss

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSwissSet(t *testing.T) {
	t.Run("strings=0", func(t *testing.T) {
		testSwissSet(t, genStringData(16, 0))
	})
	t.Run("strings=100", func(t *testing.T) {
		testSwissSet(t, genStringData(16, 101))
	})
	t.Run("strings=1000", func(t *testing.T) {
		testSwissSet(t, genStringData(16, 1000))
	})
	t.Run("strings=10_000", func(t *testing.T) {
		testSwissSet(t, genStringData(16, 10_000))
	})
	t.Run("strings=100_000", func(t *testing.T) {
		testSwissSet(t, genStringData(16, 100_000))
	})
	t.Run("uint32=0", func(t *testing.T) {
		testSwissSet(t, genUint32Data(0))
	})
	t.Run("uint32=100", func(t *testing.T) {
		testSwissSet(t, genUint32Data(100))
	})
	t.Run("uint32=1000", func(t *testing.T) {
		testSwissSet(t, genUint32Data(1000))
	})
	t.Run("uint32=10_000", func(t *testing.T) {
		testSwissSet(t, genUint32Data(10_000))
	})
	t.Run("uint32=100_000", func(t *testing.T) {
		testSwissSet(t, genUint32Data(100_000))
	})
	t.Run("string capacity", func(t *testing.T) {
		testSwissSetCapacity(t, func(n int) []string {
			return genStringData(16, n)
		})
	})
	t.Run("uint32 capacity", func(t *testing.T) {
		testSwissSetCapacity(t, genUint32Data)
	})
}

func testSwissSet[K comparable](t *testing.T, keys []K) {
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
	t.Run("probe stats", func(t *testing.T) {
		testProbeStats(t, keys)
	})
}

/*
	func uniq[K comparable](keys []K) []K {
		s := make(Set[K], len(keys))
		for _, k := range keys {
			s[k] = struct{}
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
*/
func testSetPut[K comparable](t *testing.T, keys []K) {
	m := NewSet[K](uint32(len(keys)))
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
	m := NewSet[K](uint32(len(keys)))
	for _, key := range keys {
		m.Add(key)
	}
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
}

func testSetDelete[K comparable](t *testing.T, keys []K) {
	m := NewSet[K](uint32(len(keys)))
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
	m := NewSet[K](0)
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
	m := NewSet[K](uint32(len(keys)))
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
	m := NewSet[K](n / 10)
	for _, key := range keys {
		m.Add(key)
	}
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
}

func testSwissSetCapacity[K comparable](t *testing.T, gen func(n int) []K) {
	// Capacity() behavior depends on |groupSize|
	// which varies by processor architecture.
	caps := []uint32{
		1 * maxAvgGroupLoad,
		2 * maxAvgGroupLoad,
		3 * maxAvgGroupLoad,
		4 * maxAvgGroupLoad,
		5 * maxAvgGroupLoad,
		10 * maxAvgGroupLoad,
		25 * maxAvgGroupLoad,
		50 * maxAvgGroupLoad,
		100 * maxAvgGroupLoad,
	}
	for _, c := range caps {
		m := NewSet[K](c)
		assert.Equal(t, int(c), m.Capacity())
		keys := gen(rand.Intn(int(c)))
		for _, k := range keys {
			m.Add(k)
		}
		assert.Equal(t, int(c)-len(keys), m.Capacity())
		assert.Equal(t, int(c), m.Count()+m.Capacity())
	}
}

func testProbeStatsSet[K comparable](t *testing.T, keys []K) {
	runTest := func(load float32) {
		n := uint32(len(keys))
		sz, k := loadFactorSample(n, load)
		m := NewSet[K](sz)
		for _, key := range keys[:k] {
			m.Add(key)
		}
		// todo: assert stat invariants?
		stats := getProbeStatsSet(t, m, keys)
		t.Log(fmtProbeStats(stats))
	}
	t.Run("load_factor=0.5", func(t *testing.T) {
		runTest(0.5)
	})
	t.Run("load_factor=0.75", func(t *testing.T) {
		runTest(0.75)
	})
	t.Run("load_factor=max", func(t *testing.T) {
		runTest(maxLoadFactor)
	})
}

/*
// calculates the sample size and Set size necessary to
// create a load factor of |load| given |n| data points

	func loadFactorSample(n uint32, targetLoad float32) (SetSz, sampleSz uint32) {
		if targetLoad > maxLoadFactor {
			targetLoad = maxLoadFactor
		}
		// tables are assumed to be power of two
		sampleSz = uint32(float32(n) * targetLoad)
		SetSz = uint32(float32(n) * maxLoadFactor)
		return
	}

	type probeStats struct {
		groups     uint32
		loadFactor float32
		presentCnt uint32
		presentMin uint32
		presentMax uint32
		presentAvg float32
		absentCnt  uint32
		absentMin  uint32
		absentMax  uint32
		absentAvg  float32
	}

	func fmtProbeStats(s probeStats) string {
		g := fmt.Sprintf("groups=%d load=%f\n", s.groups, s.loadFactor)
		p := fmt.Sprintf("present(n=%d): min=%d max=%d avg=%f\n",
			s.presentCnt, s.presentMin, s.presentMax, s.presentAvg)
		a := fmt.Sprintf("absent(n=%d):  min=%d max=%d avg=%f\n",
			s.absentCnt, s.absentMin, s.absentMax, s.absentAvg)
		return g + p + a
	}
*/
func getProbeLengthSet[K comparable](t *testing.T, m *Set[K], key K) (length uint32, ok bool) {
	var end uint64
	hi, _ := splitHash(m.hashFunction.Hash(key))
	start := uint64(probeStart(hi, len(m.group)))
	end, _, ok = m.find(key)
	if end < start { // wrapped
		end += uint64(len(m.group))
	}
	length = uint32((end - start) + 1)
	require.True(t, length > 0)
	return
}

func getProbeStatsSet[K comparable](t *testing.T, m *Set[K], keys []K) (stats probeStats) {
	stats.groups = uint32(len(m.group))
	stats.loadFactor = m.loadFactor()
	var presentSum, absentSum float32
	stats.presentMin = math.MaxInt32
	stats.absentMin = math.MaxInt32
	for _, key := range keys {
		l, ok := getProbeLengthSet(t, m, key)
		if ok {
			stats.presentCnt++
			presentSum += float32(l)
			if stats.presentMin > l {
				stats.presentMin = l
			}
			if stats.presentMax < l {
				stats.presentMax = l
			}
		} else {
			stats.absentCnt++
			absentSum += float32(l)
			if stats.absentMin > l {
				stats.absentMin = l
			}
			if stats.absentMax < l {
				stats.absentMax = l
			}
		}
	}
	if stats.presentCnt == 0 {
		stats.presentMin = 0
	} else {
		stats.presentAvg = presentSum / float32(stats.presentCnt)
	}
	if stats.absentCnt == 0 {
		stats.absentMin = 0
	} else {
		stats.absentAvg = absentSum / float32(stats.absentCnt)
	}
	return
}

func TestNumGroupsSet(t *testing.T) {
	assert.Equal(t, expected(0), numGroups(0))
	assert.Equal(t, expected(1), numGroups(1))
	// max load factor 0.875
	assert.Equal(t, expected(14), numGroups(14))
	assert.Equal(t, expected(15), numGroups(15))
	assert.Equal(t, expected(28), numGroups(28))
	assert.Equal(t, expected(29), numGroups(29))
	assert.Equal(t, expected(56), numGroups(56))
	assert.Equal(t, expected(57), numGroups(57))
}

/*
func expected(x int) (groups uint32) {
	groups = uint32(math.Ceil(float64(x) / float64(maxAvgGroupLoad)))
	if groups == 0 {
		groups = 1
	}
	return
}
*/
