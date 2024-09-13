/*
Set3 is an efficient set implementation in plain Go. Unlike many other set implementations, Set3 does not rely on Go's internal map data structure.
Instead, it implements a hash set based on the "Fast, Efficient, Cache-friendly Hash Table" found in Abseil, Google's C++ libraries.
As a result, lookups are 15% faster and the data structure uses 40% less memory than implementations based on `map[type]struct{}`.
*/
package Set3
