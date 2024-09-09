# Set3

Set3 is a fast and pure set implmentation in and for Golang. I wrote it as an alternative to set implementations based on `map[type]struct{}`. Set3 is a little faster and only uses about 60% of the memory compared to `map[type]struct{}`.

The code is derived from [SwissMap](https://github.com/dolthub/swiss) and it is implementing the "Fast, Efficient, Cache-friendly Hash Table" found in [Abseil](https://abseil.io/blog/20180927-swisstables). For details on the algorithm see the [CppCon 2017 talk by Matt Kulukundis](https://www.youtube.com/watch?v=ncHmEUmJZf4). The dependency on x86 assembler for [SSE2/SSE3](https://en.wikipedia.org/wiki/Streaming_SIMD_Extensions) instructions has been removed for portability and speed (yes, the code is faster without SSE).

The name "Set3" comes from the fact that this was the 3rd attempt for an optimized datastructure/code-layout to get the best runtime performance.

## Performance

The following benchmarks have been performed with [v0.1.0](https://github.com/TomTonic/Set3/releases/tag/v0.1.0) to compare `Set3[uint32]` with `map[uint32]struct{}`. Run with the command:

```
go test -benchmem -benchtime=8s -timeout 115m -run="^$" -bench "^(BenchmarkSet3Fill|BenchmarkNativeMapFill|BenchmarkSet3Find|BenchmarkNativeMapFind)$" github.com/TomTonic/Set3
```

Please note the logarithmic scales for most axes.
![Speed comparison between Set3 and built-in `map[type]struct{}`](https://github.com/user-attachments/assets/ea809575-a8b3-449a-ae35-b6aa0081487e)
![Memory comparison between Set3 and built-in `map[type]struct{}`](https://github.com/user-attachments/assets/8471cbaa-18b6-4197-8687-c03cb03ac6a9)
![Memory allocations comparison between Set3 and built-in `map[type]struct{}`](https://github.com/user-attachments/assets/34af5032-a34e-4385-b7e5-3690262ef427)
