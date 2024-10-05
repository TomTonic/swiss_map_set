package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	set3 "github.com/TomTonic/Set3"
	"github.com/TomTonic/Set3/benchmark"
	"github.com/TomTonic/objectsize"
)

func doBenchmark(rounds, numberOfSets, initialAlloc, setSize uint32, seed uint64) {
	fmt.Printf("%d;%d;%d;%d;", rounds, numberOfSets, initialAlloc, setSize)
	set := make([]*set3.Set3[uint64], numberOfSets)
	timePerAdd := make([]float64, rounds)
	prng := benchmark.PRNG{State: seed}
	var startMem, endMem runtime.MemStats
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
	reqMem2 := float64(val) / float64(numberOfSets*setSize)
	fmt.Printf("%.3f;", reqMem2)
	runtime.GC()
	for r := range rounds {
		for s := range numberOfSets {
			set[s].Clear()
		}
		startTime := time.Now().UnixNano()
		for s := range numberOfSets {
			currentSet := set[s]
			for range setSize {
				currentSet.Add(prng.Uint64())
			}
		}
		endTime := time.Now().UnixNano()
		timePerAdd[r] = float64(endTime-startTime) / float64(numberOfSets*setSize)
	}
	debug.SetGCPercent(100)
	med := benchmark.Median(timePerAdd)
	avg, variance, stddev := benchmark.Statistics(timePerAdd)
	fmt.Printf("%.4f;%.4f;%.4f;%.4f\n", med, avg, variance, stddev)
}

func main() {
	fmt.Printf("rounds;numberOfSets;initialAlloc;setSize;memPerElem;medianForAdd;avgForAdd;var;stddev\n")
	for i := range 200 {
		doBenchmark(200_001, 10, 200+uint32(i), 200, 0xABCDEF0123456789)
	}
}
