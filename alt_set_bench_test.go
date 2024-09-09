package swiss

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
)

type testMapType map[uint32]struct{}

/*
	func prepareDataUint32(initialSetSize, finalSetSize, searchListSize int, minimalHitRatio float32) (resultSet *Set[uint32], resultMap testMapType, searchElements []uint32) {
		resultSet = NewSet[uint32](uint32(initialSetSize))
		resultMap = make(testMapType, initialSetSize)
		for n := 0; n < finalSetSize; n++ {
			element := rand.Uint32()
			resultSet.Add(element)
			resultMap[element] = 1
		}
		nrOfElemToCopy := int(minimalHitRatio * float32(searchListSize))
		tempList := make([]uint32, 0, searchListSize)
		countCopied := 0
		for countCopied < nrOfElemToCopy {
			resultSet.Iter(func(e uint32) (stop bool) {
				tempList = append(tempList, e)
				countCopied++
				return countCopied >= nrOfElemToCopy
			})
		}
		for n := countCopied; n < searchListSize; n++ {
			element := rand.Uint32()
			tempList = append(tempList, element)
		}
		perm := rand.Perm(searchListSize)
		searchElements = make([]uint32, searchListSize)
		for i, idx := range perm {
			searchElements[i] = tempList[idx]
		}
		return
	}
*/
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
	{inintialSetSize: 10, finalSetSize: 10, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 20, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 30, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 40, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 50, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
}

func BenchmarkSet1Fill(b *testing.B) {
	for _, cfg := range config {
		setValues, _ := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := NewSet[uint32](uint32(cfg.inintialSetSize))
				for j := 0; j < len(setValues); j++ {
					resultSet.Add(setValues[j])
				}
			}
		})
	}
}

func BenchmarkSet2Fill(b *testing.B) {
	for _, cfg := range config {
		setValues, _ := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		b.Run(fmt.Sprintf("inintial(%d);final(%d);search(%d);hit(%f)", cfg.inintialSetSize, cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSet := NewSet2[uint32](uint32(cfg.inintialSetSize))
				for j := 0; j < len(setValues); j++ {
					resultSet.Add(setValues[j])
				}
			}
		})
	}
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

func BenchmarkSet1Find(b *testing.B) {
	for _, cfg := range config {
		setValues, searchElements := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		resultSet := NewSet[uint32](uint32(cfg.inintialSetSize))
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

func BenchmarkSet2Find(b *testing.B) {
	for _, cfg := range config {
		setValues, searchElements := prepareDataUint32(cfg.finalSetSize, cfg.searchListSize, cfg.minimalHitRatio, cfg.seed)
		resultSet := NewSet2[uint32](uint32(cfg.inintialSetSize))
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
