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
	"regexp"
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
	m := NewSet3WithSize[K](uint32(len(keys)))
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
	m := NewSet3WithSize[K](uint32(len(keys)))
	for _, key := range keys {
		m.Add(key)
	}
	for _, key := range keys {
		ok := m.Contains(key)
		assert.True(t, ok)
	}
}

func testSetDelete[K comparable](t *testing.T, keys []K) {
	m := NewSet3WithSize[K](uint32(len(keys)))
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
	m := NewSet3WithSize[K](0)
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
	m := NewSet3WithSize[K](uint32(len(keys)))
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
	m := NewSet3WithSize[K](n / 10)
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
	set1 := NewSet3WithSize[int](10)
	if !set1.Equals(set1) {
		t.Errorf("Test case 1: Both sets are the same instance: Expected true, got false")
	}

	set2 := NewSet3WithSize[int](10)
	if !set1.Equals(set2) {
		t.Errorf("Test case 2: Both sets are empty but different instances: Expected true, got false")
	}

	set3 := NewSet3WithSize[int](10)
	set3.Add(1)
	if set1.Equals(set3) {
		t.Errorf("Test case 3: Sets with different countsExpected false, got true")
	}

	set4 := NewSet3WithSize[int](10)
	set4.Add(1)
	if !set3.Equals(set4) {
		t.Errorf("Test case 4: Sets with same elements: Expected true, got false")
	}

	set5 := NewSet3WithSize[int](20)
	set5.Add(1)
	if !set3.Equals(set4) {
		t.Errorf("Test case 5: Sets with same elements but different capacities: Expected true, got false")
	}

	set6 := NewSet3WithSize[int](10)
	set6.Add(2)
	if set3.Equals(set6) {
		t.Errorf("Test case 6: Sets with different elements: Expected false, got true")
	}
}

func TestAsSet3(t *testing.T) {
	empty := NewSet3[int]()
	s1 := AsSet3([]int{})
	eq := empty.Equals(s1)
	assert.Equal(t, eq, true)
	s1.Add(1)
	s2 := AsSet3([]int{1})
	eq = s1.Equals(s2)
	assert.Equal(t, eq, true)
	s1.Add(2)
	s3 := AsSet3([]int{2, 1})
	eq = s1.Equals(s3)
	assert.Equal(t, eq, true)
}

func TestSet3GroupString(t *testing.T) {
	tests := []struct {
		name  string
		group set3Group[int]
		want  string
	}{
		{
			name:  "Empty slots",
			group: set3Group[int]{ctrl: set3AllEmpty, slot: [set3groupSize]int{0, 0, 0, 0, 0, 0, 0, 0}},
			want:  "[__|__|__|__|__|__|__|__]->{0|0|0|0|0|0|0|0}",
		},
		{
			name:  "Deleted slots",
			group: set3Group[int]{ctrl: set3AllDeleted, slot: [set3groupSize]int{0, 0, 0, 0, 0, 0, 0, 0}},
			want:  "[XX|XX|XX|XX|XX|XX|XX|XX]->{0|0|0|0|0|0|0|0}",
		},
		{
			name:  "Mixed slots",
			group: set3Group[int]{ctrl: 0x71727374757680FE, slot: [set3groupSize]int{1, 2, 3, 4, 5, 6, 7, 8}},
			want:  "[XX|__|76|75|74|73|72|71]->{1|2|3|4|5|6|7|8}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.group.String(); got != tt.want {
				t.Errorf("set3Group.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet3String(t *testing.T) {
	tests := []struct {
		name string
		set  Set3[int]
		want string
	}{
		{
			name: "Empty set",
			set:  *NewSet3[int](),
			want: "^\\{\\}$",
		},
		{
			name: "Single element",
			set:  *AsSet3([]int{1}),
			want: "^\\{1\\}$",
		},
		{
			name: "Multiple elements",
			set:  *AsSet3([]int{1, 2, 3}),
			want: "^\\{[1-3],[1-3],[1-3]\\}$",
		},
		{
			name: "Multiple groups",
			set:  *AsSet3([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			want: "^\\{[1-9],[1-9],[1-9],[1-9],[1-9],[1-9],[1-9],[1-9],[1-9]\\}$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.set.String()
			pattern := tt.want
			re, err := regexp.Compile(pattern)
			if err != nil {
				t.Errorf("Error compiling regex %v: %v", pattern, err)
			}
			match := re.MatchString(got)
			if !match {
				t.Errorf("Set3.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet3Clone(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	set2 := set1.Clone()
	assert.False(t, set1 == set2, "set2 shall not be identical with set1")
	assert.True(t, set1.Equals(set2), "set2 shall be equal to set1")
	set1.Add(4)
	set2.Add(5)
	assert.False(t, set1.Equals(set2), "set2 shall not be equal to set1 anymore")
}

func TestSet3ContainsAll(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	set2 := AsSet3([]int{1, 2, 3})
	assert.True(t, set1.ContainsAll(set2), "set2 shall be a subset of set1")
	set2.Remove(3)
	assert.True(t, set1.ContainsAll(set2), "set2 shall be a subset of set1")
	set2.Remove(2)
	assert.True(t, set1.ContainsAll(set2), "set2 shall be a subset of set1")
	set2.Remove(1)
	assert.True(t, set1.ContainsAll(set2), "set2 shall be a subset of set1")
	empty := NewSet3[int]()
	assert.True(t, empty.ContainsAll(set2), "set2 shall be a subset of an empty set")
	set3 := set1.Clone()
	set3.Add(4)
	assert.False(t, set1.ContainsAll(set3), "set3 shall not be a subset of set1")
	assert.False(t, empty.ContainsAll(set3), "set3 shall be a subset of an empty set")
	assert.True(t, set1.ContainsAll(set1), "set1 shall be a subset of set1")
	set4 := set1.Clone()
	assert.True(t, set1.ContainsAll(set4), "set4 shall be a subset of set1")
	set5 := AsSet3([]int{3, 4, 5})
	assert.False(t, set1.ContainsAll(set5), "set5 shall not be a subset of set1")
}

func TestSet3ContainsAllFrom(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	assert.True(t, set1.ContainsAllFrom([]int{1, 2}), "[]int{1,2} shall be a subset of set1")
	assert.True(t, set1.ContainsAllFrom([]int{2}), "[]int{2} shall be a subset of set1")
	assert.True(t, set1.ContainsAllFrom([]int{}), "[]int{} shall be a subset of set1")
	empty := NewSet3[int]()
	assert.True(t, empty.ContainsAllFrom([]int{}), "[]int{} shall be a subset of empty")
	assert.True(t, set1.ContainsAllFrom([]int{2, 2, 2, 2, 2, 2}), "[]int{2,2,2,2,2,2} shall be a subset of set1")
	assert.False(t, set1.ContainsAllFrom([]int{2, 4}), "[]int{2,4} shall be a subset of set1")
	assert.False(t, set1.ContainsAllFrom([]int{4}), "[]int{4} shall be a subset of set1")
}

func TestSet3Union(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	set2 := AsSet3([]int{4, 5, 6})
	set3 := AsSet3([]int{1, 2, 3, 4, 5, 6})
	u1 := set1.Union(set2)
	assert.False(t, set1.Equals(u1), "set1 shall not be altered by union")
	assert.True(t, u1.Equals(set3), "u1 and set3 shall be equal")
	u2 := set2.Union(set1)
	assert.True(t, u2.Equals(set3), "u2 and set3 shall be equal")
	empty := NewSet3[int]()
	u3 := set1.Union(empty)
	assert.True(t, set1.Equals(u3), "set1 shall be equal to u3")
	u4 := empty.Union(set1)
	assert.True(t, set1.Equals(u4), "set1 shall be equal to u4")
	set4 := AsSet3([]int{2, 3, 4, 5, 6})
	set5 := AsSet3([]int{1, 2, 3, 4, 5, 6})
	u5 := set1.Union(set4)
	assert.True(t, u5.Equals(set5), "u5 and set5 shall be equal")
}

func TestSet3AddAll(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	set2 := AsSet3([]int{4, 5, 6})
	set3 := AsSet3([]int{1, 2, 3, 4, 5, 6})
	set1.AddAll(set2)
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	empty := NewSet3[int]()
	set1.AddAll(empty)
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	set1.AddAll(set2)
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
}

func TestSet3AddAllFrom(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	set3 := AsSet3([]int{1, 2, 3, 4, 5, 6})
	set1.AddAllFrom([]int{4, 5, 6})
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	set1.AddAllFrom([]int{})
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	set1.AddAllFrom([]int{4, 5, 6})
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
}

func TestSet3Intersection(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3, 4})
	set2 := AsSet3([]int{3, 4, 5, 6})
	set3 := AsSet3([]int{3, 4})
	i1 := set1.Intersection(set2)
	assert.False(t, set1.Equals(i1), "set1 shall not be altered by intersection")
	assert.True(t, i1.Equals(set3), "i1 and set3 shall be equal")

	empty := NewSet3[int]()
	i2 := set1.Intersection(empty)
	assert.True(t, empty.Equals(i2), "empty shall be equal to i2")
	i3 := empty.Intersection(set1)
	assert.True(t, empty.Equals(i3), "empty shall be equal to i3")

	set4 := AsSet3([]int{1, 2, 3})
	set5 := AsSet3([]int{4, 5, 6})
	i4 := set4.Intersection(set5)
	assert.True(t, empty.Equals(i4), "empty shall be equal to i4")
}

func TestSet3RemoveAll(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3})
	set2 := AsSet3([]int{3, 4, 5, 6})
	set3 := AsSet3([]int{1, 2})
	set1.RemoveAll(set2)
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	empty := NewSet3[int]()
	set1.RemoveAll(empty)
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	set1.RemoveAll(set3)
	assert.True(t, set1.Equals(empty), "set1 and empty shall be equal")
}

func TestSet3RemoveAllFrom(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3, 4, 5, 6})
	set3 := AsSet3([]int{1, 2, 3, 4})
	set1.RemoveAllFrom([]int{5, 6, 7, 8})
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	set1.RemoveAllFrom([]int{})
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
	set1.RemoveAllFrom([]int{9, 10})
	assert.True(t, set1.Equals(set3), "set1 and set3 shall be equal")
}

func TestSet3Difference(t *testing.T) {
	set1 := AsSet3([]int{1, 2, 3, 4})
	set2 := AsSet3([]int{3, 4, 5, 6})
	set3 := AsSet3([]int{1, 2})
	i1 := set1.Difference(set2)
	assert.False(t, set1.Equals(i1), "set1 shall not be altered by intersection")
	assert.True(t, i1.Equals(set3), "i1 and set3 shall be equal")

	empty := NewSet3[int]()
	i2 := set1.Difference(empty)
	assert.True(t, set1.Equals(i2), "set1 shall be equal to i2")
}

func TestRehash(t *testing.T) {
	data := genUint32Data(53)
	set := AsSet3(data)
	assert.True(t, len(set.group) == 9, "set shall contain 9 groups")
	set.RehashTo(200)
	assert.True(t, len(set.group) == 31, "set shall contain 30 groups")
	for _, e := range data {
		assert.True(t, set.Contains(e), "set shall contain %v", e)
	}
	set.RehashTo(20)
	assert.True(t, len(set.group) == 31, "set shall contain 30 groups")
	for _, e := range data {
		assert.True(t, set.Contains(e), "set shall contain %v", e)
	}
	set.Rehash()
	assert.True(t, len(set.group) == 9, "set shall contain 9 groups")
	for _, e := range data {
		assert.True(t, set.Contains(e), "set shall contain %v", e)
	}
}

func TestToArray(t *testing.T) {
	set := AsSet3([]int{1, 2, 2, 3})
	ary := set.ToArray()
	assert.True(t, len(ary) == 3, "the array shall contain 3 elements")
	new_set := AsSet3(ary)
	assert.True(t, set.Equals(new_set), "both sets shall contain the same 3 elements")
	empty := NewSet3[int]()
	ary = empty.ToArray()
	assert.True(t, len(ary) == 0, "the array shall contain 0 elements")
}
