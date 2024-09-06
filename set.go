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

//const (
// maxLoadFactor = float32(maxAvgGroupLoad) / float32(groupSize)
//)

// Set is an open-addressing set
// based on Abseil's flat_hash_map.
type Set[K comparable] struct {
	//	ctrl     []metadata
	//	slots    []slot[K]
	data     []Group[K]
	hash     maphash.Hasher[K]
	resident uint32
	dead     uint32
	limit    uint32
}

// metadata is the h2 metadata array for a group.
// find operations first probe the controls bytes
// to filter candidates before matching keys
//type metadata [groupSize]int8

// group is a group of 16 keys
//
//	type slot[K comparable] struct {
//		keys [groupSize]K
//	}
type slot[K comparable] [groupSize]K

type Group[K comparable] struct {
	ctrl  [groupSize]int8
	slots [groupSize]K
}

const (
	kEmpty    = -128 // 0b10000000
	kDeleted  = -2   // 0b11111110
	kSentinel = -1   // 0b11111111
	// kFull= ... // 0b0hhhhhhh, h = bit from hash value
)

//const (
// h1Mask    uint64 = 0xffff_ffff_ffff_ff80
// h2Mask    uint64 = 0x0000_0000_0000_007f
// empty     int8   = -128 // 0b1000_0000
// tombstone int8   = -2   // 0b1111_1110
//)

// h1 is a 57 bit hash prefix
//type h1 uint64

// h2 is a 7 bit hash suffix
//type h2 int8

// NewSet constructs a Set.
func NewSet[K comparable](sz uint32) (s *Set[K]) {
	groups := numGroups(sz)
	s = &Set[K]{
		//ctrl:  make([]metadata, groups),
		//slots: make([]slot[K], groups),
		data:  make([]Group[K], groups),
		hash:  maphash.NewHasher[K](),
		limit: groups * maxAvgGroupLoad,
	}
	for i := range s.data {
		for j := range groupSize {
			s.data[i].ctrl[j] = kEmpty
		}
	}
	return
}

// Has returns true if |key| is present in |set|.
func (set *Set[K]) Has(key K) (ok bool) {
	hash := set.hash.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	group := H1 % uint64(len(set.data))
	for { // inlined find loop
		ctrl := &(set.data[group].ctrl)
		slot := &(set.data[group].slots)
		matches := metaMatchH2_64(ctrl, H2)
		for matches != 0 {
			s := nextMatch_64(&matches)
			if key == slot[s] {
				ok = true
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(ctrl)
		if matches != 0 {
			ok = false
			return
		}
		group += 1 // linear probing
		if group >= uint64(len(set.data)) {
			group = 0
		}
	}
}

// Put attempts to insert |key| and |value|
func (set *Set[K]) Put(key K) {
	if set.resident >= set.limit {
		set.rehash(set.nextSize())
	}
	//hi, lo := splitHash(set.hash.Hash(key))
	//_, lo := splitHash(set.hash.Hash(key))
	hash := set.hash.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	//g := probeStart(hi, len(set.slots))
	g := H1 % uint64(len(set.data))
	for { // inlined find loop
		ctrl := &(set.data[g].ctrl)
		slot := &(set.data[g].slots)

		matches := metaMatchH2_64(ctrl, H2)
		for matches != 0 {
			s := nextMatch_64(&matches)
			if key == slot[s] {
				// found - already is set
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(ctrl)
		if matches != 0 { // insert
			s := nextMatch_64(&matches)
			ctrl[s] = int8(H2)
			slot[s] = key
			set.resident++
			return
		}
		g += 1 // linear probing
		if g >= uint64(len(set.data)) {
			g = 0
		}
	}
}

// Delete attempts to remove |key|, returns true successful.
func (set *Set[K]) Delete(key K) (ok bool) {
	hash := set.hash.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	g := H1 % uint64(len(set.data))

	for {
		matches := metaMatchH2_64(&set.data[g].ctrl, H2)
		for matches != 0 {
			s := nextMatch_64(&matches)
			if key == set.data[g].slots[s] {
				ok = true
				// optimization: if |m.ctrl[g]| contains any empty
				// metadata bytes, we can physically delete |key|
				// rather than placing a tombstone.
				// The observation is that any probes into group |g|
				// would already be terminated by the existing empty
				// slot, and therefore reclaiming slot |s| will not
				// cause premature termination of probes into |g|.
				if metaMatchEmpty_64(&set.data[g].ctrl) != 0 {
					set.data[g].ctrl[s] = kEmpty
					set.resident--
				} else {
					set.data[g].ctrl[s] = kDeleted
					set.dead++
				}
				var k K
				set.data[g].slots[s] = k
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty_64(&set.data[g].ctrl)
		if matches != 0 { // |key| absent
			ok = false
			return
		}
		g += 1 // linear probing
		if g >= uint64(len(set.data)) {
			g = 0
		}
	}
}

// Iter iterates the elements of the Map, passing them to the callback.
// It guarantees that any key in the Map will be visited only once, and
// for un-mutated Maps, every key will be visited once. If the Map is
// Mutated during iteration, mutations will be reflected on return from
// Iter, but the set of keys visited by Iter is non-deterministic.
func (set *Set[K]) Iter(cb func(k K) (stop bool)) {
	// take a consistent view of the table in case
	// we rehash during iteration
	//ctrl, groups := set.ctrl, set.slots
	data := set.data
	// pick a random starting group
	g := randIntN(len(data))
	for n := 0; n < len(data); n++ {
		for s, c := range data[g].ctrl {
			if c == kEmpty || c == kDeleted {
				continue
			}
			k := data[g].slots[s]
			if stop := cb(k); stop {
				return
			}
		}
		g++
		if g >= uint32(len(data)) {
			g = 0
		}
	}
}

// Clear removes all elements from the Map.
func (set *Set[K]) Clear() {
	var k K
	for grpidx := range len(set.data) {
		d := &(set.data[grpidx])
		for j := range groupSize {
			d.ctrl[j] = kEmpty
			d.slots[j] = k
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
	return int(set.limit - set.resident)
}

// find returns the location of |key| if present, or its insertion location if absent.
// for performance, find is manually inlined into public methods.
func (set *Set[K]) find(key K) (g, s uint64, ok bool) {
	//g = probeStart2(hi, len(set.data))
	hash := set.hash.Hash(key)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := hash & 0x0000_0000_0000_007f
	g = H1 % uint64(len(set.data))
	for {
		ctrl := &set.data[g].ctrl
		slot := &set.data[g].slots
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
		if g >= uint64(len(set.data)) {
			g = 0
		}
	}
}

func (set *Set[K]) nextSize() (n uint32) {
	n = uint32(len(set.data)) * 2
	if set.dead >= (set.resident / 2) {
		n = uint32(len(set.data))
	}
	return
}

func (set *Set[K]) rehash(n uint32) {
	//groups, ctrl := set.slots, set.ctrl
	data := set.data
	//set.slots = make([]slot[K], n)
	//set.ctrl = make([]metadata, n)
	set.data = make([]Group[K], n)
	//for i := range set.ctrl {
	//	set.ctrl[i] = newEmptyMetadata()
	//}
	for i := range set.data {
		for j := range groupSize {
			set.data[i].ctrl[j] = kEmpty
		}
	}
	set.hash = maphash.NewSeed(set.hash)
	set.limit = n * maxAvgGroupLoad
	set.resident, set.dead = 0, 0
	//for g := range ctrl {
	for _, g := range data {
		for s := range groupSize {
			c := g.ctrl[s]
			if c == kEmpty || c == kDeleted {
				continue
			}
			set.Put(g.slots[s])
		}
	}
}

func (set *Set[K]) loadFactor() float32 {
	slots := float32(len(set.data) * groupSize)
	return float32(set.resident-set.dead) / slots
}

func probeStart2(hi h1, groups int) uint32 {
	return uint32(uint64(hi) % uint64(groups))
}
