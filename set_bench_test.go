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
	"math/bits"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func BenchmarkStringSets(b *testing.B) {
	const keySz = 8
	sizes := []int{16, 128, 1024, 8192, 131072}
	for _, n := range sizes {
		b.Run("n="+strconv.Itoa(n), func(b *testing.B) {
			b.Run("runtime Set", func(b *testing.B) {
				benchmarkRuntimeSet(b, genStringData(keySz, n))
			})
			b.Run("swiss.Set", func(b *testing.B) {
				benchmarkSwissSet(b, genStringData(keySz, n))
			})
		})
	}
}

func BenchmarkInt64Sets(b *testing.B) {
	sizes := []int{16, 128, 1024, 8192, 131072}
	for _, n := range sizes {
		b.Run("n="+strconv.Itoa(n), func(b *testing.B) {
			b.Run("runtime Set", func(b *testing.B) {
				benchmarkRuntimeSet(b, generateInt64Data(n))
			})
			b.Run("swiss.Set", func(b *testing.B) {
				benchmarkSwissSet(b, generateInt64Data(n))
			})
		})
	}
}

func TestMemoryFootprintSet(t *testing.T) {
	//t.Skip("unskip for memory footprint stats")
	var samples []float64
	for n := 10; n <= 60_000; n += 10 {
		b1 := testing.Benchmark(func(b *testing.B) {
			// max load factor 7/8
			m := NewSet[int](uint32(n))
			require.NotNil(b, m)
		})
		b2 := testing.Benchmark(func(b *testing.B) {
			// max load factor 6.5/8
			m := make(map[int]struct{}, n)
			require.NotNil(b, m)
		})
		x := float64(b1.MemBytes) / float64(b2.MemBytes)
		samples = append(samples, x)
	}
	t.Logf("mean size ratio: %.3f", mean(samples))
}

func benchmarkRuntimeSet[K comparable](b *testing.B, keys []K) {
	n := uint32(len(keys))
	mod := n - 1 // power of 2 fast modulus
	require.Equal(b, 1, bits.OnesCount32(n))
	m := make(map[K]K, n)
	for _, k := range keys {
		m[k] = k
	}
	b.ResetTimer()
	var ok bool
	for i := 0; i < b.N; i++ {
		_, ok = m[keys[uint32(i)&mod]]
	}
	//	assert.True(b, ok)
	for i := b.N; i < b.N*2; i++ {
		_, ok = m[keys[uint32(i-b.N)&mod]]
	}
	assert.True(b, ok)
	b.ReportAllocs()
}

func benchmarkSwissSet[K comparable](b *testing.B, keys []K) {
	n := uint32(len(keys))
	mod := n - 1 // power of 2 fast modulus
	require.Equal(b, 1, bits.OnesCount32(n))
	m := NewSet[K](n)
	for _, k := range keys {
		m.Put(k)
	}
	b.ResetTimer()
	var ok bool
	for i := 0; i < b.N; i++ {
		ok = m.Has(keys[uint32(i)&mod])
	}
	//	assert.True(b, ok)
	for i := b.N; i < b.N*2; i++ {
		ok = m.Has(keys[uint32(i-b.N)&mod])
	}
	assert.True(b, ok)
	b.ReportAllocs()
}
