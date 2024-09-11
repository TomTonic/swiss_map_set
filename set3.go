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
	"iter"
	"math/bits"

	"github.com/dolthub/maphash"
)

const (
	set3groupSize       = 8
	set3maxAvgGroupLoad = 6.5
	//set3maxAvgGroupLoad = float64(20) / float64(3) // 6.66667

	set3maxLoadFactor = float32(set3maxAvgGroupLoad) / float32(set3groupSize)

	set3loBits uint64 = 0x0101010101010101
	set3hiBits uint64 = 0x8080808080808080

	set3AllEmpty   uint64 = 0x8080808080808080
	set3AllDeleted uint64 = 0xFEFEFEFEFEFEFEFE
	set3Empty      uint64 = 0b0000_1000_0000
	set3Deleted    uint64 = 0b0000_1111_1110
	set3Sentinel   uint64 = 0b0000_1111_1111
)

func set3ctlrMatchH2(m uint64, h uint64) uint64 {
	// https://graphics.stanford.edu/~seander/bithacks.html##ValueInWord
	return set3hasZeroByte(m ^ (set3loBits * h))
}

func set3ctlrMatchEmpty(m uint64) uint64 {
	return set3hasZeroByte(m ^ set3hiBits)
}

func set3nextMatch(b *uint64) int {
	s := bits.TrailingZeros64(*b)
	*b &= ^(1 << s) // clear bit |s|
	return s >> 3   // div by 8
}

func set3hasZeroByte(x uint64) uint64 {
	return ((x - set3loBits) & ^(x)) & set3hiBits
}

type set3Group[K comparable] struct {
	ctrl uint64
	slot [set3groupSize]K
}

// Set3 is an open-addressing Set3
// based on Abseil's flat_hash_map.
type Set3[K comparable] struct {
	hashFunction maphash.Hasher[K]
	resident     uint32
	dead         uint32
	elementLimit uint32
	group        []set3Group[K]
}

// NewSet3 constructs a Set3.
func NewSet3[K comparable](sz uint32) (s *Set3[K]) {
	reqNrOfGroups := int((float64(sz) + set3maxAvgGroupLoad - 1) / set3maxAvgGroupLoad)
	if reqNrOfGroups == 0 {
		reqNrOfGroups = 1
	}
	s = &Set3[K]{
		hashFunction: maphash.NewHasher[K](),
		elementLimit: uint32(float64(reqNrOfGroups) * set3maxAvgGroupLoad),
		group:        make([]set3Group[K], reqNrOfGroups),
	}
	for i := range len(s.group) {
		s.group[i].ctrl = set3AllEmpty
	}
	return
}

func (set3 *Set3[K]) Clone() (s *Set3[K]) {
	s = &Set3[K]{
		hashFunction: set3.hashFunction,
		elementLimit: set3.elementLimit,
		resident:     set3.resident,
		dead:         set3.dead,
		group:        set3.group,
	}
	return
}

// Contains returns true if |element| is present in the |Set3|.
func (Set3 *Set3[K]) Contains(element K) bool {
	hash := Set3.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpCnt := uint64(len(Set3.group))
	grpIdx := H1 % grpCnt
	for {
		group := &Set3.group[grpIdx]
		ctrl := group.ctrl
		slot := &(group.slot)
		matches := set3ctlrMatchH2(ctrl, H2)
		//matches := simd.MatchCRTLhash(ctrl, H2)
		for matches != 0 {
			s := set3nextMatch(&matches)
			if element == slot[s] {
				return true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = set3ctlrMatchEmpty(ctrl)
		//matches = simd.MatchCRTLempty(ctrl)
		if matches != 0 {
			// there is an empty slot - the element, if it had been added, hat either
			// been found until now or it had been added in the next empty spot -
			// well, this is the next empty spot...
			return false
		}
		grpIdx += 1 // carousel through all groups
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

func (this *Set3[K]) Equals(that *Set3[K]) bool {
	if this == that {
		return true
	}
	if this.Count() != that.Count() {
		return false
	}
	for elem := range that.MutableRange() {
		if !this.Contains(elem) {
			return false
		}
	}
	return true
}

// Iterates over all elements in this |Set3|.
// Caution: If the |Set3| is changed during the iteration, the result may be arbitrary.
// If the |Set3| might get changed during the itration, you might prefer ImmutableRange()
func (s *Set3[K]) MutableRange() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, group := range s.group {
			ctrl := group.ctrl
			if ctrl&set3hiBits != set3hiBits { // not all empty or deleted
				slot := &(group.slot)
				for i := 0; i < set3groupSize; i++ {
					if elementAt(ctrl, i) {
						if !yield(slot[i]) {
							return
						}
					}
				}
			}
		}
	}
}

// Iterates over all elements in this |Set3|. Makes an internal copy of the stored elements first.
// To avoid this copy, choose MutableRange()
func (s *Set3[K]) ImmutableRange() iter.Seq[K] {
	return func(yield func(K) bool) {
		groups := s.group
		for _, group := range groups {
			ctrl := group.ctrl
			if ctrl&set3hiBits != set3hiBits { // not all empty or deleted
				slot := &(group.slot)
				for i := 0; i < set3groupSize; i++ {
					if elementAt(ctrl, i) {
						if !yield(slot[i]) {
							return
						}
					}
				}
			}
		}
	}
}

// Add attempts to insert |key| and |value|
func (Set3 *Set3[K]) Add(element K) {
	if Set3.resident >= Set3.elementLimit {
		Set3.rehash(Set3.nextSize())
	}
	hash := Set3.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpCnt := uint64(len(Set3.group))
	grpIdx := H1 % grpCnt
	for {
		group := &Set3.group[grpIdx]
		ctrl := group.ctrl
		slot := &(group.slot)

		matches := set3ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := set3nextMatch(&matches)
			if element == slot[s] {
				// found - already in Set3, just return
				return
			}
		}

		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = set3ctlrMatchEmpty(ctrl)

		if matches != 0 {
			// empty spot -> element can't be in Set3 (see Contains) -> insert
			s := set3nextMatch(&matches)
			group.ctrl = setCTRLat(ctrl, H2, s)
			slot[s] = element
			Set3.resident++
			return

		}
		grpIdx += 1 // carousel through all groups
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

func setCTRLat(ctrl, val uint64, pos int) uint64 {
	shift := pos << 3                // *8
	ctrl &= ^(uint64(0xFF) << shift) // clear byte
	ctrl |= val << shift             // set byte to given value
	return ctrl
}

func elementAt(ctrl uint64, pos int) bool {
	shift := pos << 3               // *8
	ctrl &= (uint64(0x80) << shift) // clear all other bits
	// if a bit is set, the according byte represented either set3Empty or set3Deleted
	// -> if ctlr is 0 now, the according position stores a value
	return ctrl == 0
}

// Remove attempts to remove |element|, returns true if the |element| was in the |Set3|
func (Set3 *Set3[K]) Remove(element K) bool {
	hash := Set3.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpCnt := uint64(len(Set3.group))
	grpIdx := H1 % grpCnt
	for {
		group := &Set3.group[grpIdx]
		ctrl := group.ctrl
		slot := &(group.slot)
		matches := set3ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := set3nextMatch(&matches)
			if element == slot[s] {
				// found - already in Set3, just return
				// optimization: if |m.ctrl[g]| contains any empty
				// metadata bytes, we can physically delete |element|
				// rather than placing a tombstone.
				// The observation is that any probes into group |g|
				// would already be terminated by the existing empty
				// slot, and therefore reclaiming slot |s| will not
				// cause premature termination of probes into |g|.
				if set3ctlrMatchEmpty(ctrl) != 0 {
					group.ctrl = setCTRLat(ctrl, set3Empty, s)
					Set3.resident--
				} else {
					group.ctrl = setCTRLat(ctrl, set3Deleted, s)
					Set3.dead++
				}
				var k K
				slot[s] = k
				return true
			}
		}

		// |element| is not in group |g|,
		// stop probing if we see an empty slot
		matches = set3ctlrMatchEmpty(ctrl)
		if matches != 0 {
			// |element| absent
			return false
		}
		grpIdx += 1 // linear probing
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

// Clear removes all elements from the Map.
func (Set3 *Set3[K]) Clear() {
	var k K
	for grpidx := range len(Set3.group) {
		d := &(Set3.group[grpidx])
		d.ctrl = set3AllEmpty
		for j := range set3groupSize {
			d.slot[j] = k
		}
	}
	Set3.resident, Set3.dead = 0, 0
}

// Count returns the number of elements in the Map.
func (Set3 *Set3[K]) Count() int {
	return int(Set3.resident - Set3.dead)
}

// Capacity returns the number of additional elements
// the can be added to the Map before resizing.
func (Set3 *Set3[K]) Capacity() int {
	return int(Set3.elementLimit - Set3.resident)
}

// find returns the location of |key| if present, or its insertion location if absent.
// for performance, find is manually inlined into public methods.
func (Set3 *Set3[K]) find(key K) (g uint64, s int, ok bool) {
	//g = probeStart2(hi, len(Set3.data))
	hash := Set3.hashFunction.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	g = H1 % uint64(len(Set3.group))
	for {
		group := &Set3.group[g]
		ctrl := group.ctrl
		slot := &(group.slot)
		matches := set3ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s = set3nextMatch(&matches)
			if key == slot[s] {
				return g, s, true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = set3ctlrMatchEmpty(ctrl)
		if matches != 0 {
			s = set3nextMatch(&matches)
			return g, s, false
		}
		g += 1 // linear probing
		if g >= uint64(len(Set3.group)) {
			g = 0
		}
	}
}

func (Set3 *Set3[K]) nextSize() (n uint32) {
	n = uint32(len(Set3.group)) * 2
	if Set3.dead >= (Set3.resident / 2) {
		n = uint32(len(Set3.group))
	}
	return
}

func (Set3 *Set3[K]) rehash(n uint32) {
	old_groups := Set3.group
	Set3.hashFunction = maphash.NewSeed(Set3.hashFunction)
	Set3.elementLimit = uint32(float64(n) * set3maxAvgGroupLoad)
	Set3.resident, Set3.dead = 0, 0
	Set3.group = make([]set3Group[K], n)
	for i := range len(Set3.group) {
		Set3.group[i].ctrl = set3AllEmpty
	}
	grpCnt := uint64(len(Set3.group))
	for _, old_grp := range old_groups {
		if old_grp.ctrl&set3hiBits != set3hiBits { // not all empty or deleted
			for s := range set3groupSize {
				if elementAt(old_grp.ctrl, s) {
					// inlined and reduced Add instead of Set3.Add(old_grp.slot[s])
					element := old_grp.slot[s]

					hash := Set3.hashFunction.Hash(element)
					H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
					H2 := (hash & 0x0000_0000_0000_007f)
					grpIdx := H1 % uint64(len(Set3.group))
					stillSearchingSpace := true
					for stillSearchingSpace {
						group := &Set3.group[grpIdx]
						ctrl := group.ctrl
						slot := &(group.slot)

						// optimization: we know it cannot exist in the Set3 already so skip
						// searching for the hashcode and start searching for an empty slot
						// immediately
						matches := set3ctlrMatchEmpty(ctrl)

						if matches != 0 {
							// empty spot -> element can't be in Set3 (see Contains) -> insert
							s := set3nextMatch(&matches)
							group.ctrl = setCTRLat(ctrl, H2, s)
							slot[s] = element
							Set3.resident++
							stillSearchingSpace = false

						}
						grpIdx += 1 // carousel through all groups
						if grpIdx >= grpCnt {
							grpIdx = 0
						}
					}
				}
			}
		}
	}
}
