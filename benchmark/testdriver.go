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

package benchmark

import (
	"math"
)

// see https://en.wikipedia.org/wiki/Xorshift#xorshift*
// This PRNG is deterministic and has a period of 2^64-1. This way we can ensure, we always get a new 'random' number, that is unknown to the set.
type PRNG struct {
	State uint64
	Round uint64 // for debugging purposes
}

func (thisState *PRNG) Uint64() uint64 {
	x := thisState.State
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	thisState.State = x
	thisState.Round++
	return x * 0x2545F4914F6CDD1D
}

/*
func (thisState *prngState) Uint32() uint32 {
	x := thisState.Uint64()
	x >>= 32 // the upper 32 bit have better 'randomness', see https://en.wikipedia.org/wiki/Xorshift#xorshift*
	return uint32(x)
}
*/

type SearchDataDriver struct {
	rng       *PRNG
	SetValues []uint64
	hitRatio  float64
}

func NewSearchDataDriver(setSize int, targetHitRatio float64, seed uint64) *SearchDataDriver {
	s := PRNG{State: seed}
	vals := uniqueRandomNumbers(setSize, &s)
	result := &SearchDataDriver{
		rng: &s,
		// setValues: shuffleArray(vals, &s, 3),
		SetValues: vals,
		hitRatio:  float64(targetHitRatio),
	}
	return result
}

// this function is designed in a way that both paths - table lookup and random number generation only - are about equaly fast/slow.
// the current implementation differs in only 1-2% execution speed for the two paths.
func (thisCfg *SearchDataDriver) NextSearchValue() uint64 {
	x := uint64(float64(math.MaxUint64) * thisCfg.hitRatio)
	rndVal := thisCfg.rng.Uint64()
	upper32 := uint32(rndVal >> 32)
	idx := upper32 % uint32(len(thisCfg.SetValues))
	tableVal := thisCfg.SetValues[idx]
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

func uniqueRandomNumbers(setSize int, rng *PRNG) []uint64 {
	result := make([]uint64, setSize)
	for i := 0; i < setSize; i++ {
		result[i] = rng.Uint64()
	}
	return result
}

/*
func shuffleArray(input []uint64, rng *prngState, rounds int) []uint64 {
	a := input // copy array
	for r := 0; r < rounds; r++ {
		for i := len(a) - 1; i > 0; i-- {
			j := rng.Uint32() % uint32(i+1)
			temp := a[i]
			a[i] = a[j]
			a[j] = temp
		}
	}
	return a
}
*/
