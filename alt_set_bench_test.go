package set3

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

type testMapType map[uint32]struct{}

func prepareDataUint32(setSize, searchListSize int, minimalHitRatio float32, seed int64) (setValues []uint32, searchElements []uint32) {
	rng := rand.New(rand.NewSource(seed))
	setValues = make([]uint32, setSize)
	for n := 0; n < setSize; n++ {
		element := rng.Uint32()
		setValues = append(setValues, element)
	}
	nrOfElemToCopy := int(minimalHitRatio * float32(searchListSize))
	tempList := make([]uint32, 0, searchListSize)
	countCopied := 0
	for countCopied < nrOfElemToCopy {
		for _, e := range setValues {
			tempList = append(tempList, e)
			countCopied++
			if countCopied >= nrOfElemToCopy {
				break
			}
		}
	}
	for n := countCopied; n < searchListSize; n++ {
		element := rng.Uint32()
		tempList = append(tempList, element)
	}
	perm := rng.Perm(searchListSize)
	searchElements = make([]uint32, searchListSize)
	for i, idx := range perm {
		searchElements[i] = tempList[idx]
	}
	return
}

var config = []struct {
	inintialSetSize int
	finalSetSize    int
	searchListSize  int
	minimalHitRatio float32
	seed            int64
}{
	{inintialSetSize: 21, finalSetSize: 1, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 2, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 3, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 4, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 5, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 6, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 7, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 8, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 9, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 10, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 11, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 12, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 13, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 14, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 15, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 16, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 17, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 18, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 19, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 20, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 21, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 22, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 23, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 24, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 25, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 26, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 27, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 28, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 29, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 30, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 31, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 32, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 33, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 34, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 35, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 36, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 37, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 38, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 39, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 40, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 41, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 42, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 43, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 44, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 45, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 46, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 47, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 48, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 49, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 50, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 51, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 52, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 53, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 54, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 55, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 56, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 57, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 58, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 59, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 60, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 61, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 62, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 63, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 64, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 65, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 66, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 67, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 68, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 69, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 70, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 71, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 72, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 73, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 74, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 75, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 76, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 77, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 78, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 79, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 80, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 81, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 82, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 83, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 84, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 85, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 86, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 87, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 88, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 89, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 90, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 91, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 92, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 93, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 94, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 95, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 96, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 97, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 98, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 99, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 100, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 101, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 102, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 103, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 104, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 105, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 106, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 107, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 108, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 109, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 110, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 111, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 112, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 113, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 114, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 115, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 116, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 117, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 118, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 119, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 120, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 121, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 122, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 123, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 124, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 125, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 126, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 127, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 128, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 129, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 130, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 131, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 132, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 133, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 134, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 135, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 136, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 137, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 138, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 139, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 140, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 141, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 142, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 143, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 144, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 145, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 146, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 147, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 148, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 149, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 150, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 151, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 152, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 153, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 154, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 155, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 156, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 157, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 158, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 159, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 160, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 161, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 162, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 163, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 164, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 165, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 166, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 167, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 168, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 169, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 170, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 171, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 172, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 173, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 174, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 175, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 176, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 177, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 178, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 179, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 180, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 181, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 182, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 183, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 184, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 185, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 186, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 187, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 188, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 189, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 190, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 191, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 192, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 193, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 194, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 195, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 196, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 197, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 198, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 199, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 200, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 201, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 202, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 203, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 204, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 205, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 206, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 207, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 208, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 209, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 210, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 211, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 212, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 213, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 214, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 215, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 216, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 217, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 218, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 219, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 220, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 221, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 222, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 223, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 224, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 225, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 226, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 227, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 228, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 229, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 230, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 231, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 232, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 233, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 234, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 235, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 236, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 237, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 238, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 239, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 240, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 241, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 242, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 243, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 244, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 245, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 246, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 247, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 248, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 249, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 250, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 251, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 252, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 253, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 254, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 255, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 256, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 257, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 258, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 259, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 260, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 261, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 262, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 263, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 264, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 265, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 266, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 267, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 268, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 269, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 270, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 271, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 272, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 273, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 274, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 275, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 276, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 277, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 278, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 279, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 280, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 281, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 282, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 283, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 284, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 285, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 286, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 287, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 288, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 289, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 290, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 291, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 292, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 293, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 294, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 295, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 296, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 297, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 298, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 299, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 21, finalSetSize: 300, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
}

func BenchmarkSet3Fill(b *testing.B) {
	for _, cfg := range config {
		setValues, _ := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := EmptyWithCapacity[uint32](uint32(cfg.inintialSetSize))
				for j := 0; j < len(setValues); j++ {
					resultSet.Add(setValues[j])
				}
			}
		})
	}
}

func BenchmarkNativeMapFill(b *testing.B) {
	for _, cfg := range config {
		setValues, _ := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := make(testMapType, cfg.inintialSetSize)
				for j := 0; j < len(setValues); j++ {
					resultSet[setValues[j]] = struct{}{}
				}
			}
		})
	}
}

func BenchmarkSet3Find(b *testing.B) {
	for _, cfg := range config {
		setValues, searchElements := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		resultSet := EmptyWithCapacity[uint32](uint32(cfg.inintialSetSize))
		for j := 0; j < len(setValues); j++ {
			resultSet.Add(setValues[j])
		}
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		var x uint64
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x = 0
				for _, e := range searchElements {
					if resultSet.Contains(e) {
						x++
					}
				}
			}
		})
		// println(x)
	}
}

func BenchmarkNativeMapFind(b *testing.B) {
	for _, cfg := range config {
		setValues, searchElements := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		resultSet := make(testMapType, cfg.inintialSetSize)
		for j := 0; j < len(setValues); j++ {
			resultSet[setValues[j]] = struct{}{}
		}
		// Force garbage collection
		runtime.GC()
		// Give the garbage collector some time to complete
		time.Sleep(1 * time.Second)
		var x uint64
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x = 0
				for _, e := range searchElements {
					_, b := resultSet[e]
					if b {
						x++
					}
				}
			}
		})
		// println(x)
	}
}

/*
func myGenStringData(size, count int) (result []string) {
	src := rand.New(rand.NewSource(int64(size * count)))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := make([]rune, size*count)
	for i := range r {
		r[i] = letters[src.Intn(len(letters))]
	}
	result = make([]string, count)
	for i := range result {
		result[i] = string(r[:size])
		r = r[size:]
	}
	return
}

func myGenUint32Data(count int) (result []uint32) {
	result = make([]uint32, count)
	for i := range result {
		result[i] = rand.Uint32()
	}
	return
}

func myGenerateInt64Data(n int) (data []int64) {
	data = make([]int64, n)
	for i := range data {
		data[i] = rand.Int63()
	}
	return
}
*/
