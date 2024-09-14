package Set3

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rand_ops_config = []struct {
	numSets int
	setSize int
	setVar  int
	mod     int
}{
	{numSets: 7, setSize: 20, setVar: 50, mod: 140},
	{numSets: 21, setSize: 5, setVar: 5, mod: 128},
	{numSets: 7, setSize: 1000, setVar: 50, mod: 1200},
	{numSets: 21, setSize: 5, setVar: 5, mod: 1280},
	{numSets: 70, setSize: 20, setVar: 50, mod: 140},
	{numSets: 210, setSize: 3000, setVar: 3000, mod: 9999},
}

func TestRandomOps(t *testing.T) {
	t.Skip("unskip for cpu profiling - runs 0.5-1 minutes")
	// Start profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Could not create CPU profile:", err)
		return
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("Could not start CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	for _, cfg := range rand_ops_config {
		for iterations := 0; iterations < 100; iterations++ {
			sets := make([]*Set3[uint32], cfg.numSets)
			// fill all sets
			for i := 0; i < cfg.numSets; i++ {
				sets[i] = NewSet3[uint32]()
				targetSize := cfg.setSize + rand.Intn(cfg.setVar)
				for j := 0; j < targetSize; j++ {
					sets[i].Add(rand.Uint32() % uint32(cfg.mod))
				}
			}
			// delete some elements
			for i := 0; i < cfg.numSets/3; i++ {
				idx := rand.Intn(cfg.numSets)
				for i := 0; i < cfg.mod*2; i++ {
					sets[idx].Remove(rand.Uint32() % uint32(cfg.mod))
				}
				sets[idx].AddAllFrom([]uint32{rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod),
					rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod)})
			}
			// delete some sets completely
			for i := 0; i < cfg.numSets/4; i++ {
				idx := rand.Intn(cfg.numSets)
				clone := sets[idx].Clone()
				sets[idx].Clear()
				sets[idx].AddAll(clone)
			}
			// refill all sets
			for i := 0; i < cfg.numSets; i++ {
				targetSize := uint32(cfg.setSize + rand.Intn(cfg.setVar))
				for j := sets[i].Count(); j < targetSize; j++ {
					sets[i].Add(rand.Uint32() % uint32(cfg.mod))
				}
			}

			// add all elements to a superset
			superset := NewSet3[uint32]()
			for i := 0; i < cfg.numSets; i++ {
				superset.AddAll(sets[i])
			}

			for i := 0; i < cfg.numSets; i++ {
				intersect := superset.Intersection(sets[i])
				equal := intersect.Equals(sets[i])
				stringIntersect := intersect.String()
				stringSet := sets[i].String()
				assert.True(t, equal, stringIntersect+" != "+stringSet)
				newSet := AsSet3([]uint32{rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod),
					rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod), rand.Uint32() % uint32(cfg.mod)})
				stringNewSet := newSet.String()
				unequal := intersect.Equals(newSet)
				assert.False(t, unequal, stringIntersect+" == "+stringNewSet)
			}
		}
	}
}
