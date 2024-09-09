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
	"math/rand/v2"

	"github.com/dolthub/maphash"
)

// Set2 is an open-addressing Set2
// based on Abseil's flat_hash_map.
type Set2[K comparable] struct {
	hashFunction maphash.Hasher[K]
	groupctrl    [][groupSize]int8
	groupslot    [][groupSize]K
	resident     uint32
	dead         uint32
	elementLimit uint32
}

// NewSet2 constructs a Set2.
func NewSet2[K comparable](sz uint32) (s *Set2[K]) {
	reqNrOfGroups := numGroups(sz)
	s = &Set2[K]{
		hashFunction: maphash.NewHasher[K](),
		groupctrl:    make([][groupSize]int8, reqNrOfGroups),
		groupslot:    make([][groupSize]K, reqNrOfGroups),
		elementLimit: reqNrOfGroups * maxAvgGroupLoad,
	}
	for i := range len(s.groupctrl) {
		g := &s.groupctrl[i]
		for j := range groupSize {
			g[j] = kEmpty
		}
	}
	return
}

// Contains returns true if |element| is present in the |Set2|.
func (Set2 *Set2[K]) Contains(element K) bool {
	hash := Set2.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpIdx := H1 % uint64(len(Set2.groupctrl))
	grpCnt := uint64(len(Set2.groupctrl))
	for {
		ctrl := &(Set2.groupctrl[grpIdx])
		slot := &(Set2.groupslot[grpIdx])
		matches := ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := nextMatch(&matches)
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

// Add attempts to insert |key| and |value|
func (Set2 *Set2[K]) Add(element K) {
	if Set2.resident >= Set2.elementLimit {
		Set2.rehash(Set2.nextSize())
	}
	hash := Set2.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpIdx := H1 % uint64(len(Set2.groupctrl))
	grpCnt := uint64(len(Set2.groupctrl))
	for {
		ctrl := &(Set2.groupctrl[grpIdx])
		slot := &(Set2.groupslot[grpIdx])

		matches := ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := nextMatch(&matches)
			if element == slot[s] {
				// found - already in Set2, just return
				return
			}
		}

		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = ctlrMatchEmpty(ctrl)

		if matches != 0 {
			// empty spot -> element can't be in Set2 (see Contains) -> insert
			s := nextMatch(&matches)
			ctrl[s] = int8(H2)
			slot[s] = element
			Set2.resident++
			return

		}
		grpIdx += 1 // carousel through all groups
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

// Remove attempts to remove |element|, returns true if the |element| was in the |Set2|
func (Set2 *Set2[K]) Remove(element K) bool {
	hash := Set2.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpIdx := H1 % uint64(len(Set2.groupctrl))
	grpCnt := uint64(len(Set2.groupctrl))
	for {
		ctrl := &(Set2.groupctrl[grpIdx])
		slot := &(Set2.groupslot[grpIdx])
		matches := ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s := nextMatch(&matches)
			if element == slot[s] {
				// found - already in Set2, just return
				// optimization: if |m.ctrl[g]| contains any empty
				// metadata bytes, we can physically delete |element|
				// rather than placing a tombstone.
				// The observation is that any probes into group |g|
				// would already be terminated by the existing empty
				// slot, and therefore reclaiming slot |s| will not
				// cause premature termination of probes into |g|.
				if ctlrMatchEmpty(ctrl) != 0 {
					ctrl[s] = kEmpty
					Set2.resident--
				} else {
					ctrl[s] = kDeleted
					Set2.dead++
				}
				var k K
				slot[s] = k
				return true
			}
		}

		// |element| is not in group |g|,
		// stop probing if we see an empty slot
		matches = ctlrMatchEmpty(ctrl)
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

// Iter iterates the elements of the Map, passing them to the callback.
// It guarantees that any key in the Map will be visited only once, and
// for un-mutated Maps, every key will be visited once. If the Map is
// Mutated during iteration, mutations will be reflected on return from
// Iter, but the Set2 of keys visited by Iter is non-deterministic.
func (Set2 *Set2[K]) Iter(callBack func(element K) (stop bool)) {
	// take a consistent view of the table in case
	// we rehash during iteration
	groupctrl := Set2.groupctrl
	groupslot := Set2.groupslot
	// pick a random starting group
	grpIdx := rand.Uint32N(uint32(len(groupctrl)))
	for n := 0; n < len(groupctrl); n++ {
		ctrl := &(groupctrl[grpIdx])
		slot := &(groupslot[grpIdx])
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
		if grpIdx >= uint32(len(groupctrl)) {
			grpIdx = 0
		}
	}
}

// Clear removes all elements from the Map.
func (Set2 *Set2[K]) Clear() {
	var k K
	for grpidx := range len(Set2.groupctrl) {
		groupctrl := &(Set2.groupctrl[grpidx])
		for j := range groupSize {
			groupctrl[j] = kEmpty
		}
	}
	for grpidx := range len(Set2.groupslot) {
		groupslot := &(Set2.groupslot[grpidx])
		for j := range groupSize {
			groupslot[j] = k
		}
	}
	Set2.resident, Set2.dead = 0, 0
}

// Count returns the number of elements in the Map.
func (Set2 *Set2[K]) Count() int {
	return int(Set2.resident - Set2.dead)
}

// Capacity returns the number of additional elements
// the can be added to the Map before resizing.
func (Set2 *Set2[K]) Capacity() int {
	return int(Set2.elementLimit - Set2.resident)
}

// find returns the location of |key| if present, or its insertion location if absent.
// for performance, find is manually inlined into public methods.
func (Set2 *Set2[K]) find(key K) (g uint64, s int, ok bool) {
	//g = probeStart2(hi, len(Set2.data))
	hash := Set2.hashFunction.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	g = H1 % uint64(len(Set2.groupctrl))
	for {
		ctrl := &Set2.groupctrl[g]
		slot := &Set2.groupslot[g]
		matches := ctlrMatchH2(ctrl, H2)
		for matches != 0 {
			s = nextMatch(&matches)
			if key == slot[s] {
				return g, s, true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = ctlrMatchEmpty(ctrl)
		if matches != 0 {
			s = nextMatch(&matches)
			return g, s, false
		}
		g += 1 // linear probing
		if g >= uint64(len(Set2.groupctrl)) {
			g = 0
		}
	}
}

func (Set2 *Set2[K]) nextSize() (n uint32) {
	n = uint32(len(Set2.groupctrl)) * 2
	if Set2.dead >= (Set2.resident / 2) {
		n = uint32(len(Set2.groupctrl))
	}
	return
}

func (Set2 *Set2[K]) rehash(n uint32) {
	//groups, ctrl := Set2.slots, Set2.ctrl
	old_groups_ctrl := Set2.groupctrl
	old_groups_slot := Set2.groupslot
	Set2.hashFunction = maphash.NewSeed(Set2.hashFunction)
	Set2.elementLimit = n * maxAvgGroupLoad
	Set2.resident, Set2.dead = 0, 0
	Set2.groupctrl = make([][groupSize]int8, n)
	Set2.groupslot = make([][groupSize]K, n)

	for i := range len(Set2.groupctrl) {
		groupctrl := &Set2.groupctrl[i]
		for j := range groupSize {
			groupctrl[j] = kEmpty
		}
	}
	grpCnt := uint64(len(Set2.groupctrl))
	for idx := range len(old_groups_ctrl) {
		old_ctrl := &old_groups_ctrl[idx]
		old_slot := &old_groups_slot[idx]
		for s := range groupSize {
			c := old_ctrl[s]
			if c == kEmpty || c == kDeleted {
				continue
			}
			// inlined and reduced Add instead of Set2.Add(old_grp.slot[s])

			element := old_slot[s]

			hash := Set2.hashFunction.Hash(element)
			H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
			H2 := int64(hash & 0x0000_0000_0000_007f)
			grpIdx := H1 % uint64(len(Set2.groupctrl))
			stillSearchingSpace := true
			for stillSearchingSpace {
				ctrl := &(Set2.groupctrl[grpIdx])
				slot := &(Set2.groupslot[grpIdx])

				// optimization: we know it cannot exist in the Set2 already so skip
				// searching for the hashcode and start searching for an empty slot
				// immediately
				matches := ctlrMatchEmpty(ctrl)

				if matches != 0 {
					// empty spot -> element can't be in Set2 (see Contains) -> insert
					s := nextMatch(&matches)
					ctrl[s] = int8(H2)
					slot[s] = element
					Set2.resident++
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

func (Set2 *Set2[K]) loadFactor() float32 {
	slots := float32(len(Set2.groupctrl) * groupSize)
	return float32(Set2.resident-Set2.dead) / slots
}

/*
// numGroups returns the minimum number of groups needed to store |n| elems.
func numGroups(n uint32) (groups uint32) {
	groups = (n + maxAvgGroupLoad - 1) / maxAvgGroupLoad
	if groups == 0 {
		groups = 1
	}
	return
}
*/
