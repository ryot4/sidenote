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

func FindStorageDir(name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}

	// Find name upward from the current directory.
	wd, err := os.Getwd()
	if err != nil {
		return name, err
	}

	separator := string(filepath.Separator)
	for cur := wd; cur != "." && cur != separator; cur = filepath.Dir(cur) {
		d := filepath.Join(cur, name)
		fi, err := os.Stat(d)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		if fi.IsDir() {
			return d, nil
		} else {
			return "", fmt.Errorf("%s exists but is not a directory", d)
		}
	}

	// Not found. name is assumed to be in the current directory.
	return name, nil
}

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{path: filepath.Clean(path)}
}

func (s *Storage) Init() error {
	_, err := os.Stat(s.path)
	if err == nil {
		return &os.PathError{
			Op:   "initialize",
			Path: s.path,
			Err:  os.ErrExist,
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(s.path, os.ModePerm)
}

func (s *Storage) JoinPath(name string) (string, error) {
	separator := string(filepath.Separator)
	for elem := range strings.SplitSeq(name, separator) {
		if strings.HasPrefix(elem, ".") {
			return "", &DotFileError{Path: name}
		}
	}
	return filepath.Join(s.path, name), nil
}

func (s *Storage) Readdir(name string) ([]os.FileInfo, error) {
	path, err := s.JoinPath(name)
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
