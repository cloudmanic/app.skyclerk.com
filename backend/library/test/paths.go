package test

import (
	"path/filepath"
	"runtime"
)

// GetProjectRoot returns the absolute path to the project root
func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// Go up two directories from library/test to get to backend root
	return filepath.Join(basepath, "..", "..")
}

// GetTestFilePath returns the absolute path to a test file
func GetTestFilePath(filename string) string {
	return filepath.Join(GetProjectRoot(), "library", "test", "files", filename)
}