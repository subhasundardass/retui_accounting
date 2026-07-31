package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/subhasundardass/retui/internal/app"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	_ "modernc.org/sqlite"
)

// rootFn holds the element-tree function that App.Run() actually calls.
// It starts out pointing at loadingScreen and is swapped, under mu, once
// bootstrap finishes — this is the replacement for the SetRoot() method
// that App does not have.
var (
	mu     sync.RWMutex
	rootFn func(props retui.Props) retui.Element = loadingScreen
)

// setRoot atomically swaps the function App.Run() renders.
func setRoot(fn func(props retui.Props) retui.Element) {
	mu.Lock()
	rootFn = fn
	mu.Unlock()
}

// dispatch is the single, fixed function passed to tui.Run(). It always
// delegates to whatever rootFn currently points at.
func dispatch(props retui.Props) retui.Element {
	mu.RLock()
	fn := rootFn
	mu.RUnlock()
	return fn(props)
}

func main() {
	retui.Info("🔧 Starting application...")
	tui := retui.NewApp(0, 0)

	var (
		bootstrap *app.Bootstrap
		once      sync.Once
	)

	shutdown := func() {
		once.Do(func() {
			retui.Info("🧹 Cleaning up...")
			if bootstrap != nil {
				if err := bootstrap.Shutdown(); err != nil {
					retui.Errorf("Shutdown failed: %v", err)
				}
			}
		})
	}
	defer shutdown()

	defer func() {
		if r := recover(); r != nil {
			shutdown()
			retui.Errorf("Application panicked: %v", r)
			os.Exit(1)
		}
	}()

	handleSignals(shutdown)

	// Bootstrap in background.
	go func() {
		retui.Info("🔧 Initializing application...")
		b, err := app.NewBootstrap()
		if err != nil {
			retui.Errorf("Failed to initialize application: %v", err)
			retui.Exit() // App has no Stop(); Exit() signals Run()'s exitCh
			return
		}
		bootstrap = b
		app.SetBootstrap(b)

		setRoot(func(props retui.Props) retui.Element {
			return app.Root(b.AppCtx, props)
		})
	}()

	// Show loading screen immediately; dispatch switches over to the real
	// root as soon as bootstrap calls setRoot above.
	tui.Run(dispatch, retui.Props{
		Width:  retui.Grow(1),
		Height: retui.Grow(1),
	})

	retui.Info("Application exited.")
}

func handleSignals(shutdown func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		retui.Infof("Received signal: %v", sig)
		shutdown()
		retui.Exit() // App has no Stop(); Exit() signals Run()'s exitCh
	}()
}

func loadingScreen(props retui.Props) retui.Element {
	return retui.Box(
		retui.Props{
			Width:   retui.Grow(1),
			Height:  retui.Grow(1),
			Align:   retui.AlignCenter,
			Justify: retui.JustifyCenter,
		},
		retui.NewStyle(),
		retui.Box(
			retui.Props{
				Padding: [4]int{0, 2, 0, 2},
			},
			retui.NewStyle().Border(retui.Border{Top: true, Right: true, Bottom: true, Left: true, Color: retui.Blue}),
			components.Spinner("Loading... "),
		),
	)
}
