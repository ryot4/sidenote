package main

import (
	"fmt"
)

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
		dir, err = FindDirectory()
		if err != nil {
			return
		}
	} else {
		dir = NewDirectory(path)
	}
	return
}
