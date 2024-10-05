package benchmark

import (
	"runtime"
)

// RuntimeNano returns the current value of the runtime clock in nanoseconds.
// from time.go, Copyright 2009 The Go Authors, BSD-style License
//
// go:linkname RuntimeNano runtime.nanotime
//func RuntimeNano() int64

func CurrentMemUsage() uint64 {
	var startMem runtime.MemStats
	runtime.ReadMemStats(&startMem)
	result := startMem.HeapAlloc + startMem.StackInuse + startMem.StackSys
	return result
}

/*
func CurrentMemUsage() uint64 {
	// "github.com/shirou/gopsutil/v4/mem"
	v, _ := mem.VirtualMemory()
	result := v.Total - v.Free
	return result
}
*/
