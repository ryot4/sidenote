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

func checkDirectory(path string) (dir *Directory, err error) {
	dir, err = getDirectory(path)
	if err != nil {
		return
	}
	isDir, err := dir.IsDir()
	if err != nil {
		return nil, err
	} else if !isDir {
		return nil, fmt.Errorf("%s exists but is not a directory", path)
	}
	return
}

func getDirectory(path string) (dir *Directory, err error) {
	if path == "" {
		dir, err = findDirectory()
		if err != nil {
			return
		}
	} else {
		dir = NewDirectory(path)
	}
	return
}

// findDirectory searches the directory for notes upward from the current directory.
func findDirectory() (*Directory, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	separator := string(filepath.Separator)
	for cur := wd; cur != "." && cur != separator; cur = filepath.Dir(cur) {
		dir := NewDirectory(filepath.Join(cur, NoteDirName))
		isDir, err := dir.IsDir()
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Fprintln(os.Stderr, err)
		}
		if isDir {
			return dir, nil
		}
	}
	// If not found, return the default one.
	return NewDirectory(NoteDirName), nil
}
