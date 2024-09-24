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

package set3

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

type testMapType map[uint64]struct{}

// see https://en.wikipedia.org/wiki/Xorshift#xorshift*
// This PRNG is deterministic and has a period of 2^64-1. This way we can ensure, we always get a new 'random' number, that is unknown to the set.
type prngState uint64
func (thisState *prngState) Uint64() uint64 {
	x := *thisState
    x ^= x >> 12
    x ^= x << 25
    x ^= x >> 27
	*thisState = x
	return uint64(x) * 0x2545F4914F6CDD1D
}

func (thisState *prngState) Uint32() uint32 {
	x := thisState.Uint64()
	x >>= 32 // the upper 32 bit have better 'randomness', see https://en.wikipedia.org/wiki/Xorshift#xorshift*
	return uint32(x)
}

type searchDataDriver struct {
	rng *prngState
	setValues []uint64
	hitRatio float64
}

func newSearchDataDriver(setSize int, targetHitRatio float32, seed int64) *searchDataDriver {
	s := prngState(seed)
	vals := uniqueRandomNumbers(setSize, &s)
	result := &searchDataDriver{
		rng: &s,
		setValues: shuffleArray(vals, rand.New(rand.NewSource(seed)), 3),
		hitRatio: float64(targetHitRatio),
	}
	return result
}

// tis function is designed in a way that both paths - table lookup and random number generation only are about equaly fast/slow
func (thisCfg *searchDataDriver) nextSearchValue() uint64 {
	x := uint64(float64(math.MaxUint64) * thisCfg.hitRatio)
	rndVal := thisCfg.rng.Uint64()
	upper32 := uint32(rndVal>>32)
	idx := upper32 % uint32(len(thisCfg.setValues))
	tableVal := thisCfg.setValues[idx]
	var result uint64
	if thisCfg.rng.Uint64() < x {
		// this shall be a hit
		result = rndVal ^ tableVal ^ rndVal // use both values to make both paths equally slow/fast
	} else {
		// this shall be a miss
		result = tableVal ^ rndVal ^ tableVal // use both values to make both paths equally slow/fast
	}
	return result
}

func uniqueRandomNumbers(setSize int, rng *prngState) []uint64 {
	tmpSet := EmptyWithCapacity[uint64](2 * uint32(setSize))
	for tmpSet.Count() < uint32(setSize) {
		tmpSet.Add(rng.Uint64())
	}
	result := tmpSet.ToArray()
	return result
}

func shuffleArray(input []uint64, rng *rand.Rand, rounds int) (output []uint64) {
	a := input
	b := make([]uint64, len(input))
	for i := 0; i<rounds; i++ {
		perm := rng.Perm(len(input))
		for j := 0; j<len(input); j++ {
			b[j] = a[perm[j]];
		}
		temp := a
		a = b
		b = temp
	}
	output = a
	return
}

var config = []struct {
	initSetSize int
	finalSetSize    int
	targetHitRatio  float32
	seed            int64
}{
	{initSetSize: 21, finalSetSize: 10, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 20, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 30, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 40, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 50, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
/*	{initSetSize: 21, finalSetSize: 6, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 7, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 8, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 9, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 10, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 11, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 12, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 13, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 14, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 15, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 16, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 17, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 18, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 19, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 20, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 21, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 22, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 23, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 24, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 25, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 26, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 27, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 28, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 29, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 30, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 31, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 32, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 33, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 34, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 35, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 36, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 37, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 38, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 39, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 40, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 41, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 42, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 43, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 44, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 45, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 46, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 47, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 48, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 49, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 50, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 51, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 52, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 53, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 54, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 55, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 56, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 57, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 58, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 59, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 60, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 61, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 62, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 63, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 64, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 65, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 66, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 67, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 68, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 69, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 70, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 71, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 72, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 73, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 74, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 75, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 76, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 77, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 78, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 79, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 80, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 81, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 82, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 83, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 84, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 85, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 86, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 87, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 88, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 89, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 90, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 91, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 92, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 93, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 94, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 95, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 96, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 97, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 98, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 99, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 100, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 101, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 102, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 103, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 104, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 105, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 106, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 107, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 108, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 109, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 110, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 111, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 112, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 113, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 114, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 115, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 116, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 117, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 118, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 119, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 120, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 121, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 122, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 123, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 124, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 125, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 126, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 127, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 128, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 129, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 130, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 131, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 132, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 133, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 134, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 135, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 136, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 137, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 138, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 139, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 140, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 141, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 142, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 143, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 144, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 145, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 146, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 147, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 148, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 149, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 150, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 151, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 152, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 153, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 154, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 155, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 156, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 157, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 158, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 159, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 160, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 161, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 162, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 163, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 164, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 165, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 166, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 167, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 168, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 169, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 170, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 171, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 172, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 173, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 174, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 175, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 176, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 177, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 178, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 179, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 180, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 181, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 182, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 183, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 184, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 185, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 186, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 187, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 188, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 189, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 190, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 191, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 192, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 193, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 194, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 195, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 196, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 197, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 198, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 199, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 200, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 201, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 202, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 203, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 204, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 205, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 206, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 207, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 208, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 209, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 210, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 211, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 212, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 213, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 214, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 215, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 216, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 217, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 218, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 219, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 220, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 221, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 222, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 223, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 224, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 225, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 226, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 227, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 228, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 229, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 230, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 231, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 232, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 233, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 234, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 235, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 236, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 237, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 238, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 239, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 240, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 241, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 242, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 243, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 244, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 245, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 246, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 247, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 248, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 249, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 250, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 251, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 252, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 253, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 254, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 255, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 256, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 257, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 258, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 259, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 260, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 261, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 262, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 263, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 264, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 265, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 266, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 267, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 268, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 269, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 270, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 271, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 272, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 273, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 274, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 275, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 276, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 277, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 278, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 279, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 280, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 281, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 282, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 283, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 284, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 285, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 286, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 287, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 288, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 289, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 290, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 291, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 292, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 293, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 294, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 295, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 296, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 297, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 298, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 299, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{initSetSize: 21, finalSetSize: 300, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF},
*/}

func BenchmarkSearchDataDriver(b *testing.B) {
	b.Skip("unskip to benchmark nextSearchValue paths")
	for _, cfg := range config {
		sdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed)
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		var x uint64 = 0
		sdd.hitRatio = 0.0
		b.Run(fmt.Sprintf("setSize(%d);hit(%f)", len(sdd.setValues), sdd.hitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x ^= sdd.nextSearchValue()
			}
		})
		sdd.hitRatio = 1.0
		b.Run(fmt.Sprintf("setSize(%d);hit(%f)", len(sdd.setValues), sdd.hitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x ^= sdd.nextSearchValue()
			}
		})
		// println(x)
	}
}

func BenchmarkSet3Fill(b *testing.B) {
	for _, cfg := range config {
		sdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed)
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		b.Run(fmt.Sprintf("init(%d);final(%d);hit(%f)", cfg.initSetSize, cfg.finalSetSize, cfg.targetHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := EmptyWithCapacity[uint64](uint32(cfg.initSetSize))
				for j := 0; j < len(sdd.setValues); j++ {
					resultSet.Add(sdd.setValues[j])
				}
			}
		})
	}
}

func BenchmarkNativeMapFill(b *testing.B) {
	for _, cfg := range config {
		sdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed)
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		b.Run(fmt.Sprintf("init(%d);final(%d);hit(%f)", cfg.initSetSize, cfg.finalSetSize, cfg.targetHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := make(testMapType, cfg.initSetSize)
				for j := 0; j < len(sdd.setValues); j++ {
					resultSet[sdd.setValues[j]] = struct{}{}
				}
			}
		})
	}
}

func BenchmarkSet3Find(b *testing.B) {
	for _, cfg := range config {
		for round := 0; round < 5 ; round ++ {
			sdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed+int64(round))
			resultSet := FromArray(sdd.setValues)
			// Force garbage collection
			runtime.GC()
			// Give the garbage collector some time to complete
			time.Sleep(1 * time.Second)
			var hit uint64 = 0
			var total uint64 = 0
			b.Run(fmt.Sprintf("init(%d);final(%d);hit(%f)", len(sdd.setValues), cfg.finalSetSize, cfg.targetHitRatio), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					search := sdd.nextSearchValue()
					if resultSet.Contains(search) {
						hit++
					}
					total++
				}
			})
			b.Logf("Actual hit ratio: %.3f", float32(hit)/float32(total))
		}
	}
}

func BenchmarkNativeMapFind(b *testing.B) {
	for _, cfg := range config {
		sdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed)
		resultSet := make(testMapType, len(sdd.setValues))
		for j := 0; j < len(sdd.setValues); j++ {
			resultSet[sdd.setValues[j]] = struct{}{}
		}
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		var hit uint64 = 0
		var total uint64 = 0
		b.Run(fmt.Sprintf("init(%d);final(%d);hit(%f)", len(sdd.setValues), cfg.finalSetSize, cfg.targetHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				search := sdd.nextSearchValue()
				_, b := resultSet[search]
				if b {
					hit++
				}
				total++
			}
		})
		b.Logf("Actual hit ratio: %.3f", float32(hit)/float32(total))
	}
}
