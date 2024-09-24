// Copyright 2024 TomTonic
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Set3 is an efficient set implementation in plain Go. Unlike many other set implementations, Set3 does not rely on Go's internal map data structure.
Instead, it implements a hash set based on the "Fast, Efficient, Cache-friendly Hash Table" found in Abseil, Google's C++ libraries.
As a result, Set3 is 10%-20% faster and the data structure uses 40% less memory than implementations based on `map[type]struct{}`.
*/
package set3
