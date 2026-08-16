package debug

import (
	"sync"
	"time"
)

// FrameStats contains performance statistics for completed render frames.
//
// FPS is the actual observed number of completed frames per second over
// the FPS sampling window.
//
// FrameTime is the duration of the most recently completed frame.
//
// AvgFrameTime is the average duration of completed frames in the
// frame-duration sampling window.
//
// FrameCount is the total number of completed frames observed since
// the tracker started or was last reset.

// FPS       → How often Retui actually completes frames
// Frame     → Cost of the latest frame
// Avg       → Average cost of rendering a frame
// MEM       → Current allocated memory
// GC        → Total GC cycles
// G         → Current goroutines

type FrameStats struct {
	FPS          float64
	FrameTime    time.Duration
	AvgFrameTime time.Duration
	FrameCount   uint64
}

const (
	defaultFrameWindowSize = 60
)

// frameTracker tracks completed render frames.
//
// Duration samples and FPS samples are intentionally separate:
//
//   - window contains frame durations.
//   - fpsTimes contains timestamps of completed frames.
//
// This is important because:
//
//	1 / AvgFrameTime
//
// is frame throughput, not actual observed FPS.
type frameTracker struct {
	mu sync.Mutex

	// Current frame start time.
	start time.Time

	// Frame duration rolling window.
	window     []time.Duration
	windowPos  int
	windowFull bool
	windowSize int

	// Total completed frames.
	count uint64

	// Completed-frame timestamps used for actual FPS calculation.
	fpsTimes      []time.Time
	fpsTimePos    int
	fpsTimeFull   bool
	fpsWindowSize int
}

var frames = frameTracker{
	window:        make([]time.Duration, defaultFrameWindowSize),
	windowSize:    defaultFrameWindowSize,
	fpsTimes:      make([]time.Time, defaultFrameWindowSize),
	fpsWindowSize: defaultFrameWindowSize,
}

// BeginFrame marks the beginning of a complete render frame.
//
// It should be called once immediately before the complete render cycle:
//
//	BeginFrame()
//	    build UI
//	    layout
//	    paint
//	    flush
//	EndFrame()
//
// Do not call this around only the root component function. That would
// measure UI building time rather than a complete render frame.
func BeginFrame() {
	if !Enabled() {
		return
	}

	frames.mu.Lock()
	defer frames.mu.Unlock()

	// Protect against accidental nested BeginFrame calls.
	if !frames.start.IsZero() {
		return
	}

	frames.start = timeNow()
}

// EndFrame marks the completion of a render frame.
//
// It records:
//
//   - the frame duration,
//   - the total frame count,
//   - the completion timestamp used for FPS calculation.
func EndFrame() {
	if !Enabled() {
		return
	}

	now := timeNow()

	frames.mu.Lock()
	defer frames.mu.Unlock()

	if frames.start.IsZero() {
		return
	}

	// Calculate duration of this complete frame.
	duration := now.Sub(frames.start)

	// Record frame duration.
	frames.count++

	frames.window[frames.windowPos] = duration
	frames.windowPos = (frames.windowPos + 1) % frames.windowSize

	if frames.windowPos == 0 {
		frames.windowFull = true
	}

	// Record completion timestamp for actual FPS calculation.
	frames.fpsTimes[frames.fpsTimePos] = now
	frames.fpsTimePos = (frames.fpsTimePos + 1) % frames.fpsWindowSize

	if frames.fpsTimePos == 0 {
		frames.fpsTimeFull = true
	}

	// Mark frame as completed.
	frames.start = time.Time{}
}

// CurrentFrameStats returns the current render-frame statistics.
//
// FPS is calculated from actual completed frame timestamps:
//
//	frames completed / elapsed wall-clock time
//
// It is intentionally NOT calculated as:
//
//	time.Second / AvgFrameTime
//
// because that represents theoretical frame throughput rather than
// observed frames per second.
func CurrentFrameStats() FrameStats {
	if !Enabled() {
		return FrameStats{}
	}

	frames.mu.Lock()
	defer frames.mu.Unlock()

	if frames.count == 0 {
		return FrameStats{}
	}

	// ---------------------------------------------------------------------
	// Frame duration statistics
	// ---------------------------------------------------------------------

	n := frames.windowPos

	if frames.windowFull {
		n = frames.windowSize
	}

	var avg time.Duration

	if n > 0 {
		var sum time.Duration

		for i := 0; i < n; i++ {
			sum += frames.window[i]
		}

		avg = sum / time.Duration(n)
	}

	lastFrameIndex := frames.windowPos - 1

	if lastFrameIndex < 0 {
		lastFrameIndex = frames.windowSize - 1
	}

	frameTime := frames.window[lastFrameIndex]

	// ---------------------------------------------------------------------
	// Actual FPS
	// ---------------------------------------------------------------------

	fps := calculateFPSLocked()

	return FrameStats{
		FPS:          fps,
		FrameTime:    frameTime,
		AvgFrameTime: avg,
		FrameCount:   frames.count,
	}
}

// calculateFPSLocked calculates actual observed FPS from completed-frame
// timestamps.
//
// frames.mu must be held by the caller.
func calculateFPSLocked() float64 {
	n := frames.fpsTimePos

	if frames.fpsTimeFull {
		n = frames.fpsWindowSize
	}

	if n < 2 {
		return 0
	}

	// The ring buffer contains timestamps, but they may wrap around.
	//
	// Find the oldest and newest valid timestamps by scanning the samples.
	var oldest time.Time
	var newest time.Time

	for i := 0; i < n; i++ {
		t := frames.fpsTimes[i]

		if t.IsZero() {
			continue
		}

		if oldest.IsZero() || t.Before(oldest) {
			oldest = t
		}

		if newest.IsZero() || t.After(newest) {
			newest = t
		}
	}

	if oldest.IsZero() || newest.IsZero() {
		return 0
	}

	elapsed := newest.Sub(oldest)

	if elapsed <= 0 {
		return 0
	}

	// N completed frames have N-1 intervals between their timestamps.
	//
	// Example:
	//
	//   60 frames
	//   over ~1 second
	//
	// gives approximately 59 frame intervals, so use the elapsed period
	// between the first and last completed frame.
	return float64(n-1) / elapsed.Seconds()
}

// ResetFrameStats resets all frame statistics.
func ResetFrameStats() {
	frames.mu.Lock()
	defer frames.mu.Unlock()

	frames.start = time.Time{}

	frames.windowPos = 0
	frames.windowFull = false

	for i := range frames.window {
		frames.window[i] = 0
	}

	frames.count = 0

	frames.fpsTimePos = 0
	frames.fpsTimeFull = false

	for i := range frames.fpsTimes {
		frames.fpsTimes[i] = time.Time{}
	}
}
