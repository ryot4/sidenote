package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DotFileError struct {
	Path string
}

func (e *DotFileError) Error() string {
	return fmt.Sprintf("path %s contains dotfile", e.Path)
}

type Directory struct {
	path string
}

func NewDirectory(path string) *Directory {
	return &Directory{path: filepath.Clean(path)}
}

func (dir *Directory) Init() error {
	_, err := os.Stat(dir.path)
	if err == nil {
		return &os.PathError{
			Op:   "initialize",
			Path: dir.path,
			Err:  os.ErrExist,
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dir.path, os.ModePerm)
}

func (dir *Directory) JoinPath(name string) (string, error) {
	separator := string(filepath.Separator)
	for _, elem := range strings.Split(name, separator) {
		if strings.HasPrefix(elem, ".") {
			return "", &DotFileError{Path: name}
		}
	}
	return filepath.Join(dir.path, name), nil
}

func (dir *Directory) Readdir(name string) ([]os.FileInfo, error) {
	path, err := dir.JoinPath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	children, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return nil, err
	}

	n := 0
	for _, c := range children {
		if !strings.HasPrefix(c.Name(), ".") {
			children[n] = c
			n++
		}
	}
	return children[:n], nil
}
