package swiss

import (
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
)

func prepareDataUint32(initialSetSize, finalSetSize, searchListSize int, minimalHitRatio float32) (resultSet *Set[uint32], searchElements []uint32) {
	resultSet = NewSet[uint32](uint32(initialSetSize))
	for n := 0; n < finalSetSize; n++ {
		element := rand.Uint32()
		resultSet.Add(element)
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

func Algorithm1(resultSet *Set[uint32], searchElements []uint32) {
	x := uint64(0)
	for _, e := range searchElements {
		if resultSet.Has(e) {
			x += uint64(e)
		}
	}
}

func Algorithm2(resultSet *Set[uint32], searchElements []uint32) {
	x := uint64(0)
	for _, e := range searchElements {
		if resultSet.Contains(e) {
			x += uint64(e)
		}
	}
}

func Algorithm3(resultSet *Set[uint32], searchElements []uint32) {
	x := uint64(0)
	for _, e := range searchElements {
		if resultSet.Contains2(e) {
			x += uint64(e)
		}
	}
}

func BenchmarkAlgorithm1(b *testing.B) {
	resultSet, searchElements := prepareDataUint32(10, 5000, 50000, 0.99)
	for i := 0; i < b.N; i++ {
		Algorithm1(resultSet, searchElements)
	}
}

func BenchmarkAlgorithm2(b *testing.B) {
	resultSet, searchElements := prepareDataUint32(10, 5000, 50000, 0.99)
	for i := 0; i < b.N; i++ {
		Algorithm2(resultSet, searchElements)
	}
}

func BenchmarkAlgorithm3(b *testing.B) {
	resultSet, searchElements := prepareDataUint32(10, 5000, 50000, 0.99)
	for i := 0; i < b.N; i++ {
		Algorithm3(resultSet, searchElements)
	}
}

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	resultSet, searchElements := prepareDataUint32(10, 20, 30, 0.3)
	println("Set: %v", resultSet)
	println("Search: %v", searchElements)

	Algorithm1(resultSet, searchElements) // oder Algorithm2()
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
