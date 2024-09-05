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
	ctrl     []metadata
	groups   []sgroup[K]
	hash     maphash.Hasher[K]
	resident uint32
	dead     uint32
	limit    uint32
}

// metadata is the h2 metadata array for a group.
// find operations first probe the controls bytes
// to filter candidates before matching keys
//type metadata [groupSize]int8

// group is a group of 16 key-value pairs
type sgroup[K comparable] struct {
	keys [groupSize]K
}

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
		ctrl:   make([]metadata, groups),
		groups: make([]sgroup[K], groups),
		hash:   maphash.NewHasher[K](),
		limit:  groups * maxAvgGroupLoad,
	}
	for i := range s.ctrl {
		s.ctrl[i] = newEmptyMetadata()
	}
	return
}

// Has returns true if |key| is present in |m|.
func (set *Set[K]) Has(key K) (ok bool) {
	hi, lo := splitHash(set.hash.Hash(key))
	g := probeStart(hi, len(set.groups))
	for { // inlined find loop
		matches := metaMatchH2(&set.ctrl[g], lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == set.groups[g].keys[s] {
				ok = true
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(&set.ctrl[g])
		if matches != 0 {
			ok = false
			return
		}
		g += 1 // linear probing
		if g >= uint32(len(set.groups)) {
			g = 0
		}
	}
}

// Put attempts to insert |key| and |value|
func (set *Set[K]) Put(key K) {
	if set.resident >= set.limit {
		set.rehash(set.nextSize())
	}
	hi, lo := splitHash(set.hash.Hash(key))
	g := probeStart(hi, len(set.groups))
	for { // inlined find loop
		matches := metaMatchH2(&set.ctrl[g], lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == set.groups[g].keys[s] { // update
				// m.groups[g].keys[s] = key     // superflouos
				// m.groups[g].values[s] = value // no map - vo values
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(&set.ctrl[g])
		if matches != 0 { // insert
			s := nextMatch(&matches)
			set.groups[g].keys[s] = key
			//m.groups[g].values[s] = value
			set.ctrl[g][s] = int8(lo)
			set.resident++
			return
		}
		g += 1 // linear probing
		if g >= uint32(len(set.groups)) {
			g = 0
		}
	}
}

// Delete attempts to remove |key|, returns true successful.
func (set *Set[K]) Delete(key K) (ok bool) {
	hi, lo := splitHash(set.hash.Hash(key))
	g := probeStart(hi, len(set.groups))
	for {
		matches := metaMatchH2(&set.ctrl[g], lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == set.groups[g].keys[s] {
				ok = true
				// optimization: if |m.ctrl[g]| contains any empty
				// metadata bytes, we can physically delete |key|
				// rather than placing a tombstone.
				// The observation is that any probes into group |g|
				// would already be terminated by the existing empty
				// slot, and therefore reclaiming slot |s| will not
				// cause premature termination of probes into |g|.
				if metaMatchEmpty(&set.ctrl[g]) != 0 {
					set.ctrl[g][s] = empty
					set.resident--
				} else {
					set.ctrl[g][s] = tombstone
					set.dead++
				}
				var k K
				set.groups[g].keys[s] = k
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(&set.ctrl[g])
		if matches != 0 { // |key| absent
			ok = false
			return
		}
		g += 1 // linear probing
		if g >= uint32(len(set.groups)) {
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
	ctrl, groups := set.ctrl, set.groups
	// pick a random starting group
	g := randIntN(len(groups))
	for n := 0; n < len(groups); n++ {
		for s, c := range ctrl[g] {
			if c == empty || c == tombstone {
				continue
			}
			k := groups[g].keys[s]
			if stop := cb(k); stop {
				return
			}
		}
		g++
		if g >= uint32(len(groups)) {
			g = 0
		}
	}
}

// Clear removes all elements from the Map.
func (set *Set[K]) Clear() {
	for i, c := range set.ctrl {
		for j := range c {
			set.ctrl[i][j] = empty
		}
	}
	var k K
	for i := range set.groups {
		g := &set.groups[i]
		for i := range g.keys {
			g.keys[i] = k
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
func (set *Set[K]) find(key K, hi h1, lo h2) (g, s uint32, ok bool) {
	g = probeStart(hi, len(set.groups))
	for {
		matches := metaMatchH2(&set.ctrl[g], lo)
		for matches != 0 {
			s = nextMatch(&matches)
			if key == set.groups[g].keys[s] {
				return g, s, true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(&set.ctrl[g])
		if matches != 0 {
			s = nextMatch(&matches)
			return g, s, false
		}
		g += 1 // linear probing
		if g >= uint32(len(set.groups)) {
			g = 0
		}
	}
}

func (set *Set[K]) nextSize() (n uint32) {
	n = uint32(len(set.groups)) * 2
	if set.dead >= (set.resident / 2) {
		n = uint32(len(set.groups))
	}
	return
}

func (set *Set[K]) rehash(n uint32) {
	groups, ctrl := set.groups, set.ctrl
	set.groups = make([]sgroup[K], n)
	set.ctrl = make([]metadata, n)
	for i := range set.ctrl {
		set.ctrl[i] = newEmptyMetadata()
	}
	set.hash = maphash.NewSeed(set.hash)
	set.limit = n * maxAvgGroupLoad
	set.resident, set.dead = 0, 0
	for g := range ctrl {
		for s := range ctrl[g] {
			c := ctrl[g][s]
			if c == empty || c == tombstone {
				continue
			}
			set.Put(groups[g].keys[s])
		}
	}
}

func (set *Set[K]) loadFactor() float32 {
	slots := float32(len(set.groups) * groupSize)
	return float32(set.resident-set.dead) / slots
}

// numGroups returns the minimum number of groups needed to store |n| elems.
//func numGroups(n uint32) (groups uint32) {
//	groups = (n + maxAvgGroupLoad - 1) / maxAvgGroupLoad
//	if groups == 0 {
//		groups = 1
//	}
//	return
//}

//func newEmptyMetadata() (meta metadata) {
//	for i := range meta {
//		meta[i] = empty
//	}
//	return
//}

//func splitHash(h uint64) (h1, h2) {
//	return h1((h & h1Mask) >> 7), h2(h & h2Mask)
//}

//func probeStart(hi h1, groups int) uint32 {
//	return fastModN(uint32(hi), uint32(groups))
//}

// lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
//func fastModN(x, n uint32) uint32 {
//	return uint32((uint64(x) * uint64(n)) >> 32)
//}
