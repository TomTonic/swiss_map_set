package benchmark

import "sort"

func Median(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)
	sort.Float64s(dataCopy)

	l := len(dataCopy)
	if l%2 == 0 {
		return (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	}
	return dataCopy[l/2]
}
