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

	set3loBits uint64 = 0x0101010101010101
	set3hiBits uint64 = 0x8080808080808080

	set3AllEmpty   uint64 = 0x8080808080808080
	set3AllDeleted uint64 = 0xFEFEFEFEFEFEFEFE
	set3Empty      uint64 = 0b0000_1000_0000
	set3Deleted    uint64 = 0b0000_1111_1110
	set3Sentinel   uint64 = 0b0000_1111_1111
)

func set3ctlrMatchH2(m uint64, h uint64) uint64 {
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

type set3Group[T comparable] struct {
	ctrl uint64
	slot [set3groupSize]T
}

var set3hextable = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}

func (this *set3Group[T]) String() string {
	var builder strings.Builder
	var mask uint64 = 0xFF
	shr := 0
	builder.WriteString("[")
	for i := range set3groupSize {
		b := this.ctrl & mask
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
	for i, v := range this.slot {
		builder.WriteString(fmt.Sprintf("%v", v))
		if i < set3groupSize-1 {
			builder.WriteString("|")
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// Set3 is a hash set of type K.
type Set3[T comparable] struct {
	hashFunction maphash.Hasher[T]
	resident     uint32
	dead         uint32
	elementLimit uint32
	group        []set3Group[T]
}

/*
Returns a string representation of the elements of this Set3 in Roster notation (https://en.wikipedia.org/wiki/Set_(mathematics)#Roster_notation).
The order of the elements in the result is arbitrarily.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set) // will print "{2,3,1}" with the numbers in arbitrary order
*/
func (this *Set3[T]) String() string {
	if this == nil {
		return "{nil}"
	}
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

/*
NewSet3 creates a new and empty Set3 with a reasonable default initial capacity. Choose this constructor if you have no idea on how big your set will be.
You can add as many elements to this set as you like, the backing data structure will automatically be reorganized to fit your needs.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
*/
func NewSet3[T comparable]() *Set3[T] {
	return NewSet3WithSize[T](21)
}

/*
NewSet3WithSize creates a new and empty Set3 with a given initial capacity. Choose this constructor if you have a pretty good idea on how big your set will be.
Nonetheless, you can add as many elements to this set as you like, the backing data structure will automatically be reorganized to fit your needs.

Example:

	set1 := NewSet3[int]() // you can put 1 mio. ints in set1 - rehashing will occur
	set2 := NewSet3[int](2_000_000) // you can put 1 mio. ints in set2 - no rehashing will occur
*/
func NewSet3WithSize[T comparable](size uint32) *Set3[T] {
	reqNrOfGroups := calcReqNrOfGroups(size)
	result := &Set3[T]{
		hashFunction: maphash.NewHasher[T](),
		elementLimit: uint32(float64(reqNrOfGroups) * set3maxAvgGroupLoad),
		group:        make([]set3Group[T], reqNrOfGroups),
	}
	for i := range len(result.group) {
		result.group[i].ctrl = set3AllEmpty
	}
	return result
}

func calcReqNrOfGroups(size uint32) int {
	reqNrOfGroups := int((float64(size) + set3maxAvgGroupLoad - 1) / set3maxAvgGroupLoad)
	if reqNrOfGroups == 0 {
		reqNrOfGroups = 1
	}
	return reqNrOfGroups
}

/*
AsSet3 is a convenience constructor to directly create a Set3 from given values. It creates a Set3 with the required capacity and adds all (unique) elements to this Set3.

If the array contains duplicates, the duplicates are omitted. If data is nil, an empty Set3 is returned.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := AsSet3([]int{1, 2, 3}) // set1 and set2 are equal
*/
func AsSet3[T comparable](data []T) *Set3[T] {
	if data == nil {
		return NewSet3[T]()
	}
	result := NewSet3WithSize[T](uint32(len(data))) //nolint:gosec
	for _, e := range data {
		result.Add(e)
	}
	return result
}

/*
Clone creates an exact clone of this Set3. You can manipulate both clones independently.

Cloning is 'cheap' in comparison with creating a new set and adding all elements from this set, as only the backing data structures are copied (no rehashing is applied).

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := set1.Clone() // set2 will be an exact but independent clone of set1
*/
func (this *Set3[T]) Clone() *Set3[T] {
	result := &Set3[T]{
		hashFunction: this.hashFunction,
		elementLimit: this.elementLimit,
		resident:     this.resident,
		dead:         this.dead,
		group:        this.fullCopyGroups(),
	}
	return result
}

func (this *Set3[T]) fullCopyGroups() []set3Group[T] {
	result := make([]set3Group[T], len(this.group))
	copy(result, this.group)
	return result
}

/*
Contains returns true if the element is contained in this Set3.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	b1 := set.Contains(2) // b1 will be true
	b2 := set.Contains(4) // b2 will be false
*/
func (this *Set3[T]) Contains(element T) bool {
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
				return true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = set3ctlrMatchEmpty(ctrl)
		if matches != 0 {
			// there is an empty slot - the element, if it had been added, hat either
			// been found until now or it had been added in the next empty spot -
			// well, this is the next empty spot...
			return false
		}
		grpIdx++ // carousel through all groups
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

/*
Returns true if this Set3 contains all elements from that Set3.

If that Set3 is empty, ContainsAll returns true. If that Set3 is nil, ContainsAll returns true.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	emptyset := NewSet3[int]()
	b := set.ContainsAll(emptyset) // b will be true
*/
func (this *Set3[T]) ContainsAll(that *Set3[T]) bool {
	if this == that {
		return true
	}
	if that == nil {
		// nil is interpreted as empty set
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

/*
Returns true if this Set3 contains all elements from the given data array.

If the length of data is zero, ContainsAllFrom returns true. If data is nil, ContainsAllFrom returns true.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	b := set.ContainsAllFrom([]int{2,3,4}) // b will be false
*/
func (this *Set3[T]) ContainsAllFrom(data []T) bool {
	if data == nil {
		// nil is interpreted as empty set
		return true
	}
	for _, e := range data {
		if !this.Contains(e) {
			return false
		}
	}
	return true
}

/*
Returns true if this Set3 and that Set3 have the same size and contain the same elements.

If that Set3 is nil, Equals returns true if and only if this Set3 is empty.

Example:

	set1 := NewSet3[int]()
	set2 := NewSet3[int]()
	b1 := set1.Equals(set2) // b1 will be true
	set1.Add(7)
	set2.Add(31)
	b2 := set1.Equals(set2) // b2 will be false
*/
func (this *Set3[T]) Equals(that *Set3[T]) bool {
	if this == that {
		return true
	}
	if that == nil {
		// nil is interpreted as empty set
		return this.Count() == 0
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

/*
Iterates over all elements in this Set3.

Caution: If this Set3 is changed during the iteration, the result is unpredictable. So if you want to add or remove elements to or from this Set3 during the itration, choose [ImmutableRange].

Example:

	for elem := range set.MutableRange() {
		// do something with elem...
	}
*/
func (this *Set3[T]) MutableRange() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, group := range this.group {
			ctrl := group.ctrl
			if ctrl&set3hiBits != set3hiBits { // not all empty or deleted
				slot := &(group.slot)
				for i := 0; i < set3groupSize; i++ {
					if isAnElementAt(ctrl, i) {
						if !yield(slot[i]) {
							return
						}
					}
				}
			}
		}
	}
}

/*
Iterates over all elements in this Set3.

Makes an internal copy of the stored elements first, so you can add or remove elements to or from this Set3 during the itration, for example. To avoid this extra copy, e.g., for performance reasons, choose [MutableRange].

Example:

	for elem := range set.ImmutableRange() {
		// do something with elem...
	}
*/
func (this *Set3[T]) ImmutableRange() iter.Seq[T] {
	return func(yield func(T) bool) {
		groups := this.fullCopyGroups()
		for _, group := range groups {
			ctrl := group.ctrl
			if ctrl&set3hiBits != set3hiBits { // not all empty or deleted
				slot := &(group.slot)
				for i := 0; i < set3groupSize; i++ {
					if isAnElementAt(ctrl, i) {
						if !yield(slot[i]) {
							return
						}
					}
				}
			}
		}
	}
}

/*
ToArray allocates an array of type T and adds all elements of this Set3 to it. The order of the elements in the resulting array is arbitrary.

Example:

	set := NewSet3[int]()
	set.Add(7)
	set.Add(31)
	int_array := set.ToArray() // will be an []int of length 2 containing 7 and 31 in arbitrary order
*/
func (this *Set3[T]) ToArray() []T {
	result := make([]T, this.Count())
	i := 0
	for e := range this.MutableRange() {
		result[i] = e
		i++
	}
	return result
}

/*
Inserts the element into this Set3 if it is not yet in this Set3.

Example:

	set := NewSet3[int]()
	set.Add(7)
*/
func (this *Set3[T]) Add(element T) {
	if this.resident >= this.elementLimit {
		this.rehashToNumGroups(this.nextSize())
	}
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
				return
			}
		}

		// element is not in group,
		// stop probing if we see an empty slot
		matches = set3ctlrMatchEmpty(ctrl)

		if matches != 0 {
			// empty spot -> element can't be in Set3 (see Contains) -> insert
			s := set3nextMatch(&matches)
			group.ctrl = setCTRLat(ctrl, H2, s)
			slot[s] = element
			this.resident++
			return

		}
		grpIdx++ // carousel through all groups
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

/*
Inserts all elements from that Set3 that are not yet in this Set3 into this Set3.

If that Set3 is nil, nothing is added to this Set3.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set2 := NewSet3[int]()
	set2.Add(2)
	set2.Add(3)
	set1.AddAll(set2) // set1 will now contain 1, 2, 3
*/
func (this *Set3[T]) AddAll(that *Set3[T]) {
	if that == nil {
		return
	}
	for e := range that.MutableRange() {
		this.Add(e)
	}
}

/*
Inserts all elements from the given data array that are not yet in this Set3 into this Set3.

If data is nil, nothing is added to this Set3.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.AddAllFrom([]int{2,3}]) // set will now contain 1, 2, 3
*/
func (this *Set3[T]) AddAllFrom(data []T) {
	if data == nil {
		return
	}
	for _, e := range data {
		this.Add(e)
	}
}

/*
Creates a new Set3 as a mathematical union of the elements from this Set3 and that Set3.

If that Set3 is nil, Union returns a clone of this Set3.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set2 := NewSet3[int]()
	set2.Add(2)
	set2.Add(3)

	u := set1.Union(set2) // set1 and set2 remain unchanged, u will contain 1, 2, 3
*/
func (this *Set3[T]) Union(that *Set3[T]) *Set3[T] {
	if that == nil {
		return this.Clone()
	}
	potentialSize := this.Count() + that.Count()
	result := NewSet3WithSize[T](potentialSize)
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

func isAnElementAt(ctrl uint64, pos int) bool {
	shift := pos << 3               // *8
	ctrl &= (uint64(0x80) << shift) // clear all other bits
	// if a bit is set, the according byte represented either set3Empty or set3Deleted
	// -> if ctlr is 0 now, the according position stores a value
	return ctrl == 0
}

/*
Removes the given element from this Set3 if it is in this Set3, returns whether or not the element was in this Set3.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Remove(2)	// nothing happens to set
	set.Remove(1)	// set will be empty
	set.Remove(0)	// set will still be empty
*/
func (this *Set3[T]) Remove(element T) bool {
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
				var k T
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
		grpIdx++ // linear probing
		if grpIdx >= grpCnt {
			grpIdx = 0
		}
	}
}

/*
Removes all elements from this Set3 that are in that Set3.

If that Set3 is nil, nothing happens.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := NewSet3[int]()
	set2.Add(3)
	set2.Add(4)
	set1.RemoveAll(set2) // set1 will now contain 1, 2
*/
func (this *Set3[T]) RemoveAll(that *Set3[T]) {
	if that == nil {
		return
	}
	for e := range that.MutableRange() {
		this.Remove(e)
	}
}

/*
Removes all elements from this Set3 that are in the data array.

If data is nil, nothing happens.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.RemoveAllFrom([]int{3,4}]) // set will now contain 1, 2
*/
func (this *Set3[T]) RemoveAllFrom(data []T) {
	if data == nil {
		return
	}
	for _, e := range data {
		this.Remove(e)
	}
}

/*
Creates a new Set3 as a mathematical difference between this and that. The result is a new Set3 that contains elements that are in this Set3 but not in that Set3.

If that Set3 is nil, Difference returns a clone of this Set3.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := NewSet3[int]()
	set2.Add(3)
	set2.Add(4)
	d := set1.Difference(set2) // set1 and set2 are not altered, d will contain 1, 2
*/
func (this *Set3[T]) Difference(that *Set3[T]) *Set3[T] {
	if that == nil {
		return this.Clone()
	}
	potentialSize := this.Count()
	result := NewSet3WithSize[T](potentialSize)
	for e := range this.MutableRange() {
		if !that.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

/*
Clear removes all elements from this Set3.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Clear()  // set will be empty, Count() will return 0
*/
func (this *Set3[T]) Clear() {
	var k T
	for grpidx := range len(this.group) {
		d := &(this.group[grpidx])
		d.ctrl = set3AllEmpty
		for j := range set3groupSize {
			d.slot[j] = k
		}
	}
	this.resident, this.dead = 0, 0
}

/*
Creates a new Set3 as a mathematical intersection between this and that. The result is a new Set3 that contains elements that are in both sets.

If that Set3 is nil, Intersection returns an empty Set3.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := NewSet3[int]()
	set2.Add(3)
	set2.Add(4)
	intersect := set1.Intersection(set2) // set1 and set2 are not altered, intersect will contain 3
*/
func (this *Set3[T]) Intersection(that *Set3[T]) *Set3[T] {
	if that == nil {
		return NewSet3[T]()
	}

	var smallerSet *Set3[T]
	var biggerSet *Set3[T]

	if this.Count() < that.Count() {
		smallerSet = this
		biggerSet = that
	} else {
		smallerSet = that
		biggerSet = this
	}

	potentialSize := smallerSet.Count()
	result := NewSet3WithSize[T](potentialSize)
	for e := range smallerSet.ImmutableRange() {
		if biggerSet.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

/*
Creates a new Set3 as a mathematical intersection between this Set3 and the elements of the data array. The result is a new Set3.

If data is nil, IntersectionFrom returns an empty Set3.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	intersect := set.IntersectionFrom([]int{3,4}]) // set1 and set2 are not altered, intersect will contain 3
*/
func (this *Set3[T]) IntersectionFrom(data []T) *Set3[T] {
	if data == nil {
		return NewSet3[T]()
	}

	var potentialSize uint32

	if this.Count() < uint32(len(data)) { //nolint:gosec
		potentialSize = this.Count()
	} else {
		potentialSize = uint32(len(data)) //nolint:gosec
	}

	result := NewSet3WithSize[T](potentialSize)
	for _, e := range data {
		if this.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

/*
Checks if this Set3 contains any element that is also present in that Set3. This function also provides a quick way to check if two Set3 are disjoint (i.e. !ContainsAny).

Returns false if that Set3 is nil.

Example:

	set1 := NewSet3[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := NewSet3[int]()
	set2.Add(0)
	set2.Add(1)
	b := set1.ContainsAny(set2) // b will be true
*/
func (this *Set3[T]) ContainsAny(that *Set3[T]) bool {
	if that == nil {
		return false
	}

	var smallerSet *Set3[T]
	var biggerSet *Set3[T]

	if this.Count() < that.Count() {
		smallerSet = this
		biggerSet = that
	} else {
		smallerSet = that
		biggerSet = this
	}

	for e := range smallerSet.ImmutableRange() {
		if biggerSet.Contains(e) {
			return true
		}
	}
	return false
}

/*
Checks if this Set3 contains any element fromthe given data array.

Returns false if data is nil.

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	b := set1.ContainsAnyFrom([]int{4, 5, 6}) // b will be false
*/
func (this *Set3[T]) ContainsAnyFrom(data []T) bool {
	if data == nil {
		return false
	}
	for _, d := range data {
		if this.Contains(d) {
			return true
		}
	}
	return false
}

/*
Count returns the number of elements in this Set3.

Example:

	set := NewSet3[int]()
	set.Add(7)
	set.Add(8)
	set.Add(9)
	c := set.Count()   // c will be 3
*/
func (this *Set3[T]) Count() uint32 {
	return this.resident - this.dead
}

func (this *Set3[T]) nextSize() (n uint32) {
	n = uint32(len(this.group)) * 2 //nolint:gosec
	if this.dead >= (this.resident / 2) {
		n = uint32(len(this.group)) //nolint:gosec
	}
	return
}

/*
Rorganizes the backend of this Set3 for optimal space efficiency: This call rehashes this Set3 to a size matching its current element count.

Example:

	set := NewSet3[int](1_000_000) // allocates a big hashset
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Rehash() // saves memory consumed by set
*/
func (this *Set3[T]) Rehash() {
	numGroups := uint32(calcReqNrOfGroups(this.Count()))
	this.rehashToNumGroups(numGroups)
}

/*
Rorganizes the backend of this Set3: RehashTo redistributs the elements of this Set3 onto a new hashset in its backend, e.g., to ensure faster element access.

If newSize is smaller than the current number of elements in this Set3, this function does nothing. If newSize is equal to the current number of elements in this Set3, this function does the same as [Rehash].

Example:

	set := NewSet3[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.RehashTo(1000) // ensures that you can add at least 997 more elements to set without rehashing
*/
func (this *Set3[T]) RehashTo(newSize uint32) {
	if newSize < this.Count() {
		return
	}
	newNumGroups := uint32(calcReqNrOfGroups(newSize))
	this.rehashToNumGroups(newNumGroups)
}

func (this *Set3[T]) rehashToNumGroups(newNumGroups uint32) {
	old_groups := this.fullCopyGroups()
	this.hashFunction = maphash.NewSeed(this.hashFunction)
	this.elementLimit = uint32(float64(newNumGroups) * set3maxAvgGroupLoad)
	this.resident, this.dead = 0, 0
	this.group = make([]set3Group[T], newNumGroups)
	for i := range len(this.group) {
		this.group[i].ctrl = set3AllEmpty
	}
	grpCnt := uint64(len(this.group))
	for _, oldGroup := range old_groups {
		if oldGroup.ctrl&set3hiBits != set3hiBits { // not all empty or deleted
			for s := range set3groupSize {
				if isAnElementAt(oldGroup.ctrl, s) {
					// inlined and reduced Add instead of Set3.Add(oldGroup.slot[s])
					element := oldGroup.slot[s]

					hash := this.hashFunction.Hash(element)
					H1 := (hash & 0xffff_ffff_ffff_ff80) >> 7
					H2 := (hash & 0x0000_0000_0000_007f)
					grpIdx := H1 % uint64(len(this.group))
					stillSearchingSpace := true
					for stillSearchingSpace {
						group := &this.group[grpIdx]
						ctrl := group.ctrl
						slot := &(group.slot)

						// optimization: we know it cannot exist in this Set3 already so skip
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
						grpIdx++ // carousel through all groups
						if grpIdx >= grpCnt {
							grpIdx = 0
						}
					}
				}
			}
		}
	}
}
