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
	/*
		t.Run("string capacity", func(t *testing.T) {
			testSet3Capacity(t, func(n int) []string {
				return genStringData(16, n)
			})
		})
		t.Run("uint32 capacity", func(t *testing.T) {
			testSet3Capacity(t, genUint32Data)
		})
	*/
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
	assert.Equal(t, uint32(0), m.Count())
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, uint32(len(keys)), m.Count())
	// overwrite
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, uint32(len(keys)), m.Count())
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
	assert.Equal(t, uint32(0), m.Count())
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, uint32(len(keys)), m.Count())
	for _, key := range keys {
		m.Remove(key)
		ok := m.Contains(key)
		assert.False(t, ok)
	}
	assert.Equal(t, uint32(0), m.Count())
	// put keys back after deleting them
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, uint32(len(keys)), m.Count())
}

func testSetClear[K comparable](t *testing.T, keys []K) {
	m := NewSet3[K](0)
	assert.Equal(t, uint32(0), m.Count())
	for _, key := range keys {
		m.Add(key)
	}
	assert.Equal(t, uint32(len(keys)), m.Count())
	m.Clear()
	assert.Equal(t, uint32(0), m.Count())
	for _, key := range keys {
		ok := m.Contains(key)
		assert.False(t, ok)
	}
	var calls int
	for _ = range m.ImmutableRange() {
		calls++
	}
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
	for _, k := range keys {
		visited[k] = 0
	}
	for e := range m.ImmutableRange() {
		visited[e]++
	}
	for _, c := range visited {
		assert.Equal(t, c, uint(1))
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

/*
	func testSet3Capacity[K comparable](t *testing.T, gen func(n int) []K) {
		// Capacity() behavior depends on |groupSize|
		// which varies by processor architecture.
		caps := []uint32{
			uint32(1.0 * set3maxAvgGroupLoad),
			uint32(2.0 * set3maxAvgGroupLoad),
			uint32(3.0 * set3maxAvgGroupLoad),
			uint32(4.0 * set3maxAvgGroupLoad),
			uint32(5.0 * set3maxAvgGroupLoad),
			uint32(10.0 * set3maxAvgGroupLoad),
			uint32(25.0 * set3maxAvgGroupLoad),
			uint32(50.0 * set3maxAvgGroupLoad),
			uint32(100.0 * set3maxAvgGroupLoad),
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
*/
func TestMutableRange(t *testing.T) {
	tests := []struct {
		name     string
		set      Set3[int]
		expected []int
	}{
		{
			name: "Empty set",
			set: Set3[int]{group: []set3Group[int]{
				{
					ctrl: set3AllEmpty,
					slot: [8]int{0, 0, 0, 0, 0, 0, 0, 0},
				},
			}},
			expected: []int{},
		},
		{
			name: "Single group with elements",
			set: Set3[int]{group: []set3Group[int]{
				{
					ctrl: 0x8001800180018001,
					slot: [8]int{1, 0, 3, 0, 5, 0, 7, 0},
				},
			}},
			expected: []int{1, 3, 5, 7},
		},
		{
			name: "Multiple groups with elements",
			set: Set3[int]{group: []set3Group[int]{
				{
					ctrl: 0x8001800180018001,
					slot: [8]int{1, 2, 3, 4, 5, 6, 7, 8},
				},
				{
					ctrl: 0x0180018001800180,
					slot: [8]int{9, 10, 11, 12, 13, 14, 15, 16},
				},
				{
					ctrl: set3AllEmpty,
					slot: [8]int{9, 10, 11, 12, 13, 14, 15, 16},
				},
			}},
			expected: []int{1, 3, 5, 7, 10, 12, 14, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []int
			for e := range tt.set.ImmutableRange() {
				result = append(result, e)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %v at index %d, got %v", tt.expected[i], i, v)
				}
			}
		})
	}
}

func TestImmutableRange(t *testing.T) {
	tests := []struct {
		name     string
		set      Set3[int]
		expected []int
	}{
		{
			name: "Empty set",
			set: Set3[int]{group: []set3Group[int]{
				{
					ctrl: set3AllEmpty,
					slot: [8]int{0, 0, 0, 0, 0, 0, 0, 0},
				},
			}},
			expected: []int{},
		},
		{
			name: "Single group with elements",
			set: Set3[int]{group: []set3Group[int]{
				{
					ctrl: 0x8001800180018001,
					slot: [8]int{1, 0, 3, 0, 5, 0, 7, 0},
				},
			}},
			expected: []int{1, 3, 5, 7},
		},
		{
			name: "Multiple groups with elements",
			set: Set3[int]{group: []set3Group[int]{
				{
					ctrl: 0x8001800180018001,
					slot: [8]int{1, 2, 3, 4, 5, 6, 7, 8},
				},
				{
					ctrl: 0x0180018001800180,
					slot: [8]int{9, 10, 11, 12, 13, 14, 15, 16},
				},
				{
					ctrl: set3AllEmpty,
					slot: [8]int{17, 18, 19, 20, 21, 22, 23, 24},
				},
			}},
			expected: []int{1, 3, 5, 7, 10, 12, 14, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []int
			i := 0
			for e := range tt.set.ImmutableRange() {
				tt.set.group[0].ctrl = set3AllDeleted
				tt.set.group[0].slot[i] = i * 20
				result = append(result, e)
				i++
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %v at index %d, got %v", tt.expected[i], i, v)
				}
			}
		})
	}
}

func TestEquals(t *testing.T) {
	set1 := NewSet3[int](10)
	if !set1.Equals(set1) {
		t.Errorf("Test case 1: Both sets are the same instance: Expected true, got false")
	}

	set2 := NewSet3[int](10)
	if !set1.Equals(set2) {
		t.Errorf("Test case 2: Both sets are empty but different instances: Expected true, got false")
	}

	set3 := NewSet3[int](10)
	set3.Add(1)
	if set1.Equals(set3) {
		t.Errorf("Test case 3: Sets with different countsExpected false, got true")
	}

	set4 := NewSet3[int](10)
	set4.Add(1)
	if !set3.Equals(set4) {
		t.Errorf("Test case 4: Sets with same elements: Expected true, got false")
	}

	set5 := NewSet3[int](20)
	set5.Add(1)
	if !set3.Equals(set4) {
		t.Errorf("Test case 5: Sets with same elements but different capacities: Expected true, got false")
	}

	set6 := NewSet3[int](10)
	set6.Add(2)
	if set3.Equals(set6) {
		t.Errorf("Test case 6: Sets with different elements: Expected false, got true")
	}
}
