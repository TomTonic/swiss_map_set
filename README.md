# Set3

[![Tests](https://github.com/TomTonic/Set3/actions/workflows/coverage.yml/badge.svg?branch=main)](https://github.com/TomTonic/Set3/actions/workflows/coverage.yml)
![coverage](https://raw.githubusercontent.com/TomTonic/Set3/badges/.badges/main/coverage.svg)

Set3 is a fast and pure set implmentation in and for Golang. I wrote it as an alternative to set implementations based on `map[type]struct{}`. Set3 is 10%-20% faster and uses 40% less memory than `map[type]struct{}`. As hash function, Set3 uses the built-in hash function of Golang via [dolthub/maphash](https://github.com/dolthub/maphash).

The code is derived from [SwissMap](https://github.com/dolthub/swiss) and it implements the "Fast, Efficient, Cache-friendly Hash Table" found in [Abseil](https://abseil.io/blog/20180927-swisstables). For details on the algorithm see the [CppCon 2017 talk by Matt Kulukundis](https://www.youtube.com/watch?v=ncHmEUmJZf4). The dependency on x86 assembler for [SSE2/SSE3](https://en.wikipedia.org/wiki/Streaming_SIMD_Extensions) instructions has been removed for portability and speed (yes, the code runs faster without SSE and the necessary additional stack frame).

The name "Set3" comes from the fact that this was the 3rd attempt for an optimized datastructure/code-layout to get the best runtime performance.

## Performance

The following benchmarks have been performed with [v0.2.0](https://github.com/TomTonic/Set3/releases/tag/v0.2.0) to compare `Set3[uint32]` with `map[uint32]struct{}` with the command:

```
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
