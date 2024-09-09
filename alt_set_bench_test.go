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
	{inintialSetSize: 10, finalSetSize: 10, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 13, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 18, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 24, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 32, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 42, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 56, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 75, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 100, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 133, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 177, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 236, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 315, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 420, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 559, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 746, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 994, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1325, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 1766, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 2354, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 3138, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 4182, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 5575, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 7432, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 9907, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 13205, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 17603, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 23465, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 31278, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 41694, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 55578, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 74086, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 98756, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 131642, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 175479, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 233913, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 311807, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 415638, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 554046, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
	{inintialSetSize: 10, finalSetSize: 738543, searchListSize: 5000, minimalHitRatio: 0.3, seed: 0x1234567890ABCDEF},
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
