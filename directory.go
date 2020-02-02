package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrDotFileInPath = errors.New("contains dotfile")
	ErrNotDirectory  = errors.New("not a directory")
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
	separator := string(filepath.Separator)
	for _, elem := range strings.Split(path, separator) {
		if strings.HasPrefix(elem, ".") {
			return "", ErrDotFileInPath
		}
	}
	return filepath.Join(dir.path, path), nil
}

func (dir *Directory) Readdir(path string) ([]os.FileInfo, error) {
	realPath, err := dir.FilePath(path)
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
