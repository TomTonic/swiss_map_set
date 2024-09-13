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
	"fmt"
	"iter"
	"math/bits"
	"strings"

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

var set3hextable = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}

func (g set3Group[k]) String() string {
	var builder strings.Builder
	var mask uint64 = 0xFF
	shr := 0
	builder.WriteString("[")
	for i := range set3groupSize {
		b := g.ctrl & mask
		b >>= shr
		mask <<= 8
		shr += 8
		switch b {
		case set3Empty:
			builder.WriteString("__")
		case set3Deleted:
			builder.WriteString("XX")
		default:
			builder.WriteString(set3hextable[b>>4])
			builder.WriteString(set3hextable[b&0x0f])
		}
		if i < set3groupSize-1 {
			builder.WriteString("|")
		}
	}
	builder.WriteString("]->{")
	for i, v := range g.slot {
		builder.WriteString(fmt.Sprintf("%v", v))
		if i < set3groupSize-1 {
			builder.WriteString("|")
		}
	}
	builder.WriteString("}")
	return builder.String()
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

func (this Set3[K]) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	total := this.Count()
	cnt := uint32(0)
	for e := range this.MutableRange() {
		builder.WriteString(fmt.Sprintf("%v", e))
		if cnt < total-1 {
			builder.WriteString(",")
		}
		cnt++
	}
	builder.WriteString("}")
	return builder.String()
}

func NewSet3[K comparable]() (s *Set3[K]) {
	return NewSet3WithSize[K](21)
}

// NewSet3WithSize constructs a Set3.
func NewSet3WithSize[K comparable](size uint32) (s *Set3[K]) {
	reqNrOfGroups := calcReqNrOfGroups(size)
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

func calcReqNrOfGroups(size uint32) int {
	reqNrOfGroups := int((float64(size) + set3maxAvgGroupLoad - 1) / set3maxAvgGroupLoad)
	if reqNrOfGroups == 0 {
		reqNrOfGroups = 1
	}
	return reqNrOfGroups
}

// AsSet3 constructs a Set3 from an array/slice.
func AsSet3[K comparable](data []K) (set *Set3[K]) {
	set = NewSet3WithSize[K](uint32(len(data)))
	for _, e := range data {
		set.Add(e)
	}
	return
}

func (set3 *Set3[K]) Clone() (s *Set3[K]) {
	s = &Set3[K]{
		hashFunction: set3.hashFunction,
		elementLimit: set3.elementLimit,
		resident:     set3.resident,
		dead:         set3.dead,
		group:        set3.fullCopyGroups(),
	}
	return
}

func (set3 *Set3[K]) fullCopyGroups() []set3Group[K] {
	result := make([]set3Group[K], len(set3.group))
	copy(result, set3.group)
	return result
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

func (this *Set3[K]) ContainsAll(that *Set3[K]) bool {
	if this == that {
		return true
	}
	if this.Count() < that.Count() {
		return false
	}
	for e := range that.MutableRange() {
		if !this.Contains(e) {
			return false
		}
	}
	return true
}

func (this *Set3[K]) ContainsAllFrom(data []K) bool {
	for _, e := range data {
		if !this.Contains(e) {
			return false
		}
	}
	return true
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
		groups := s.fullCopyGroups()
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

// Add attempts to insert |element|
func (Set3 *Set3[K]) Add(element K) {
	if Set3.resident >= Set3.elementLimit {
		Set3.rehashToNumGroups(Set3.nextSize())
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

// Attempts to insert all elements from |that| into |this| Set3.
func (this *Set3[K]) AddAll(that *Set3[K]) {
	for e := range that.MutableRange() {
		this.Add(e)
	}
}

// Attempts to insert all elements from |that| into |this| Set3.
func (this *Set3[K]) AddAllFrom(data []K) {
	for _, e := range data {
		this.Add(e)
	}
}

// Creates a new Set3 as a mathematical union of the elements from |this| and |that|.
func (this *Set3[K]) Union(that *Set3[K]) *Set3[K] {
	potentialSize := this.Count() + that.Count()
	result := NewSet3WithSize[K](potentialSize)
	for e := range this.MutableRange() {
		result.Add(e)
	}
	for e := range that.MutableRange() {
		result.Add(e)
	}
	return result
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
func (this *Set3[K]) Remove(element K) bool {
	hash := this.hashFunction.Hash(element)
	H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
	H2 := (hash & 0x0000_0000_0000_007f)
	grpCnt := uint64(len(this.group))
	grpIdx := H1 % grpCnt
	for {
		group := &this.group[grpIdx]
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
					this.resident--
				} else {
					group.ctrl = setCTRLat(ctrl, set3Deleted, s)
					this.dead++
					/*
						// unfortunately, this is an invalid optimization, as the algorithm might stop searching for elements to early.
						// if they spilled over in the next group, we unfortunately need all the tumbstones...
						if group.ctrl == set3AllDeleted {
							group.ctrl = set3AllEmpty
							this.dead -= set3groupSize
							this.resident -= set3groupSize
						}
					*/
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

// Attempts to insert all elements from |that| into |this| Set3.
func (this *Set3[K]) RemoveAll(that *Set3[K]) {
	for e := range that.MutableRange() {
		this.Remove(e)
	}
}

// Attempts to insert all elements from |that| into |this| Set3.
func (this *Set3[K]) RemoveAllFrom(data []K) {
	for _, e := range data {
		this.Remove(e)
	}
}

// Creates a new Set3 as a mathematical difference between |this| and |that|; i.e. the result contains nodes that are in |this| but not in |that|.
func (this *Set3[K]) Difference(that *Set3[K]) *Set3[K] {
	potentialSize := this.Count()
	result := NewSet3WithSize[K](potentialSize)
	for e := range this.MutableRange() {
		if !that.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

// Clear removes all elements from the Set3.
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

// Creates a new Set3 as a mathematical intersection between |this| and |that|; i.e. the result contains nodes that are in |this| and in |that|.
func (this *Set3[K]) Intersection(that *Set3[K]) *Set3[K] {
	var smallerSet *Set3[K]
	var biggerSet *Set3[K]

	if this.Count() < that.Count() {
		smallerSet = this
		biggerSet = that
	} else {
		smallerSet = that
		biggerSet = this
	}

	potentialSize := smallerSet.Count()
	result := NewSet3WithSize[K](potentialSize)
	for e := range smallerSet.ImmutableRange() {
		if biggerSet.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

// Count returns the number of elements in the Map.
func (Set3 *Set3[K]) Count() uint32 {
	return Set3.resident - Set3.dead
}

func (Set3 *Set3[K]) nextSize() (n uint32) {
	n = uint32(len(Set3.group)) * 2
	if Set3.dead >= (Set3.resident / 2) {
		n = uint32(len(Set3.group))
	}
	return
}

// Rorganize Set3 for better performance for current number of elements.
func (this *Set3[K]) Rehash() {
	this.rehashToNumGroups(this.Count())
}

// Rorganize Set3 for better performance. If |newSize| is smaller than the current number of elements in the |Set3|, this function does nothing.
func (this *Set3[K]) RehashTo(newSize uint32) {
	if newSize < this.Count() {
		return
	}
	newNumGroups := uint32(calcReqNrOfGroups(newSize))
	this.rehashToNumGroups(newNumGroups)
}

func (this *Set3[K]) rehashToNumGroups(newNumGroups uint32) {
	old_groups := this.fullCopyGroups()
	this.hashFunction = maphash.NewSeed(this.hashFunction)
	this.elementLimit = uint32(float64(newNumGroups) * set3maxAvgGroupLoad)
	this.resident, this.dead = 0, 0
	this.group = make([]set3Group[K], newNumGroups)
	for i := range len(this.group) {
		this.group[i].ctrl = set3AllEmpty
	}
	grpCnt := uint64(len(this.group))
	for _, old_grp := range old_groups {
		if old_grp.ctrl&set3hiBits != set3hiBits { // not all empty or deleted
			for s := range set3groupSize {
				if elementAt(old_grp.ctrl, s) {
					// inlined and reduced Add instead of Set3.Add(old_grp.slot[s])
					element := old_grp.slot[s]

					hash := this.hashFunction.Hash(element)
					H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
					H2 := (hash & 0x0000_0000_0000_007f)
					grpIdx := H1 % uint64(len(this.group))
					stillSearchingSpace := true
					for stillSearchingSpace {
						group := &this.group[grpIdx]
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
							this.resident++
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
