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
	"github.com/dolthub/maphash"
)

// Set is an open-addressing set
// based on Abseil's flat_hash_map.
type Set[K comparable] struct {
	hashFunction maphash.Hasher[K]
	resident     uint32
	dead         uint32
	elementLimit uint32
	group        []Group[K]
}

// metadata is the h2 metadata array for a group.
// find operations first probe the controls bytes
// to filter candidates before matching keys
//type metadata [groupSize]int8

// group is a group of 16 keys

type Group[K comparable] struct {
	ctrl [groupSize]int8
	slot [groupSize]K
}

const (
	kEmpty    = -128 // 0b10000000
	kDeleted  = -2   // 0b11111110
	kSentinel = -1   // 0b11111111
	// kFull= ... // 0b0hhhhhhh, h = bit from hash value
)

// NewSet constructs a Set.
func NewSet[K comparable](sz uint32) (s *Set[K]) {
	reqNrOfGroups := numGroups(sz)
	s = &Set[K]{
		hashFunction: maphash.NewHasher[K](),
		elementLimit: reqNrOfGroups * maxAvgGroupLoad,
		group:        make([]Group[K], reqNrOfGroups),
	}
	for i := range len(s.group) {
		g := &s.group[i]
		for j := range groupSize {
			g.ctrl[j] = kEmpty
		}
	}
	return
}

// Contains returns true if |element| is present in the |Set|.
func (set *Set[K]) Contains(element K) bool {
	hash := set.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := int64(hash & 0x0000_0000_0000_007f)
	grpIdx := H1 % uint64(len(set.group))
	grpCnt := uint64(len(set.group))
	for {
		ctrl := &(set.group[grpIdx].ctrl)
		slot := &(set.group[grpIdx].slot)
		matches := ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := nextMatch_32(&matches)
			if element == slot[s] {
				return true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = ctlrMatchEmpty(ctrl)
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

// Contains returns true if |element| is present in the |Set|.
func (set *Set[K]) Contains2(element K) bool {
	hash := set.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	grpIdx := H1 % uint64(len(set.group))
	grpCnt := uint64(len(set.group))
	for {
		ctrl := &(set.group[grpIdx].ctrl)
		slot := &(set.group[grpIdx].slot)
		matches := metaMatchH2_64(ctrl, H2)
		for matches != 0 {
			s := nextMatch_64(&matches)
			if element == slot[s] {
				return true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(ctrl)
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

// Add attempts to insert |key| and |value|
func (set *Set[K]) Add(element K) {
	if set.resident >= set.elementLimit {
		set.rehash(set.nextSize())
	}
	hash := set.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := int64(hash & 0x0000_0000_0000_007f)
	grpIdx := H1 % uint64(len(set.group))
	grpCnt := uint64(len(set.group))
	for {
		ctrl := &(set.group[grpIdx].ctrl)
		slot := &(set.group[grpIdx].slot)

		matches := ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := nextMatch_32(&matches)
			if element == slot[s] {
				// found - already in Set, just return
				return
			}
		}

		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = ctlrMatchEmpty(ctrl)

		if matches != 0 {
			// empty spot -> element can't be in Set (see Contains) -> insert
			s := nextMatch_32(&matches)
			ctrl[s] = int8(H2)
			slot[s] = element
			set.resident++
			return

		}
		grpIdx += 1 // carousel through all groups
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

// Add attempts to insert |key| and |value|
func (set *Set[K]) Add2(element K) {
	if set.resident >= set.elementLimit {
		set.rehash(set.nextSize())
	}
	hash := set.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	grpIdx := H1 % uint64(len(set.group))
	for {
		ctrl := &(set.group[grpIdx].ctrl)
		slot := &(set.group[grpIdx].slot)

		matches := metaMatchH2_64(ctrl, H2)

		var bitmask uint64 = 1
		// look for matches
		for i := 0; i < groupSize; i++ {
			if (matches & bitmask) != 0 {
				if element == slot[i] {
					// found - already in, just return
					return
				}
			}
			bitmask <<= 1
		}

		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(ctrl)

		if matches != 0 {
			// empty spot -> element can't be in Set (see Contains) -> insert
			var bitmask uint64 = 1
			// find first match
			for i := 0; i < groupSize; i++ {
				if (matches & bitmask) != 0 {
					ctrl[i] = int8(H2)
					slot[i] = element
					set.resident++
					return
				}
				bitmask <<= 1
			}
		}
		grpIdx += 1 // carousel through all groups
		if grpIdx >= uint64(len(set.group)) {
			grpIdx = 0
		}
	}
}

// Remove attempts to remove |element|, returns true if the |element| was in the |Set|
func (set *Set[K]) Remove(element K) (ok bool) {
	hash := set.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	grpIdx := H1 % uint64(len(set.group))
	for {
		ctrl := &(set.group[grpIdx].ctrl)
		slot := &(set.group[grpIdx].slot)
		matches := metaMatchH2_64(ctrl, H2)
		for matches != 0 {
			var bitmask uint64 = 1
			// find first match
			for i := 0; i < groupSize; i++ {
				if (matches & bitmask) != 0 {
					if element == slot[i] {
						ok = true
						// optimization: if |m.ctrl[g]| contains any empty
						// metadata bytes, we can physically delete |key|
						// rather than placing a tombstone.
						// The observation is that any probes into group |g|
						// would already be terminated by the existing empty
						// slot, and therefore reclaiming slot |s| will not
						// cause premature termination of probes into |g|.
						if metaMatchEmpty_64(ctrl) != 0 {
							ctrl[i] = kEmpty
							set.resident--
						} else {
							ctrl[i] = kDeleted
							set.dead++
						}
						var k K
						slot[i] = k
						return
					}
				}
				bitmask <<= 1
			}
		}

		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(ctrl)
		if matches != 0 { // |key| absent
			ok = false
			return
		}
		grpIdx += 1 // linear probing
		if grpIdx >= uint64(len(set.group)) {
			grpIdx = 0
		}
	}
}

// Iter iterates the elements of the Map, passing them to the callback.
// It guarantees that any key in the Map will be visited only once, and
// for un-mutated Maps, every key will be visited once. If the Map is
// Mutated during iteration, mutations will be reflected on return from
// Iter, but the set of keys visited by Iter is non-deterministic.
func (set *Set[K]) Iter(callBack func(element K) (stop bool)) {
	// take a consistent view of the table in case
	// we rehash during iteration
	data := set.group
	// pick a random starting group
	grpIdx := randIntN(len(data))
	for n := 0; n < len(data); n++ {
		ctrl := &(data[grpIdx].ctrl)
		slot := &(data[grpIdx].slot)
		for i, ctrlByte := range ctrl {
			if ctrlByte == kEmpty || ctrlByte == kDeleted {
				continue
			}
			k := slot[i]
			if stop := callBack(k); stop {
				return
			}
		}
		grpIdx++
		if grpIdx >= uint32(len(data)) {
			grpIdx = 0
		}
	}
}

// Clear removes all elements from the Map.
func (set *Set[K]) Clear() {
	var k K
	for grpidx := range len(set.group) {
		d := &(set.group[grpidx])
		for j := range groupSize {
			d.ctrl[j] = kEmpty
			d.slot[j] = k
		}
	}
	set.resident, set.dead = 0, 0
}

// Count returns the number of elements in the Map.
func (set *Set[K]) Count() int {
	return int(set.resident - set.dead)
}

// Capacity returns the number of additional elements
// the can be added to the Map before resizing.
func (set *Set[K]) Capacity() int {
	return int(set.elementLimit - set.resident)
}

// find returns the location of |key| if present, or its insertion location if absent.
// for performance, find is manually inlined into public methods.
func (set *Set[K]) find(key K) (g, s uint64, ok bool) {
	//g = probeStart2(hi, len(set.data))
	hash := set.hashFunction.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	g = H1 % uint64(len(set.group))
	for {
		ctrl := &set.group[g].ctrl
		slot := &set.group[g].slot
		matches := metaMatchH2_64(ctrl, H2)
		for matches != 0 {
			s = nextMatch_64(&matches)
			if key == slot[s] {
				return g, s, true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(ctrl)
		if matches != 0 {
			s = nextMatch_64(&matches)
			return g, s, false
		}
		g += 1 // linear probing
		if g >= uint64(len(set.group)) {
			g = 0
		}
	}
}

func (set *Set[K]) nextSize() (n uint32) {
	n = uint32(len(set.group)) * 2
	if set.dead >= (set.resident / 2) {
		n = uint32(len(set.group))
	}
	return
}

func (set *Set[K]) rehash(n uint32) {
	//groups, ctrl := set.slots, set.ctrl
	old_groups := set.group
	set.hashFunction = maphash.NewSeed(set.hashFunction)
	set.elementLimit = n * maxAvgGroupLoad
	set.resident, set.dead = 0, 0
	set.group = make([]Group[K], n)
	for i := range len(set.group) {
		group := &set.group[i]
		for j := range groupSize {
			group.ctrl[j] = kEmpty
		}
	}
	for _, old_grp := range old_groups {
		for s := range groupSize {
			c := old_grp.ctrl[s]
			if c == kEmpty || c == kDeleted {
				continue
			}
			set.Add(old_grp.slot[s])
		}
	}
}

func (set *Set[K]) loadFactor() float32 {
	slots := float32(len(set.group) * groupSize)
	return float32(set.resident-set.dead) / slots
}
