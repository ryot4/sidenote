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

type NotDirectoryError struct {
	Path string
}

func (e *NotDirectoryError) Error() string {
	return fmt.Sprintf("%s is not a directory", e.Path)
}

func IsNotDirectory(err error) bool {
	_, ok := err.(*NotDirectoryError)
	return ok
}

type Directory struct {
	path string
}

func newDirectory(path string) *Directory {
	return &Directory{path: filepath.Clean(path)}
}

func InitDirectory(path string) (*Directory, error) {
	dir := newDirectory(path)

	_, err := os.Stat(dir.path)
	if err == nil {
		return dir, &os.PathError{"initialize", dir.path, os.ErrExist}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	err = os.MkdirAll(dir.path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return dir, nil
}

func OpenDirectory(path string) (*Directory, error) {
	dir := newDirectory(path)

	fi, err := os.Stat(dir.path)
	if err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return nil, &NotDirectoryError{Path: dir.path}
	}
	return dir, nil
}

func (dir *Directory) AbsPath(path string) (string, error) {
	separator := string(filepath.Separator)
	for _, elem := range strings.Split(path, separator) {
		if strings.HasPrefix(elem, ".") {
			return "", &DotFileError{Path: path}
		}
	}
	return filepath.Join(dir.path, path), nil
}

func (dir *Directory) Readdir(path string) ([]os.FileInfo, error) {
	realPath, err := dir.AbsPath(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(realPath)
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
