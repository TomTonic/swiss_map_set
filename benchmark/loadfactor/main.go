package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
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

var rngOverhead = getPRNGOverhead()

func getPRNGOverhead() float64 {
	calibrationCalls := 500_000_000 // prng.Uint64() is about 1-2ns, timer resolution is 100ns (windows)
	prng := benchmark.PRNG{State: 0x1234567890abcde}
	start := hrtime.Now()
	for i := 0; i < calibrationCalls; i++ {
		prng.Uint64()
	}
	stop := hrtime.Now()
	nowOverhead := hrtime.Overhead()
	result := float64(stop-start-nowOverhead) / float64(calibrationCalls)
	return result
}

func addBenchmark(rounds int, numberOfSets, initialAlloc, setSize uint32, seed uint64) (measurements []time.Duration) {
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
		timePerRound[r] = endTime - startTime - hrtime.Overhead() - time.Duration((rngOverhead * float64(numberOfSets*setSize)))
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

func printTotalRuntime(start time.Time) {
	end := time.Now()
	fmt.Printf("\nTotal runtime of benchmark: %v\n", end.Sub(start))
}

// Percent is a custom flag type for parsing percent values.
type Step struct {
	isSet       bool
	isPercent   bool
	percent     float64
	integerStep uint32
}

// Set parses the flag value and sets it.
func (p *Step) Set(value string) error {
	if strings.HasSuffix(value, "%") {
		value = strings.TrimSuffix(value, "%")
		p.isPercent = true
	} else {
		p.isPercent = false
	}
	if p.isPercent {
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		p.percent = parsedValue
		p.isSet = true
	} else {
		parsedValue, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		p.integerStep = uint32(parsedValue)
		p.isSet = true
	}
	return nil
}

// String returns the string representation of the flag value.
func (p *Step) String() string {
	if !p.isSet {
		return "1"
	}
	if p.isPercent {
		return fmt.Sprintf("%f%%", p.percent)
	}
	return fmt.Sprintf("%d", p.integerStep)
}

func getNumberOfSteps(setSizeTo uint32, step Step) uint32 {
	if step.isPercent {
		pval := step.percent
		count := uint32(0)
		for f := 0.0; f < 100.0+pval; f += pval {
			count++
		}
		return count
	} else {
		numberOfSteps := setSizeTo
		return numberOfSteps + 1
	}
}

func columnHeadings(setSizeTo uint32, step Step) []string {
	result := make([]string, 0, 100)
	if step.isPercent {
		pval := step.percent
		for f := 0.0; f < 100.0+pval; f += pval {
			result = append(result, fmt.Sprintf("+%.2f%%%% ", f)) // caution: strings are used in fmt.Printf() so encode %% twice
		}
	} else {
		numberOfSteps := setSizeTo
		ival := step.integerStep
		for i := uint32(0); i <= numberOfSteps; i += ival {
			result = append(result, fmt.Sprintf("+%d ", i))
		}
	}
	return result
}

func initSizeValues(currentSetSize, setSizeTo uint32, step Step) []uint32 {
	result := make([]uint32, 0, 100)
	if step.isPercent {
		pval := step.percent
		for f := 0.0; f <= 100.0; f += pval {
			retval := currentSetSize + uint32(math.Round(f*float64(currentSetSize)))
			result = append(result, retval)
		}
	} else {
		numberOfSteps := setSizeTo
		ival := step.integerStep
		for i := currentSetSize; i <= currentSetSize+numberOfSteps; i += ival {
			result = append(result, i)
		}
	}
	return result
}

func main() {
	var fromSetSize, toSetSize, addsPerRound uint
	var assumeAddNs, secondsPerConfig float64

	flag.UintVar(&fromSetSize, "from", 100, "First set size to benchmark (inclusive)")
	flag.UintVar(&toSetSize, "to", 200, "Last set size to benchmark (inclusive)")
	// 50_000 x ~8ns = ~400_000ns; Timer resolution 100ns (Windows) => 0,25% error, i.e. 0,02ns per Add()
	flag.UintVar(&addsPerRound, "apr", 50_000, "AddsPerRound - instructions between two measurements. Balance between memory consumption (cache!) and timer resolution (Windows: 100ns)")
	flag.Float64Var(&secondsPerConfig, "spc", 1.0, "SecondsPerConfig - estimated benchmark time per configuration in seconds")
	flag.Float64Var(&assumeAddNs, "arpa", 8.0, "AssumedRuntimePerAdd - in nanoseconds per instruction. Used to predcict runtimes")
	var step Step
	flag.Var(&step, "step", "Step to increment headroom of pre-allocated sets. Either percent of set size (e.g. \"2.5%\") or absolut value (e.g. \"2\") (default: 1)")

	flag.Parse()

	if !step.isSet {
		step.isSet = true
		step.isPercent = false
		step.integerStep = 1
	}

	totalAddsPerConfig := secondsPerConfig * (1_000_000_000.0 / float64(assumeAddNs))

	fmt.Printf("Architecture:\t\t%s\n", runtime.GOARCH)
	fmt.Printf("OS:\t\t\t%s\n", runtime.GOOS)
	fmt.Printf("Timer precision:\t%fns\n", hrtime.NowPrecision())
	fmt.Printf("hrtime.Now() overhead:\t%v (informative, already subtracted from below measurement values)\n", hrtime.Overhead())
	fmt.Printf("prng.Uint64() overhead:\t%fns (informative, already subtracted from below measurement values)\n", rngOverhead)
	fmt.Printf("Add()'s per round:\t%d\n", addsPerRound)
	fmt.Printf("Add()'s per config:\t%.0f (should result in a benchmarking time of %fs per config)\n", totalAddsPerConfig, secondsPerConfig)
	fmt.Printf("Set3 sizes:\t\tfrom %d to %d, stepsize %v\n", fromSetSize, toSetSize, step.String())
	numberOfStepsPerSetSize := getNumberOfSteps(uint32(toSetSize), step)
	fmt.Printf("Number of configs:\t%d\n", numberOfStepsPerSetSize*uint32(toSetSize-fromSetSize+1))
	totalduration := time.Duration(assumeAddNs * totalAddsPerConfig) // total ns per round
	totalduration *= time.Duration(numberOfStepsPerSetSize)          // different headroom sizes per setSize
	totalduration *= time.Duration(toSetSize - fromSetSize + 1)      // number of setSizes to evaluate
	totalduration = time.Duration(float64(totalduration) * 1.12)     // overhead
	fmt.Printf("Expected total runtime:\t%v (assumption: %fns per Add(prng.Uint64()) and 12%% overhead for housekeeping)\n", totalduration, assumeAddNs)
	fmt.Printf("\n")

	start := time.Now()
	defer printTotalRuntime(start)

	fmt.Printf("setSize ")
	for _, columnH := range columnHeadings(uint32(toSetSize), step) {
		fmt.Printf(columnH)
	}
	fmt.Printf("\n")
	for setSize := uint32(fromSetSize); setSize <= uint32(toSetSize); setSize++ {
		fmt.Printf("%d ", setSize)
		for _, initSize := range initSizeValues(setSize, uint32(toSetSize), step) {
			initSizeVal := uint32(math.Round(float64(setSize) + (float64(setSize*initSize) / 100.0)))
			numOfSets := uint32(math.Round(float64(addsPerRound) / float64(setSize)))
			actualInsertsPerRound := numOfSets * setSize
			rounds := int(math.Round(totalAddsPerConfig / float64(actualInsertsPerRound)))
			measurements := addBenchmark(rounds, numOfSets, initSizeVal, setSize, 0xABCDEF0123456789)
			nsValues := toNSperAdd(measurements, actualInsertsPerRound)
			median := benchmark.Median(nsValues)
			fmt.Printf("%.3f ", median)
		}
		fmt.Printf("\n")
	}
}
