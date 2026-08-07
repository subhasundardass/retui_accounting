//go:build !windows

package retui // match your actual package name

import (
	"os"
	"os/signal"
	"syscall"
)

// newResizeChan returns a channel that fires whenever the terminal is
// resized, plus a cleanup func to stop listening.
func newResizeChan() (chan os.Signal, func()) {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	return resize, func() { signal.Stop(resize) }
}
