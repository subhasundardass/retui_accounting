package layout

import (
	"fmt"
	"sync"
	"time"

	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/debug"
)

var (
	clockMu   sync.RWMutex
	clockNow  = time.Now()
	clockOnce sync.Once
)

func startClockOnce() {
	clockOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				clockMu.Lock()
				clockNow = time.Now()
				clockMu.Unlock()

				retui.RequestRender() // <-- see note below
			}
		}()
	})
}

func readClock() time.Time {
	clockMu.RLock()
	defer clockMu.RUnlock()
	return clockNow
}

func Footer(props retui.Props) retui.Element {
	startClockOnce()

	now := readClock()
	blinkOn := now.Second()%2 == 0

	dotStyle := retui.NewStyle().Foreground(retui.Green)
	if !blinkOn {
		dotStyle = retui.NewStyle().Foreground(retui.Gray(1))
	}

	// Left side of footer
	leftChildren := []retui.Element{
		retui.Text("Ready", retui.Style{}),
	}

	if debug.Enabled() {
		stats := debug.CurrentFrameStats()
		mem := debug.ReadMemStats()

		debugText := fmt.Sprintf(
			"FPS: %.1f  Frame: %s  Avg: %s  MEM: %.1fMB  Go: %d",
			stats.FPS,
			formatDuration(stats.FrameTime),
			formatDuration(stats.AvgFrameTime),
			float64(mem.Alloc)/1024/1024,
			mem.Goroutines,
		)

		leftChildren = append(
			leftChildren,
			retui.Text(
				"|",
				retui.NewStyle().Foreground(retui.BrightWhite),
			),
			retui.Text(
				debugText,
				retui.NewStyle().Foreground(retui.BrightBlack),
			),
		)
	}

	footer := retui.Box(
		retui.Props{
			Width:   retui.Grow(1),
			Justify: retui.JustifySpaceBetween,
			Align:   retui.AlignCenter,
		},
		retui.NewStyle().Background(retui.Gray(1)),

		// LEFT
		retui.Box(
			retui.Props{Gap: 2},
			retui.NewStyle(),
			leftChildren...,
		),

		// RIGHT
		retui.Box(
			retui.Props{Gap: 2},
			retui.NewStyle(),
			retui.Text("●", dotStyle),
			retui.Text("v0.8.0", retui.Style{}),
			retui.Text(
				now.Format("02/01/2006 03:04:05 PM"),
				retui.Style{},
			),
		),
	)

	return footer
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%.0fns", float64(d.Nanoseconds()))

	case d < time.Millisecond:
		return fmt.Sprintf("%.0fµs", float64(d.Microseconds()))

	default:
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
}
