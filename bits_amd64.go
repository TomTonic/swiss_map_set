// Copyright 2023 Dolthub, Inc.
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

//go:build amd64 && !nosimd

package swiss

import (
	"math/bits"
	_ "unsafe"

	"github.com/dolthub/swiss/simd"
)

const (
	groupSize       = 16
	maxAvgGroupLoad = 14
)

func ctlrMatchH2(m *[16]int8, h uint64) uint64 {
	b := simd.MatchCRTLhash(m, h)
	return b
}

func ctlrMatchEmpty(m *[16]int8) uint64 {
	b := simd.MatchCRTLempty(m)
	return b
}

func nextMatch(b *uint64) (s int) {
	s = bits.TrailingZeros16(uint16(*b))
	*b &= ^(1 << s) // clear bit |s|
	return
}
