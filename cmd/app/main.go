package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/subhasundardass/retui/internal/app"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
	"github.com/subhasundardass/retui/retui/debug"

	_ "modernc.org/sqlite"
)

// rootFn holds the element-tree function that App.Run() actually calls.
//
// It starts out pointing at loadingScreen and is swapped, under mu, once
// bootstrap finishes.
var (
	mu        sync.RWMutex
	rootFn    func(props retui.Props) retui.Element = loadingScreen
	bootstrap *app.Bootstrap

	// exitCode records the process exit code a background goroutine wants
	// once tui.Run() has unwound.
	//
	// Goroutines never call os.Exit directly. They store the code here
	// and signal retui.Exit(). main() exits with this code only after
	// Run() has returned in the main goroutine.
	exitCode atomic.Int32
)

// setRoot atomically swaps the function App.Run() renders.
func setRoot(fn func(props retui.Props) retui.Element) {
	mu.Lock()
	rootFn = fn
	mu.Unlock()
}

// setBootstrap atomically records the bootstrap instance for shutdown().
func setBootstrap(b *app.Bootstrap) {
	mu.Lock()
	bootstrap = b
	mu.Unlock()
}

// getBootstrap atomically reads the bootstrap instance.
func getBootstrap() *app.Bootstrap {
	mu.RLock()
	defer mu.RUnlock()

	return bootstrap
}

// dispatch is the single fixed function passed to tui.Run().
//
// IMPORTANT:
// This measures only UI tree construction, not the complete render.
//
// Build time:
//
//	BeginBuild()
//	    ↓
//	rootFn(props)
//	    ↓
//	EndBuild()
//
// Actual FPS measurement belongs around the complete render cycle inside
// App.Run(), where the returned Element is actually rendered.
func dispatch(props retui.Props) retui.Element {
	// debug.BeginBuild()
	// defer debug.EndBuild()

	debug.BeginFrame()
	defer debug.EndFrame()

	mu.RLock()
	fn := rootFn
	mu.RUnlock()

	return fn(props)
}

// safeGo runs fn in a goroutine with panic recovery.
//
// On panic:
//   - RecoverAndReraise logs the panic and stack trace.
//   - The outer recover catches the re-panic.
//   - exitCode is set.
//   - cleanup is executed.
//   - retui.Exit() signals the application to unwind cleanly.
//
// It deliberately never calls os.Exit from the background goroutine.
func safeGo(name string, cleanup func(), fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				exitCode.Store(1)

				if cleanup != nil {
					cleanup()
				}

				retui.Exit()
			}
		}()

		defer debug.RecoverAndReraise(name)()

		fn()
	}()
}

func main() {
	// -------------------------------------------------------------------------
	// Debug configuration
	// -------------------------------------------------------------------------

	debugEnv := strings.EqualFold(os.Getenv("RETUI_DEBUG"), "1") ||
		strings.EqualFold(os.Getenv("RETUI_DEBUG"), "true")

	verbose := flag.Bool(
		"debug",
		debugEnv,
		"enable verbose debug logging (also settable via RETUI_DEBUG=1)",
	)

	flag.Parse()

	retui.SetDebugMode(*verbose)

	if !*verbose {
		debug.Disable()
		debug.SetLevel(debug.LevelOff)
	} else {
		retui.Info("🐞 Debug mode enabled — verbose logging on")
	}

	retui.Info("🔧 Starting application...")

	// -------------------------------------------------------------------------
	// Application
	// -------------------------------------------------------------------------

	tui := retui.NewApp(0, 0)

	var once sync.Once

	shutdown := func() {
		once.Do(func() {
			retui.Info("Cleaning up...")

			if b := getBootstrap(); b != nil {
				if err := b.Shutdown(); err != nil {
					retui.Errorf("Shutdown failed: %v", err)
				}
			}
		})
	}

	defer shutdown()

	// -------------------------------------------------------------------------
	// Main panic recovery
	// -------------------------------------------------------------------------

	defer func() {
		if r := recover(); r != nil {
			shutdown()
			os.Exit(1)
		}
	}()

	defer debug.RecoverAndReraise("main")()

	// -------------------------------------------------------------------------
	// OS signals
	// -------------------------------------------------------------------------

	handleSignals(shutdown)

	// -------------------------------------------------------------------------
	// Bootstrap
	// -------------------------------------------------------------------------
	//
	// Bootstrap happens in the background so the loading screen can be
	// displayed immediately.
	//
	// Once bootstrap finishes, setRoot() atomically switches dispatch()
	// from loadingScreen to the real application root.
	// -------------------------------------------------------------------------

	safeGo("bootstrap", shutdown, func() {
		retui.Info("Initializing application...")

		b, err := app.NewBootstrap()
		if err != nil {
			retui.Errorf(
				"Failed to initialize application: %v",
				err,
			)

			exitCode.Store(1)
			retui.Exit()
			return
		}

		setBootstrap(b)
		app.SetBootstrap(b)

		setRoot(func(props retui.Props) retui.Element {
			return app.Root(b.AppCtx, props)
		})

		retui.Info("Application initialized.")
	})

	// -------------------------------------------------------------------------
	// Run application
	// -------------------------------------------------------------------------
	//
	// dispatch() initially renders loadingScreen.
	//
	// Once bootstrap calls setRoot(), the next render automatically uses
	// app.Root().
	//
	// dispatch() measures BUILD performance only.
	// Full FPS measurement must happen inside retui.App.Run(), around the
	// actual renderer invocation.
	// -------------------------------------------------------------------------

	tui.Run(
		dispatch,
		retui.Props{
			Width:  retui.Grow(1),
			Height: retui.Grow(1),
		},
	)

	// -------------------------------------------------------------------------
	// Application exited
	// -------------------------------------------------------------------------

	retui.Info("Application exited.")

	if code := exitCode.Load(); code != 0 {
		os.Exit(int(code))
	}
}

// handleSignals starts the signal-handling goroutine.
//
// SIGINT / SIGTERM trigger a clean application shutdown followed by
// retui.Exit(), allowing tui.Run() to unwind normally.
func handleSignals(shutdown func()) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	safeGo("signal-handler", nil, func() {
		sig := <-sigChan

		retui.Infof("Received signal: %v", sig)

		shutdown()
		retui.Exit()
	})
}

// loadingScreen is displayed while application bootstrap is running.
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
			retui.NewStyle().Border(
				retui.Border{
					Top:    true,
					Right:  true,
					Bottom: true,
					Left:   true,
					Color:  retui.Blue,
				},
			),

			components.Spinner("Loading... "),
		),
	)
}
