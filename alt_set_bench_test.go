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
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// see https://en.wikipedia.org/wiki/Xorshift#xorshift*
// This PRNG is deterministic and has a period of 2^64-1. This way we can ensure, we always get a new 'random' number, that is unknown to the set.
type prngState struct {
	state uint64
	round uint64 // for debugging purposes
}

func (thisState *prngState) Uint64() uint64 {
	x := thisState.state
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	thisState.state = x
	thisState.round++
	return uint64(x) * 0x2545F4914F6CDD1D
}

func (thisState *prngState) Uint32() uint32 {
	x := thisState.Uint64()
	x >>= 32 // the upper 32 bit have better 'randomness', see https://en.wikipedia.org/wiki/Xorshift#xorshift*
	return uint32(x)
}

type searchDataDriver struct {
	rng       *prngState
	setValues []uint64
	hitRatio  float64
}

func newSearchDataDriver(setSize int, targetHitRatio float64, seed uint64) *searchDataDriver {
	s := prngState{state: seed}
	vals := uniqueRandomNumbers(setSize, &s)
	result := &searchDataDriver{
		rng: &s,
		// setValues: shuffleArray(vals, &s, 3),
		setValues: vals,
		hitRatio:  float64(targetHitRatio),
	}
	return result
}

// this function is designed in a way that both paths - table lookup and random number generation only - are about equaly fast/slow.
// the current implementation differs in only 1-2% execution speed for the two paths.
func (thisCfg *searchDataDriver) nextSearchValue() uint64 {
	x := uint64(float64(math.MaxUint64) * thisCfg.hitRatio)
	rndVal := thisCfg.rng.Uint64()
	upper32 := uint32(rndVal >> 32)
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
	// tmpSet := EmptyWithCapacity[uint64](2 * uint32(setSize))
	// for tmpSet.Count() < uint32(setSize) {
	//	tmpSet.Add(rng.Uint64())
	// }
	// result := tmpSet.ToArray()
	result := make([]uint64, setSize)
	for i := 0; i < setSize; i++ {
		result[i] = rng.Uint64()
	}
	return result
}

func shuffleArray(input []uint64, rng *prngState, rounds int) []uint64 {
	a := input // copy array
	for r := 0; r < rounds; r++ {
		for i := len(a) - 1; i > 0; i-- {
			j := rng.Uint32() % uint32(i+1)
			a[i], a[j] = a[j], a[i]
		}
	}
	return a
}

var config = []struct {
	initSetSize    int
	finalSetSize   int
	targetHitRatio float64
	seed           uint64
	itersPerRound  int
	rounds         int
}{
	{initSetSize: 21, finalSetSize: 1, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 2, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 3, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 4, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 5, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 6, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 7, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 8, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 9, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 10, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 11, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 12, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 13, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 14, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 15, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 16, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 17, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 18, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 19, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 20, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 21, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 22, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 23, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 24, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 25, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 26, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 27, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 28, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 29, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 30, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 31, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 32, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 33, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 34, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 35, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 36, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 37, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 38, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 39, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 40, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 41, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 42, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 43, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 44, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 45, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 46, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 47, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 48, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 49, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 50, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 51, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 52, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 53, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 54, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 55, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 56, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 57, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 58, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 59, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 60, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 61, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 62, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 63, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 64, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 65, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 66, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 67, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 68, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 69, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 70, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 71, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 72, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 73, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 74, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 75, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 76, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 77, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 78, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 79, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 80, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 81, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 82, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 83, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 84, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 85, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 86, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 87, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 88, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 89, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 90, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 91, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 92, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 93, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 94, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 95, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 96, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 97, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 98, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 99, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
	{initSetSize: 21, finalSetSize: 100, targetHitRatio: 0.3, seed: 0x1234567890ABCDEF, itersPerRound: 10_000_000, rounds: 11},
}

func TestPrngSeqLength(t *testing.T) {
	state := prngState{state: 0x1234567890ABCDEF}
	limit := uint32(30_000_000)
	set := EmptyWithCapacity[uint64](limit * 2)
	counter := uint32(0)
	for set.Count() < limit {
		set.Add(state.Uint64())
		counter++
	}
	assert.True(t, counter == limit, "sequence < limit")
}

func TestPrngDeterminism(t *testing.T) {
	state1 := prngState{state: 0x1234567890ABCDEF}
	state2 := prngState{state: 0x1234567890ABCDEF}
	limit := 30_000_000
	for i := 0; i < limit; i++ {
		v1 := state1.Uint64()
		v2 := state2.Uint64()
		assert.True(t, v1 == v2, "in sync: values not equal in round %d", i)
	}
	_ = state2.Uint64() // skip one value to get both prng out of sync
	for i := 0; i < limit; i++ {
		v1 := state1.Uint64()
		v2 := state2.Uint64()
		assert.False(t, v1 == v2, "out of sync: values equal in round %d", i)
	}
	_ = state1.Uint64() // get both prng back in sync
	for i := 0; i < limit; i++ {
		v1 := state1.Uint64()
		v2 := state2.Uint64()
		assert.True(t, v1 == v2, "back in sync: values not equal in round %d", i)
	}
}

func TestUniqueRandomNumbersDeterministic(t *testing.T) {
	s1 := prngState{state: 0x1234567890ABCDEF}
	s2 := prngState{state: 0x1234567890ABCDEF}

	urn1 := uniqueRandomNumbers(5000, &s1)
	urn2 := uniqueRandomNumbers(5000, &s2)
	assert.True(t, slicesEqual(urn1, urn2), "slices not equal")

}

func TestSearchDataDriver(t *testing.T) {
	setSize := 500_000
	targetHitRatio := 0.3
	seed := uint64(0x1234567890ABCDEF)

	sdd1 := newSearchDataDriver(setSize, targetHitRatio, seed)
	sdd2 := newSearchDataDriver(setSize, targetHitRatio, seed)
	assert.True(t, slicesEqual(sdd1.setValues, sdd2.setValues), "slices not equal")

	set := FromArray(sdd1.setValues)

	rounds := 50_000_000
	hits := 0
	for i := 0; i < rounds; i++ {
		v1 := sdd1.nextSearchValue()
		v2 := sdd2.nextSearchValue()
		assert.True(t, v1 == v2, "values not equal in round %d", i)
		if set.Contains(v1) {
			hits++
		}
	}
	actualHitRatio := float64(hits) / float64(rounds)
	lowerBound := targetHitRatio * 0.99
	upperBound := targetHitRatio * 1.01
	assert.True(t, actualHitRatio > lowerBound && actualHitRatio < upperBound, "actual hit ratio (%d) missed target hit ratio by more than 1 percent", actualHitRatio)
}

func slicesEqual(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

var sddConfig = []struct {
	setSize int
	seed    uint64
}{
	{setSize: 1, seed: 0x1234567890ABCDEF},
	{setSize: 10, seed: 0x1234567890ABCDEF},
	{setSize: 10_000, seed: 0x1234567890ABCDEF},
	{setSize: 10_000_000, seed: 0x1234567890ABCDEF},
}

func BenchmarkSearchDataDriver(b *testing.B) {
	// b.Skip("unskip to benchmark nextSearchValue paths")
	for _, cfg := range sddConfig {
		sdd := newSearchDataDriver(cfg.setSize, 0.0, cfg.seed)
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(2 * time.Second)
		var x uint64 = 0
		b.Run(fmt.Sprintf("setSize(%d);hit(%.1f)", len(sdd.setValues), sdd.hitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x ^= sdd.nextSearchValue()
			}
		})
		sdd.hitRatio = 1.0
		b.Run(fmt.Sprintf("setSize(%d);hit(%.1f)", len(sdd.setValues), sdd.hitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x ^= sdd.nextSearchValue()
			}
		})
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
				resultSet := emptyNativeWithCapacity[uint64](uint32(cfg.initSetSize))
				for j := 0; j < len(sdd.setValues); j++ {
					resultSet.add(sdd.setValues[j])
				}
			}
		})
	}
}

func BenchmarkSet3FindVariance(b *testing.B) {
	for _, cfg := range config {
		for seedUp := 0; seedUp < 10; seedUp++ {
			for round := 0; round < 10; round++ {
				sdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed+uint64(seedUp*51))
				resultSet := FromArray(sdd.setValues)
				// Force garbage collection
				runtime.GC()
				// Give the garbage collector some time to complete
				time.Sleep(1 * time.Second)
				var hit uint64 = 0
				var total uint64 = 0
				b.Run(fmt.Sprintf("init(%d);final(%d);hit(%f)-s(%d)", len(sdd.setValues), cfg.finalSetSize, cfg.targetHitRatio, seedUp), func(b *testing.B) {
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
}

func median(data []float64) float64 {
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)
	sort.Float64s(dataCopy)

	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		return (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		return dataCopy[l/2]
	}
}

func TestSet3Find(t *testing.T) {
	fmt.Printf("Impl;Operation;init;size;hit_target;total_iter;hit_rea;time/iter[ns]\n")
	for _, cfg := range config {
		timePerIter := make([]float64, cfg.rounds)
		var hit uint64 = 0
		var total uint64 = 0
		for i := 0; i < cfg.rounds; i++ {
			currentSdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed+uint64(i*53))
			currentSet := FromArray(currentSdd.setValues)
			testdata := make([]uint64, cfg.itersPerRound)
			for j := range cfg.itersPerRound {
				testdata[j] = currentSdd.nextSearchValue()
			}
			currentSdd = nil
			runtime.GC()
			startTime := time.Now().UnixNano()
			for j := 0; j < cfg.itersPerRound; j++ {
				// search := currentSdd.nextSearchValue()
				search := testdata[j]
				if currentSet.Contains(search) {
					hit++
				}
				total++
			}
			endTime := time.Now().UnixNano()
			timePerIter[i] = float64(endTime-startTime) / float64(cfg.itersPerRound)
		}
		hitRea := float32(hit) / float32(total)
		med := median(timePerIter)
		fmt.Printf("Set3;Contains;%d;%d;%.3f;%d;%.3f;%.3f\n", cfg.finalSetSize, cfg.finalSetSize, cfg.targetHitRatio, cfg.itersPerRound, hitRea, med)
	}
}

func TestNativeMapFind(t *testing.T) {
	fmt.Printf("Impl;Operation;init;size;hit_target;total_iter;hit_rea;time/iter[ns]\n")
	for _, cfg := range config {
		timePerIter := make([]float64, cfg.rounds)
		var hit uint64 = 0
		var total uint64 = 0
		for i := 0; i < cfg.rounds; i++ {
			currentSdd := newSearchDataDriver(cfg.finalSetSize, cfg.targetHitRatio, cfg.seed+uint64(i*53))
			currentSet := emptyNativeWithCapacity[uint64](uint32(len(currentSdd.setValues)))
			for j := 0; j < len(currentSdd.setValues); j++ {
				currentSet.add(currentSdd.setValues[j])
			}
			testdata := make([]uint64, cfg.itersPerRound)
			for j := range cfg.itersPerRound {
				testdata[j] = currentSdd.nextSearchValue()
			}
			currentSdd = nil
			runtime.GC()
			startTime := time.Now().UnixNano()
			for j := 0; j < cfg.itersPerRound; j++ {
				// search := currentSdd.nextSearchValue()
				search := testdata[j]
				if currentSet.contains(search) {
					hit++
				}
				total++
			}
			endTime := time.Now().UnixNano()
			timePerIter[i] = float64(endTime-startTime) / float64(cfg.itersPerRound)
		}
		hitRea := float32(hit) / float32(total)
		med := median(timePerIter)
		fmt.Printf("map;Contains;%d;%d;%.3f;%d;%.3f;%.3f\n", cfg.finalSetSize, cfg.finalSetSize, cfg.targetHitRatio, cfg.itersPerRound, hitRea, med)
	}
}
