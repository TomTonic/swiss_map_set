# Set3

Set3 is a fast and pure set implmentation in and for Golang.

The code is derived from [SwissMap](https://github.com/dolthub/swiss) and it is implementing the "Fast, Efficient, Cache-friendly Hash Table" found in [Abseil](https://abseil.io/blog/20180927-swisstables). For details on the algorithm see the [CppCon 2017 talk by Matt Kulukundis](https://www.youtube.com/watch?v=ncHmEUmJZf4). The dependency on x86 assembler for [SSE](https://en.wikipedia.org/wiki/Streaming_SIMD_Extensions) instructions has been removed for portability and speed (yes, the code is faster without SSE).
