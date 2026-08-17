// config/paths.go
package config

import (
	"os"
	"path/filepath"
)

// ExecutableDir returns the directory containing the currently running
// binary, with symlinks resolved (handles the common case of a symlinked
// /usr/local/bin entry pointing at a versioned install).
//
// This exists so that config/data file resolution is anchored to *where
// the binary lives* rather than the process's current working directory.
// A packaged app is frequently launched from a different cwd (double
// clicked, invoked via an absolute path from a shell profile, started by
// a service manager, etc.) and relative-path lookups like "./config.yml"
// silently fail or read the wrong file in those cases.
//
// Falls back to the current working directory if the executable path
// can't be determined — this happens under `go run`, where os.Executable
// returns a path inside a temporary build cache.
func ExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return fallbackWd()
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}

	dir := filepath.Dir(resolved)

	// Under `go run`, os.Executable() points into a temp dir
	// (e.g. /tmp/go-buildXXXX/b001/exe/app). That's never where a
	// developer expects config/data to live, so prefer the working
	// directory in that case.
	if isGoRunTempDir(dir) {
		return fallbackWd()
	}

	return dir
}

func fallbackWd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func isGoRunTempDir(dir string) bool {
	tmp := os.TempDir()
	rel, err := filepath.Rel(tmp, dir)
	if err != nil {
		return false
	}
	return rel != ".." && len(rel) > 0 && rel[0] != '.'
}

// fileExists reports whether path exists and is a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
