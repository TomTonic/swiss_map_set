# Set3

[![Go Report Card](https://goreportcard.com/badge/github.com/TomTonic/Set3)](https://goreportcard.com/report/github.com/TomTonic/Set3)
[![Go Reference](https://pkg.go.dev/badge/github.com/TomTonic/Set3.svg)](https://pkg.go.dev/github.com/TomTonic/Set3)
[![Tests](https://github.com/TomTonic/Set3/actions/workflows/coverage.yml/badge.svg?branch=main)](https://github.com/TomTonic/Set3/actions/workflows/coverage.yml)
![coverage](https://raw.githubusercontent.com/TomTonic/Set3/badges/.badges/main/coverage.svg)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9470/badge)](https://www.bestpractices.dev/projects/9470)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/TomTonic/Set3/badge)](https://scorecard.dev/viewer/?uri=github.com/TomTonic/Set3)

Set3 is a fast and pure set implmentation in and for Golang. I wrote it as an alternative to set implementations based on `map[type]struct{}`. Set3 is 10%-20% faster and uses 40% less memory than `map[type]struct{}`. As hash function, Set3 uses the built-in hash function of Golang via [dolthub/maphash](https://github.com/dolthub/maphash).

The code is derived from [SwissMap](https://github.com/dolthub/swiss) and it implements the "Fast, Efficient, Cache-friendly Hash Table" found in [Abseil](https://abseil.io/blog/20180927-swisstables).
For details on the algorithm see the [CppCon 2017 talk by Matt Kulukundis](https://www.youtube.com/watch?v=ncHmEUmJZf4).
The dependency on x86 assembler for [SSE2/SSE3](https://en.wikipedia.org/wiki/Streaming_SIMD_Extensions) instructions has been removed for portability and speed; the code runs faster without SSE and the necessary additional stack frame.

The name "Set3" comes from the fact that this was the 3rd attempt for an optimized datastructure/code-layout to get the best runtime performance.

## Installation

To use the `Set3` package in your Go project, follow these steps:

1. **Initialize a Go module** (if you haven't already):

   ```sh
   go mod init your-module-name
   ```

2. **Add the package**: Simply import the package in your Go code, and Go modules will handle the rest:

   ```go
   import "github.com/TomTonic/Set3"
   ```

3. **Download dependencies**: Run the following command to download the dependencies:

   ```sh
   go mod tidy
   ```

   This will automatically download and install the Set3 package along with any other dependencies.

## Using Set3

The following test case creates two sets and demonstrates some operations. For a full list of operations on the `Set3` type, see [API doc](https://pkg.go.dev/github.com/TomTonic/Set3#Set3).

```go
func TestExample(t *testing.T) {
    // create a new Set3
    set1 := NewSet3[int]()
    // add some elements
    set1.Add(1)
    set1.Add(2)
    set1.Add(3)
    // add some more elements from an array
    set1.AddAllFrom([]int{4, 5, 6})
    // create a second set directly from an array
    set2 := AsSet3([]int{2, 3, 4, 5})
    // check if set2 is a subset of set1. must be true in this case
    isSubset := set1.ContainsAll(set2)
    assert.True(t, isSubset, "%v is not a subset of %v", set2, set1)
    // mathematical operations like Union, Difference and Intersect
    // do not manipulate a Set3 but return a new set
    intersect := set1.Intersection(set2)
    // compare sets. as set2 is a subset of set1, intersect must be equal to set2
    equal := intersect.Equals(set2)
    assert.True(t, equal, "%v is not equal to %v", intersect, set2)
}
```

## Performance

The following benchmarks have been performed with [v0.2.0](https://github.com/TomTonic/Set3/releases/tag/v0.2.0) to compare `Set3[uint32]` with `map[uint32]struct{}` with the command:

```sh
go test -benchmem -benchtime=6s -timeout 480m -run="^$" -bench "^(BenchmarkSet3Fill|BenchmarkNativeMapFill|BenchmarkSet3Find|BenchmarkNativeMapFind)$" github.com/TomTonic/Set3 > benchresult.txt
```

(Raw benchmark results are available [here](https://raw.githubusercontent.com/TomTonic/Set3/main/benchresult.txt). Go version 1.23.1, no PGO.)

### Inserting Nodes into an Empty Set

The total time for inserting random elements into newly allocated sets with an initial capacity of 21 elements; i.e. rehashing took place multiple times for larger sets.

n = 1 ... 300 (step size +1, linear scale)
![Inserting Nodes into an Empty Set, n = 1 ... 300 (step size +1, linear scale)](https://github.com/user-attachments/assets/3dbf7f75-8859-46da-9512-c61e151db2fd)

n = 50 ... 6,193,533 (step size +5%, log scale)
![Inserting Nodes into an Empty Set, n = 50 ... 6,193,533 (step size +5%, log scale)](https://github.com/user-attachments/assets/17bfa1d8-403f-460f-9f15-1c6394ea0c9e)

### Searching Nodes in a Populated Set

Total time for searching 5000 elements in sets of different sizes, 30% hit ratio.

n = 1 ... 300 (step size +1, linear scale)
![Searching Nodes in a Populated Set, n = 1 ... 300 (step size +1, linear scale)](https://github.com/user-attachments/assets/52ba00e2-32e8-41f4-ae6c-a2a237990fb1)

n = 50 ... 6,193,533 (step size +5%, linear scale)
![Searching Nodes in a Populated Set, n = 50 ... 6,193,533 (step size +5%, linear scale)](https://github.com/user-attachments/assets/0aec23ec-52b7-45c1-ae75-f2397651bd85)
