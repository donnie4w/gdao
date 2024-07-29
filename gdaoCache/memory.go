// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoCache

import (
	"runtime"
	"sync"
	"time"
)

type memoryMonitor struct {
	lastNumGC        uint32
	lastCheckTime    time.Time
	mux              sync.Mutex
	threshold        float64
	pauselimit       time.Duration
	gcFrequencyLimit float64
}

// newMemoryMonitor creates a new instance of MemoryMonitor with the specified thresholds for memory pressure detection.
//
// Parameters:
//
//	threshold (float64): The maximum acceptable heap usage rate (HeapAlloc / HeapSys).
//	                     If the current rate exceeds this value, memory pressure is considered high.
//	pauselimit (time.Duration): The maximum acceptable duration for the most recent GC pause.
//	                            If the last GC pause duration exceeds this limit, memory pressure is considered high.
//	gcFrequencyLimit (float64): The maximum acceptable frequency of garbage collections (GCs) in GCs per second.
//	                            If the current GC frequency exceeds this value, memory pressure is considered high.
//
// The function initializes the MemoryMonitor with the current GC count and the current time, which are used
// for subsequent memory pressure checks to calculate the GC frequency.
func newMemoryMonitor(threshold float64, pauselimit time.Duration, gcFrequencyLimit float64) *memoryMonitor {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return &memoryMonitor{
		lastNumGC:        mem.NumGC,
		lastCheckTime:    time.Now(),
		threshold:        threshold,
		pauselimit:       pauselimit,
		gcFrequencyLimit: gcFrequencyLimit,
	}
}

// CheckMemoryPressure checks if the current memory usage indicates high memory pressure.
// It evaluates the heap usage rate, the duration of the last GC pause, and the GC frequency.
//
// Returns:
//
//	bool: true if the memory pressure is high based on the given thresholds, false otherwise.
//
// The method performs the following checks:
//  1. Heap Usage Rate: If the ratio of HeapAlloc to HeapSys exceeds the threshold, memory pressure is considered high.
//  2. Last GC Pause Duration: If the duration of the last GC pause exceeds the given pauselimit, memory pressure is considered high.
//  3. GC Frequency: If the frequency of garbage collections (GCs) exceeds the gcFrequencyLimit, memory pressure is considered high.
//
// The method uses a mutex to ensure thread-safe access to the state variables and updates them
// with the latest GC statistics for subsequent checks.
func (m *memoryMonitor) CheckMemoryPressure() bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	heapUsageRate := float64(mem.HeapAlloc) / float64(mem.HeapSys)

	var lastPauseDuration time.Duration
	if mem.NumGC > 0 {
		lastPauseIndex := (mem.NumGC - 1) % uint32(len(mem.PauseNs))
		lastPauseDuration = time.Duration(mem.PauseNs[lastPauseIndex])
	}

	currentTime := time.Now()
	timeElapsed := currentTime.Sub(m.lastCheckTime).Seconds()
	gcCount := mem.NumGC - m.lastNumGC
	gcFrequency := float64(gcCount) / timeElapsed

	m.lastNumGC = mem.NumGC
	m.lastCheckTime = currentTime

	return heapUsageRate > m.threshold || lastPauseDuration > m.pauselimit || gcFrequency > m.gcFrequencyLimit
}
