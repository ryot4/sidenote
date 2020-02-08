package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func exitWithSyntaxError(message string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(2)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}

func openDirectory(path string) (*Directory, error) {
	if path == "" {
		path = findDirectory()
	}
	return OpenDirectory(path)
}

// findDirectory searches the directory for notes upward from the current directory.
func findDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		exitWithError(err)
	}
	separator := string(filepath.Separator)
	for dir := wd; dir != "." && dir != separator; dir = filepath.Dir(dir) {
		noteDir := filepath.Join(dir, NoteDirName)
		fi, err := os.Stat(noteDir)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "cannot stat %s: %s. ignoring\n", noteDir, err)
			}
			continue
		}
		if fi.IsDir() {
			return noteDir
		}
	}
	return filepath.Join(wd, NoteDirName)
}
