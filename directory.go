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
		return &os.PathError{"initialize", dir.path, os.ErrExist}
	} else if !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(dir.path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (dir *Directory) IsDir() (bool, error) {
	fi, err := os.Stat(dir.path)
	if err != nil {
		return false, err
	} else if !fi.IsDir() {
		return false, nil
	}
	return true, nil
}

func (dir *Directory) JoinPath(path string) (string, error) {
	separator := string(filepath.Separator)
	for _, elem := range strings.Split(path, separator) {
		if strings.HasPrefix(elem, ".") {
			return "", &DotFileError{Path: path}
		}
	}
	return filepath.Join(dir.path, path), nil
}

func (dir *Directory) Readdir(path string) ([]os.FileInfo, error) {
	realPath, err := dir.JoinPath(path)
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
