// Copyright 2023 Dolthub, Inc.
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
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestMatchMetadata(t *testing.T) {
	var meta [groupSize]int8
	for i := range meta {
		meta[i] = int8(i)
	}
	t.Run("metaMatchH2", func(t *testing.T) {
		for _, x := range meta {
			mask := ctlrMatchH2(&meta, (uint64(x) & 0x0000_0000_0000_007f))
			assert.NotZero(t, mask)
			assert.Equal(t, int(x), nextMatch(&mask))
		}
	})
	t.Run("metaMatchEmpty", func(t *testing.T) {
		mask := ctlrMatchEmpty(&meta)
		assert.Equal(t, mask, uint64(0))
		for i := range meta {
			meta[i] = kEmpty
			mask = ctlrMatchEmpty(&meta)
			assert.NotZero(t, mask)
			assert.Equal(t, int(i), nextMatch(&mask))
			meta[i] = int8(i)
		}
	})
	t.Run("nextMatch", func(t *testing.T) {
		// test iterating multiple matches
		for j := range groupSize {
			meta[j] = kEmpty
		}
		mask := ctlrMatchEmpty(&meta)
		for i := range meta {
			assert.Equal(t, int(i), nextMatch(&mask))
		}
		for i := 0; i < len(meta); i += 2 {
			meta[i] = int8(42)
		}
		mask = ctlrMatchH2(&meta, (uint64(42) & 0x0000_0000_0000_007f))
		for i := 0; i < len(meta); i += 2 {
			assert.Equal(t, int(i), nextMatch(&mask))
		}
	})
}

func BenchmarkMatchMetadata(b *testing.B) {
	var meta [groupSize]int8
	for i := range meta {
		meta[i] = int8(i)
	}
	var mask uint64
	for i := 0; i < b.N; i++ {
		mask = ctlrMatchH2(&meta, (uint64(i) & 0x0000_0000_0000_007f))
	}
	b.Log(mask)
}

func TestNextPow2(t *testing.T) {
	assert.Equal(t, 0, int(nextPow2(0)))
	assert.Equal(t, 1, int(nextPow2(1)))
	assert.Equal(t, 2, int(nextPow2(2)))
	assert.Equal(t, 4, int(nextPow2(3)))
	assert.Equal(t, 8, int(nextPow2(7)))
	assert.Equal(t, 8, int(nextPow2(8)))
	assert.Equal(t, 16, int(nextPow2(9)))
}

func nextPow2(x uint32) uint32 {
	return 1 << (32 - bits.LeadingZeros32(x-1))
}

func TestConstants(t *testing.T) {
	c1, c2 := kEmpty, kDeleted
	assert.Equal(t, byte(0b1000_0000), byte(c1))
	assert.Equal(t, byte(0b1000_0000), reinterpretCast(int8(c1)))
	assert.Equal(t, byte(0b1111_1110), byte(c2))
	assert.Equal(t, byte(0b1111_1110), reinterpretCast(int8(c2)))
}

func reinterpretCast(i int8) byte {
	return *(*byte)(unsafe.Pointer(&i))
}
