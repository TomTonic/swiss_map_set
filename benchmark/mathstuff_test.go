package benchmark

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedian(t *testing.T) {
	testCases := []struct {
		data     []float64
		expected float64
	}{
		{[]float64{}, 0},
		{[]float64{1}, 1},
		{[]float64{1, 2, 3}, 2},
		{[]float64{1, 2, 3, 4}, 2.5},
		{[]float64{3, 1, 2}, 2},
		{[]float64{4, 1, 3, 2}, 2.5},
		{[]float64{1, 2, 2, 3, 4}, 2},
		{[]float64{1.5, 3.5, 2.5}, 2.5},
		{[]float64{1.1, 2.2, 3.3, 4.4}, 2.75},
	}

	for _, tc := range testCases {
		result := Median(tc.data)
		assert.True(t, result == tc.expected, "FAIL: data=%v, expected=%v, got=%v\n", tc.data, tc.expected, result)
	}
}

func TestStatistics(t *testing.T) {
	testCases := []struct {
		data     []float64
		expected struct {
			mean     float64
			variance float64
			stddev   float64
		}
	}{
		{[]float64{}, struct{ mean, variance, stddev float64 }{0, -1, -1}},
		{[]float64{1}, struct{ mean, variance, stddev float64 }{1, 0, 0}},
		{[]float64{1, 2, 3}, struct{ mean, variance, stddev float64 }{2, 2 / 3.0, math.Sqrt(2 / 3.0)}},
		{[]float64{1, 2, 3, 4}, struct{ mean, variance, stddev float64 }{2.5, 1.25, math.Sqrt(1.25)}},
		{[]float64{1, 1, 1, 1}, struct{ mean, variance, stddev float64 }{1, 0, 0}},
		{[]float64{1.5, 2.5, 3.5}, struct{ mean, variance, stddev float64 }{2.5, 2 / 3.0, math.Sqrt(2 / 3.0)}},
		{[]float64{3, 53, 512, 11, 75, 201, 335}, struct{ mean, variance, stddev float64 }{170, 31576.285714285714, math.Sqrt(31576.285714285714)}},
	}

	for _, tc := range testCases {
		mean, variance, stddev := Statistics(tc.data)
		assert.True(t, mean == tc.expected.mean && variance == tc.expected.variance && stddev == tc.expected.stddev,
			"FAIL: data=%v, expected=(%v, %v, %v), got=(%v, %v, %v)\n", tc.data, tc.expected.mean, tc.expected.variance, tc.expected.stddev, mean, variance, stddev)
	}
}

func TestFloatsEqualWithTolerance(t *testing.T) {
	testCases := []struct {
		f1, f2, tolerance float64
		expected          bool
	}{
		{1.0, 1.0, 10, true},   // Exact match
		{1.0, 1.05, 10, true},  // Within tolerance
		{1.0, 1.15, 10, false}, // Outside tolerance
		{1.0, 0.95, 10, true},  // Within tolerance
		{1.0, 0.85, 10, false}, // Outside tolerance
		{1.0, 1.1, 10, true},   // On the edge of tolerance
		{1.0, 0.9, 10, true},   // On the edge of tolerance
		{2.0, 2.15, 10, true},  // Inside tolerance
		{2.0, 1.85, 10, true},  // Inside tolerance
		{2.0, 2.21, 10, true},  // Inside tolerance of second parameter
	}

	for _, tc := range testCases {
		result := FloatsEqualWithTolerance(tc.f1, tc.f2, tc.tolerance)
		assert.True(t, result == tc.expected, "%f == %f with a tolerance of %f should be %v", tc.f1, tc.f2, tc.tolerance, tc.expected)
	}
}

/*
func TestCalcNumberOfSamplesForConfidence(t *testing.T) {
	testCases := []struct {
		data     []float64
		expected int32
	}{
		{[]float64{5.0}, 0},
		{[]float64{2.0, 2.0001, 2.00005, 2.0002}, 2},
		{[]float64{0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 1.1, 1.2, 1.3, 1.4}, 385},
		{[]float64{10.0, 10.1, 10.2, 10.3, 10.4, 10.5, 10.6, 10.7, 10.8, 10.9}, 385},
		{[]float64{100.0, 100.1, 100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.9}, 385},
	}

	for _, tc := range testCases {
		result := CalcNumberOfSamplesForConfidence(tc.data)
		if result != tc.expected {
			t.Errorf("Expected %d, but got %d", tc.expected, result)
		}
	}
}
*/
