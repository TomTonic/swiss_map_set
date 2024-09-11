package Set3

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
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
	{inintialSetSize: 10, finalSetSize: 20, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 21, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 22, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 23, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 24, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 26, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 27, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 28, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 30, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 31, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 33, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 34, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 36, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 38, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 40, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 42, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 44, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 46, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 48, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 51, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 53, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 56, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 59, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 61, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 65, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 68, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 71, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 75, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 78, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 82, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 86, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 91, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 95, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 100, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 105, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 110, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 116, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 122, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 128, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 134, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 141, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 148, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 155, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 163, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 171, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 180, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 189, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 198, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 208, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 218, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 229, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 241, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 253, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 265, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 279, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 293, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 307, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 323, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 339, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 356, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 374, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 392, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 412, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 432, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 454, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 477, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 501, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 526, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 552, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 580, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 609, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 639, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 671, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 704, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 740, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 777, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 815, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 856, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 899, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 944, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 991, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1041, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1093, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1147, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1205, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1265, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1328, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1395, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1464, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1538, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1615, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1695, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1780, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1869, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1963, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2061, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2164, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2272, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2386, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2505, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2630, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2762, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2900, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3045, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3197, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3357, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3524, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3701, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3886, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 4080, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 4284, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 4498, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 4723, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 4959, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 5207, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 5468, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 5741, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 6028, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 6329, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 6646, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 6978, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 7327, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 7694, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 8078, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 8482, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 8906, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 9352, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 9819, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 10310, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 10826, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 11367, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 11935, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 12532, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 13159, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 13816, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 14507, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 15233, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 15994, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 16794, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 17634, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 18515, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 19441, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 20413, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 21434, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 22506, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 23631, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 24812, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 26053, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 27356, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 28723, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 30160, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 31668, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 33251, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 34913, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 36659, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 38492, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 40417, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 42438, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 44559, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 46787, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 49127, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 51583, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 54162, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 56870, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 59714, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 62700, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 65835, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 69126, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 72583, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 76212, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 80022, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 84023, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 88225, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 92636, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 97268, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 102131, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 107237, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 112599, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 118229, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 124141, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 130348, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 136865, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 143708, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 150894, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 158439, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 166361, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 174679, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 183412, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 192583, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 202212, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 212323, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 222939, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 234086, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 245790, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 258080, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 270984, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 284533, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 298760, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 313698, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 329382, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
}

func BenchmarkSet3Fill(b *testing.B) {
	for _, cfg := range config {
		setValues, _ := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := NewSet3[uint32](uint32(cfg.inintialSetSize))
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
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := make(testMapType, 10)
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
		resultSet := NewSet3[uint32](uint32(cfg.inintialSetSize))
		for j := 0; j < len(setValues); j++ {
			resultSet.Add(setValues[j])
		}
		var x uint64
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x = 0
				for _, e := range searchElements {
					if resultSet.Contains(e) {
						x += 1
					}
				}
			}
		})
		//println(x)
	}
}

func BenchmarkNativeMapFind(b *testing.B) {

	for _, cfg := range config {
		setValues, searchElements := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		resultSet := make(testMapType, cfg.inintialSetSize)
		for j := 0; j < len(setValues); j++ {
			resultSet[setValues[j]] = struct{}{}
		}
		var x uint64
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				x = 0
				for _, e := range searchElements {
					_, b := resultSet[e]
					if b {
						x += 1
					}
				}
			}
		})
		//println(x)
	}
}

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	/*
		resultSet, resultMap, searchElements := prepareDataUint32(10, 20, 30, 0.3)
		println("Set: %v", resultSet)
		println("Map: %v", resultMap)
		println("Search: %v", searchElements)

		Algorithm1(resultSet, resultMap, searchElements) // oder Algorithm2()
	*/
}

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
