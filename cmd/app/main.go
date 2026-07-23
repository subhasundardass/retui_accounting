package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/subhasundardass/retui/internal/app"
	"github.com/subhasundardass/retui/retui"
	_ "modernc.org/sqlite"
)

func main() {
	retui.Info("🔧 Initializing application...")

	defer func() {
		if r := recover(); r != nil {
			retui.Errorf("❌ Application panicked: %v", r)
			os.Exit(1)
		}
	}()

	bootstrap, err := app.NewBootstrap()
	if err != nil {
		retui.Errorf("❌ Failed to initialize app: %v", err)
		os.Exit(1)
	}
	app.SetBootstrap(bootstrap)

	shutdown := func() {
		retui.Info("🧹 Cleaning up...")
		if err := bootstrap.Shutdown(); err != nil {
			retui.Errorf("❌ Error during shutdown: %v", err)
		}
	}

	//Sutdown
	defer shutdown()
	//Signal
	handleSignals(shutdown)
	//Run App
	run(bootstrap)

	retui.Debug("✅ Application exited successfully")
}

// handleSignals shuts the app down cleanly on SIGINT/SIGTERM instead of
// letting os.Exit bypass main()'s deferred cleanup.
func handleSignals(shutdown func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		retui.Infof("Received signal: %v", sig)
		retui.Info("Shutting down gracefully...")
		shutdown()
		os.Exit(0)
	}()
}

// run starts the TUI event loop using the app's context.
func run(bootstrap *app.Bootstrap) {
	retui.Info("🚀 Starting TUI application...")

	if bootstrap.AppCtx == nil {
		retui.Error("Failed to get app context")
		os.Exit(1)
	}

	renderFn := func(props retui.Props) retui.Element {
		return app.Root(bootstrap.AppCtx, props)
	}
	props := retui.Props{
		Width:  retui.Grow(1),
		Height: retui.Grow(1),
	}

	tuiApp := retui.NewApp(0, 0) // 0,0 = full screen
	tuiApp.Run(renderFn, props)
}
