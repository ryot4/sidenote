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

// FindDirectory searches the directory for notes upward from the current directory.
func FindDirectory() (*Directory, error) {
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

func (dir *Directory) IsAbs() bool {
	return filepath.IsAbs(dir.path)
}

func (dir *Directory) FollowSymlink() (bool, error) {
	fi, err := os.Lstat(dir.path)
	if err != nil {
		return false, err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}
	dir.path, err = os.Readlink(dir.path)
	if err != nil {
		return true, err
	}
	return true, nil
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
