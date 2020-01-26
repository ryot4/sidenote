package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	NotePath = ".sidenote"
)

var (
	ErrContainsDotFile = errors.New("contains dotfile")
	ErrNotDirectory    = errors.New("not a directory")
)

type Directory struct {
	path string
}

func newDirectory(path string) *Directory {
	return &Directory{path: filepath.Clean(path)}
}

func InitDirectory(path string) (*Directory, error) {
	dir := newDirectory(path)

	fi, err := os.Stat(dir.path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir.path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return nil, ErrNotDirectory
	}
	return dir, nil
}

func InitDirectoryLink(path, link string) (*Directory, error) {
	dir, err := InitDirectory(path)
	if err != nil {
		return nil, err
	}
	err = os.Symlink(dir.path, link)
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
		return nil, ErrNotDirectory
	}
	return dir, nil
}

func (dir *Directory) FilePath(path string) (string, error) {
	// prevent directory traversal
	separator := string(filepath.Separator)
	cleanPath := filepath.Clean(filepath.Join(separator, path))
	for _, elem := range strings.Split(cleanPath, separator) {
		if strings.HasPrefix(elem, ".") {
			return "", ErrContainsDotFile
		}
	}
	return filepath.Join(dir.path, cleanPath), nil
}
