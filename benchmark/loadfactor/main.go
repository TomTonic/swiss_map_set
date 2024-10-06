package main

import (
	"fmt"
	"math"
	"runtime"
	"runtime/debug"
	"time"

	set3 "github.com/TomTonic/Set3"
	"github.com/TomTonic/Set3/benchmark"
	"github.com/TomTonic/objectsize"
	"github.com/loov/hrtime"
)

func doBenchmark(rounds, numberOfSets, initialAlloc, setSize uint32, seed uint64) {
	fmt.Printf("%d;%d;%d;%d;", rounds, numberOfSets, initialAlloc, setSize)
	set := make([]*set3.Set3[uint64], numberOfSets)
	timePerRound := make([]float64, rounds)
	prng := benchmark.PRNG{State: seed}
	var startMem, endMem runtime.MemStats
	var startTime, endTime time.Duration
	debug.SetGCPercent(-1)
	runtime.GC()
	runtime.ReadMemStats(&startMem)
	for i := range numberOfSets {
		set[i] = set3.EmptyWithCapacity[uint64](uint32(initialAlloc))
	}
	runtime.GC()
	runtime.ReadMemStats(&endMem)
	//reqMem := float64((endMem.HeapAlloc+endMem.StackInuse+endMem.StackSys)-(startMem.HeapAlloc+startMem.StackInuse+startMem.StackSys)) / float64(numberOfSets) / float64(setSize)
	val, _ := objectsize.Of(set)
	//reqMem2 := float64(val) / float64(numberOfSets*setSize)
	fmt.Printf("%d;", val)
	runtime.GC()
	for r := range rounds {
		for s := range numberOfSets {
			set[s].Clear()
		}
		startTime = hrtime.Now()
		for s := range numberOfSets {
			currentSet := set[s]
			for range setSize {
				currentSet.Add(prng.Uint64())
			}
		}
		endTime = hrtime.Now()
		timePerRound[r] = float64(endTime - startTime)
	}
	debug.SetGCPercent(100)
	med := benchmark.Median(timePerRound)
	avg, variance, stddev := benchmark.Statistics(timePerRound)
	fmt.Printf("%.3f;%.3f;%.3f;%.3f\n", med, avg, variance, stddev)
	hist := hrtime.NewHistogram(timePerRound, &defaultOptions)
	fmt.Printf(hist.String())
}

func doBenchmark2(rounds int, numberOfSets, initialAlloc, setSize uint32, seed uint64) (measurements []time.Duration) {
	prng := benchmark.PRNG{State: seed}
	set := make([]*set3.Set3[uint64], numberOfSets)
	for i := range numberOfSets {
		set[i] = set3.EmptyWithCapacity[uint64](initialAlloc)
	}
	timePerRound := make([]time.Duration, rounds)
	runtime.GC()
	debug.SetGCPercent(-1)
	defer debug.SetGCPercent(100)
	for r := range rounds {
		for s := range numberOfSets {
			set[s].Clear()
		}
		startTime := hrtime.Now()
		for s := range numberOfSets {
			currentSet := set[s]
			for range setSize {
				currentSet.Add(prng.Uint64())
			}
		}
		endTime := hrtime.Now()
		timePerRound[r] = endTime - startTime
	}
	return timePerRound
}

var defaultOptions = hrtime.HistogramOptions{
	BinCount:        10,
	NiceRange:       true,
	ClampMaximum:    0,
	ClampPercentile: 0.95,
}

func toNSperAdd(measurements []time.Duration, addsPerRound uint32) []float64 {
	result := make([]float64, len(measurements))
	div := 1.0 / float64(addsPerRound)
	for i, m := range measurements {
		result[i] = float64(m) * div
	}
	return result
}

func main() {
	/*
		fmt.Printf("rounds;numberOfSets;initialAlloc;setSize;memForAllSets;medianNSperRound;avgPerRound;varPerRound;stddevPerRound\n")
		for i := range 200 {
			doBenchmark(20_001, 100, 200+uint32(i), 200, 0xABCDEF0123456789)
		}
	*/
	fmt.Printf("setSize ")
	for initSize := range 101 {
		fmt.Printf("+%d%% ", initSize)
	}
	fmt.Printf("\n")
	for setSize := uint32(100); setSize < 120; setSize++ {
		fmt.Printf("%d ", setSize)
		for initSize := uint32(0); initSize < 101; initSize++ {
			initSizeVal := uint32(math.Round(float64(setSize) + (float64(setSize*initSize) / 100.0)))
			numOfSets := uint32(math.Round(4000.0 / float64(setSize)))
			rounds := int(math.Round(20_000_000.0 / float64(numOfSets*setSize)))
			measurements := doBenchmark2(rounds, numOfSets, initSizeVal, setSize, 0xABCDEF0123456789)
			nsValues := toNSperAdd(measurements, numOfSets*setSize)
			median := benchmark.Median(nsValues)
			fmt.Printf("%.3f ", median)
		}
		fmt.Printf("\n")
	}
}
