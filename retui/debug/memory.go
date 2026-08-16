package debug

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MemStats is a stable snapshot of Go runtime memory statistics.
//
// The raw memory values come directly from runtime.MemStats.
// Derived rate metrics such as GCPerSec and GCPausePerSec are calculated
// by the memory sampler.
type MemStats struct {
	// Current heap allocation.
	Alloc uint64

	// Cumulative bytes allocated over the lifetime of the process.
	TotalAlloc uint64

	// Total bytes obtained by the Go runtime from the OS.
	Sys uint64

	// Number of live heap objects.
	HeapObjects uint64

	// Total number of completed garbage collections.
	NumGC uint32

	// Current number of goroutines.
	Goroutines int

	// Approximate garbage collections completed per second.
	GCPerSec float64

	// Total GC pause time accumulated since process start.
	GCPause time.Duration

	// Approximate GC pause time accumulated per second.
	GCPausePerSec time.Duration
}

// memorySampler keeps the previous runtime statistics needed to calculate
// rates without making ReadMemStats() stateful.
type memorySampler struct {
	mu sync.RWMutex

	stats MemStats

	lastSampleTime time.Time
	lastNumGC      uint32
	lastPauseNs    uint64

	started bool
}

var memSampler memorySampler

// ReadMemStats captures a fresh snapshot of Go runtime statistics.
//
// This function intentionally has no internal sampling state.
// Every call reads the current runtime statistics directly.
//
// Because runtime.ReadMemStats can briefly stop the world, avoid calling
// this on every render frame in performance-sensitive code.
//
// For repeatedly displayed debug information, prefer CurrentMemStats().
func ReadMemStats() MemStats {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	return MemStats{
		Alloc:       m.Alloc,
		TotalAlloc:  m.TotalAlloc,
		Sys:         m.Sys,
		HeapObjects: m.HeapObjects,
		NumGC:       m.NumGC,
		Goroutines:  runtime.NumGoroutine(),
		GCPause:     time.Duration(m.PauseTotalNs),
	}
}

// SampleMemStats reads the current runtime memory statistics and updates
// the internal memory sampler.
//
// This should be called periodically, for example every 500ms or 1 second,
// rather than once per render frame.
func SampleMemStats() MemStats {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	now := time.Now()

	current := MemStats{
		Alloc:       m.Alloc,
		TotalAlloc:  m.TotalAlloc,
		Sys:         m.Sys,
		HeapObjects: m.HeapObjects,
		NumGC:       m.NumGC,
		Goroutines:  runtime.NumGoroutine(),
		GCPause:     time.Duration(m.PauseTotalNs),
	}

	memSampler.mu.Lock()
	defer memSampler.mu.Unlock()

	// First sample establishes the baseline.
	if !memSampler.started {
		memSampler.stats = current
		memSampler.lastSampleTime = now
		memSampler.lastNumGC = m.NumGC
		memSampler.lastPauseNs = m.PauseTotalNs
		memSampler.started = true

		return current
	}

	elapsed := now.Sub(memSampler.lastSampleTime)

	if elapsed > 0 {
		seconds := elapsed.Seconds()

		// -------------------------------------------------------------
		// GC rate
		// -------------------------------------------------------------

		gcDelta := m.NumGC - memSampler.lastNumGC

		current.GCPerSec = float64(gcDelta) / seconds

		// -------------------------------------------------------------
		// GC pause rate
		// -------------------------------------------------------------

		pauseDelta := m.PauseTotalNs - memSampler.lastPauseNs

		current.GCPausePerSec = time.Duration(
			float64(pauseDelta) / seconds,
		)
	}

	memSampler.stats = current
	memSampler.lastSampleTime = now
	memSampler.lastNumGC = m.NumGC
	memSampler.lastPauseNs = m.PauseTotalNs

	return current
}

// CurrentMemStats returns the most recently sampled memory statistics.
//
// Unlike ReadMemStats(), this does not call runtime.ReadMemStats() and is
// therefore safe to use from the footer/render path.
func CurrentMemStats() MemStats {
	memSampler.mu.RLock()
	defer memSampler.mu.RUnlock()

	return memSampler.stats
}

// ResetMemStatsSampler resets the memory sampler.
//
// The next call to SampleMemStats() establishes a new baseline for GC/s
// and GC pause rate calculations.
func ResetMemStatsSampler() {
	memSampler.mu.Lock()
	defer memSampler.mu.Unlock()

	memSampler.stats = MemStats{}
	memSampler.lastSampleTime = time.Time{}
	memSampler.lastNumGC = 0
	memSampler.lastPauseNs = 0
	memSampler.started = false
}

// String renders the memory statistics in a human-readable format.
func (m MemStats) String() string {
	return fmt.Sprintf(
		"Alloc: %.2f MB | TotalAlloc: %.2f MB | Sys: %.2f MB | HeapObjects: %d | GC: %d | GC/s: %.2f | GC Pause: %s | GC Pause/s: %s | Goroutines: %d",
		mb(m.Alloc),
		mb(m.TotalAlloc),
		mb(m.Sys),
		m.HeapObjects,
		m.NumGC,
		m.GCPerSec,
		formatDuration(m.GCPause),
		formatDuration(m.GCPausePerSec),
		m.Goroutines,
	)
}

// PrintMemory prints the most recently sampled memory statistics.
//
// Deprecated: retained for backward compatibility. Use CurrentMemStats()
// or SampleMemStats() for new code.
func PrintMemory() {
	fmt.Fprintln(currentOutput(), CurrentMemStats().String())
}

// LogMemory reports the most recently sampled memory statistics through
// the leveled debug logger.
func LogMemory() {
	if !shouldLog(LevelInfo) {
		return
	}

	writeEntry(LevelInfo, CurrentMemStats().String())
}

// mb converts bytes to megabytes.
func mb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// formatDuration produces compact durations suitable for the debug footer.
//
// Examples:
//
//	250µs
//	1.24ms
//	2.10s
func formatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%dns", d.Nanoseconds())

	case d < time.Millisecond:
		return fmt.Sprintf("%.0fµs", float64(d)/float64(time.Microsecond))

	case d < time.Second:
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))

	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}
