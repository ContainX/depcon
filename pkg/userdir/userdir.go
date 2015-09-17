package userdir

import (
	"os"
	"runtime"
)

// Key returns the env var name for the user's home dir
func Key() string {
	if runtime.GOOS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}

// Get returns the home directory of the current user
func Get() string {
	return os.Getenv(Key())
}
