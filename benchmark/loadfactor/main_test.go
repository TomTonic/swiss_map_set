package main

import (
	"runtime/debug"
	"testing"
	"time"

	"github.com/TomTonic/Set3/benchmark"
	"github.com/loov/hrtime"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	reportedPrecision := hrtime.NowPrecision()
	var values [20_000_001]time.Duration
	debug.SetGCPercent(-1)
	for i := range 20_000_001 {
		values[i] = hrtime.Now()
	}
	debug.SetGCPercent(100)
	deltas := make([]float64, 0, 20_000_000)
	var zeros uint64
	for i := range 20_000_000 {
		di := values[i+1] - values[i]
		if di == 0 {
			zeros++
		} else {
			deltas = append(deltas, float64(di))
		}
	}
	median := benchmark.Median(deltas)
	assert.True(t, benchmark.FloatsEqualWithTolerance(reportedPrecision, median, 0.001), "reportedPrecision and median differ (%f!=%f @%f%% tolerance)", reportedPrecision, median, 0.001)
	// hist := hrtime.NewHistogram(deltas, &defaultOptions)
	// fmt.Printf(hist.String())
}

func TestDoBenchmark2(t *testing.T) {
	rounds := 72
	numberOfSets := uint32(10)
	initialAlloc := uint32(150)
	setSize := uint32(100)

	result := doBenchmark2(rounds, numberOfSets, initialAlloc, setSize, 0xabcdef)

	assert.True(t, len(result) == rounds, "Result should return %d measurements. It returned %d measurements.", rounds, len(result))
	assert.False(t, containsZero(result), "Result should not contain zeros, but it does.")
	assert.False(t, containsNegative(result), "Result should not contain negative numbers, but it does.")

	reportedPrecision := hrtime.NowPrecision()
	assert.True(t, atLeastNtimesPrecision(20.0, reportedPrecision, result),
		"Result should only contain values that exceed %fx the timer precision of %fns, but it does not. The minimum Value is %v.", 20.0, reportedPrecision, minVal(result))
}

func containsZero(measurements []time.Duration) bool {
	for _, d := range measurements {
		if d == 0 {
			return true
		}
	}
	return false
}

func containsNegative(measurements []time.Duration) bool {
	for _, d := range measurements {
		if d < 0 {
			return true
		}
	}
	return false
}

func atLeastNtimesPrecision(nTimes float64, precision float64, measurements []time.Duration) bool {
	for _, d := range measurements {
		if float64(d) < precision*nTimes {
			return false
		}
	}
	return true
}

func minVal(measurements []time.Duration) time.Duration {
	min := 48 * time.Hour
	for _, d := range measurements {
		if d < min {
			min = d
		}
	}
	return min
}
