package main

import (
	"reflect"
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

	result := addBenchmark(rounds, numberOfSets, initialAlloc, setSize, 0xabcdef)

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

func TestGetNumberOfSteps(t *testing.T) {
	tests := []struct {
		setSizeTo uint32
		step      Step
		expected  uint32
	}{
		{33, Step{true, true, 10.0, 0}, 11},   // expect: 0%, 10%, 20%, ..., 90%, 100%
		{19, Step{true, true, 25.0, 0}, 5},    // expect: 0%, 25%, 50%, 75%, 100%
		{19, Step{true, true, 30.0, 0}, 5},    // expect: 0%, 30%, 60%, 90%, 120%
		{234, Step{true, false, 0.0, 1}, 235}, // expect: 0, 1, 2, ..., 233, 234
		{33, Step{true, false, 0.0, 10}, 5},   // expect: 0, 10, 20, 30, 40
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := getNumberOfSteps(tt.setSizeTo, tt.step)
			if result != tt.expected {
				t.Errorf("got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestColumnHeadings(t *testing.T) {
	tests := []struct {
		setSizeTo uint32
		step      Step
		expected  []string
	}{
		{33, Step{true, true, 10.0, 0}, []string{"+0.00%% ", "+10.00%% ", "+20.00%% ", "+30.00%% ", "+40.00%% ", "+50.00%% ", "+60.00%% ", "+70.00%% ", "+80.00%% ", "+90.00%% ", "+100.00%% "}},
		{19, Step{true, true, 25.0, 0}, []string{"+0.00%% ", "+25.00%% ", "+50.00%% ", "+75.00%% ", "+100.00%% "}},
		{19, Step{true, true, 30.0, 0}, []string{"+0.00%% ", "+30.00%% ", "+60.00%% ", "+90.00%% ", "+120.00%% "}},
		{4, Step{true, false, 0.0, 1}, []string{"+0 ", "+1 ", "+2 ", "+3 ", "+4 "}},
		{33, Step{true, false, 0.0, 10}, []string{"+0 ", "+10 ", "+20 ", "+30 ", "+40 "}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := columnHeadings(tt.setSizeTo, tt.step)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInitSizeValues(t *testing.T) {
	tests := []struct {
		currentSetSize uint32
		setSizeTo      uint32
		step           Step
		expected       []uint32
	}{
		{10, 11, Step{true, true, 10.0, 0}, []uint32{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		{100, 1000, Step{true, true, 25.0, 0}, []uint32{100, 125, 150, 175, 200}},
		{10, 1, Step{true, true, 30.0, 0}, []uint32{10, 13, 16, 19, 22}},
		{2, 6, Step{true, false, 0.0, 1}, []uint32{2, 3, 4, 5, 6, 7, 8}},
		{33, 40, Step{true, false, 0.0, 10}, []uint32{33, 43, 53, 63, 73}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := initSizeValues(tt.currentSetSize, tt.setSizeTo, tt.step)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStep_Set(t *testing.T) {
	tests := []struct {
		input    string
		expected Step
		err      bool
	}{
		{"10%", Step{true, true, 10.0, 0}, false},
		{"2.5%", Step{true, true, 2.5, 0}, false},
		{"5", Step{true, false, 0.0, 5}, false},
		{"0", Step{true, false, 0.0, 0}, false},
		{"invalid%", Step{}, true},
		{"invalid", Step{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var step Step
			err := step.Set(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("Set() error = %v, expected error = %v", err, tt.err)
			}
			if err == nil && step != tt.expected {
				t.Errorf("Set() = %v, expected %v", step, tt.expected)
			}
		})
	}
}

func TestStep_String(t *testing.T) {
	tests := []struct {
		step     Step
		expected string
	}{
		{Step{true, true, 10.0, 0}, "10.000000%"},
		{Step{true, true, 25.0, 0}, "25.000000%"},
		{Step{true, false, 0.0, 5}, "5"},
		{Step{true, false, 0.0, 0}, "0"},
		{Step{false, false, 0.0, 0}, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.step.String(); got != tt.expected {
				t.Errorf("String() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestToNSperAdd(t *testing.T) {
	tests := []struct {
		measurements []time.Duration
		addsPerRound uint32
		expected     []float64
	}{
		{[]time.Duration{time.Nanosecond * 10, time.Nanosecond * 20}, 2, []float64{5, 10}},
		{[]time.Duration{time.Nanosecond * 100, time.Nanosecond * 200}, 4, []float64{25, 50}},
		{[]time.Duration{time.Nanosecond * 0, time.Nanosecond * 50}, 5, []float64{0, 10}},
		{[]time.Duration{time.Nanosecond * 1000}, 10, []float64{100}},
		{[]time.Duration{}, 1, []float64{}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := toNSperAdd(tt.measurements, tt.addsPerRound)
			if len(result) != len(tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("at index %d, got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}
