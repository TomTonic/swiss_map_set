# Set3

[![Tests](https://github.com/TomTonic/Set3/actions/workflows/coverage.yml/badge.svg?branch=main)](https://github.com/TomTonic/Set3/actions/workflows/coverage.yml)
![coverage](https://raw.githubusercontent.com/TomTonic/Set3/badges/.badges/main/coverage.svg)

Set3 is a fast and pure set implmentation in and for Golang. I wrote it as an alternative to set implementations based on `map[type]struct{}`. Set3 is 15% faster and uses 40% less memory than `map[type]struct{}`. As hash function, Set3 uses the built-in hash function of Golang via [dolthub/maphash](https://github.com/dolthub/maphash).

The code is derived from [SwissMap](https://github.com/dolthub/swiss) and it is implementing the "Fast, Efficient, Cache-friendly Hash Table" found in [Abseil](https://abseil.io/blog/20180927-swisstables). For details on the algorithm see the [CppCon 2017 talk by Matt Kulukundis](https://www.youtube.com/watch?v=ncHmEUmJZf4). The dependency on x86 assembler for [SSE2/SSE3](https://en.wikipedia.org/wiki/Streaming_SIMD_Extensions) instructions has been removed for portability and speed (yes, the code is faster without SSE).

The name "Set3" comes from the fact that this was the 3rd attempt for an optimized datastructure/code-layout to get the best runtime performance.

## Performance

The following benchmarks have been performed with [v0.1.0](https://github.com/TomTonic/Set3/releases/tag/v0.1.0) to compare `Set3[uint32]` with `map[uint32]struct{}`. Run with the command:

```
go test -benchmem -benchtime=8s -timeout 115m -run="^$" -bench "^(BenchmarkSet3Fill|BenchmarkNativeMapFill|BenchmarkSet3Find|BenchmarkNativeMapFind)$" github.com/TomTonic/Set3
```

Please note the logarithmic scales for most axes.

### Inserting Nodes in an Empty Set

Total time/additional memory/memory allocations for inserting different numbers of elements into empty sets with an initial capacity of 10 elements.

![Speed comparison between Set3 and built-in `map[type]struct{}`](https://github.com/user-attachments/assets/46e9c4d1-45b7-4487-b2d3-a220150c5cdc)
![Memory comparison between Set3 and built-in `map[type]struct{}`](https://github.com/user-attachments/assets/8471cbaa-18b6-4197-8687-c03cb03ac6a9)
![Memory allocations comparison between Set3 and built-in `map[type]struct{}`](https://github.com/user-attachments/assets/34af5032-a34e-4385-b7e5-3690262ef427)

### Searching Nodes in a Populated Set

Total time for searching 5000 elements in sets of different sizes, 30% hit ratio.
Tests performed with commit id ceed986615c10f0b1493b855f703f880c1dedc9b, reduced maxAvgGroupLoad to 6.5, i.e. same as `map[type]struct{}`. Mean size ratio drops to 0.608.
![Screenshot 2024-09-11 123416](https://github.com/user-attachments/assets/03df5b2d-4165-4a2c-b1cb-c42f8f923e45)
