# Set3

Set3 is a fast and pure set implmentation in and for Golang. I wrote it as an alternative to set implementations based on `map[type]struct{}`. Set3 is a little faster and only uses about 60% of the memory than using `map[type]struct{}`.

The code is derived from [SwissMap](https://github.com/dolthub/swiss) and it is implementing the "Fast, Efficient, Cache-friendly Hash Table" found in [Abseil](https://abseil.io/blog/20180927-swisstables). For details on the algorithm see the [CppCon 2017 talk by Matt Kulukundis](https://www.youtube.com/watch?v=ncHmEUmJZf4). The dependency on x86 assembler for [SSE2/SSE3](https://en.wikipedia.org/wiki/Streaming_SIMD_Extensions) instructions has been removed for portability and speed (yes, the code is faster without SSE).

The name "Set3" comes from the fact that this was the 3rd attempt for an optimized datastructure/code-layout to get the best runtime performance.

## Performance

The following benchmarks have been performed with v0.1.0 to compare `Set3[uint32]` with `map[uint32]struct{}`. Run with the command:

```
go test -benchmem -benchtime=8s -timeout 115m -run="^$" -bench "^(BenchmarkSet3Fill|BenchmarkNativeMapFill|BenchmarkSet3Find|BenchmarkNativeMapFind)$" github.com/TomTonic/Set3
```

Please note the logarithmic scales for most axes.
