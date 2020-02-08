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
		return findDirectory()
	} else {
		return OpenDirectory(path)
	}
}

// findDirectory searches the directory for notes upward from the current directory.
func findDirectory() (*Directory, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	separator := string(filepath.Separator)
	for cur := wd; cur != "." && cur != separator; cur = filepath.Dir(cur) {
		dir, err := OpenDirectory(filepath.Join(cur, NoteDirName))
		if err == nil {
			return dir, nil
		}
		if os.IsNotExist(err) || IsNotDirectory(err) {
			continue
		}
		fmt.Fprintln(os.Stderr, err)
	}
	return nil, fmt.Errorf("%s directory is not found", NoteDirName)
}
