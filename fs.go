package main

import (
	"io/fs"
	"strings"
)

type dotFileHidingFs struct {
	fs.FS
}

func (fsys dotFileHidingFs) Open(name string) (fs.File, error) {
	// "." represents the root directory of fsys
	if name != "." {
		for part := range strings.SplitSeq(name, "/") {
			if strings.HasPrefix(part, ".") {
				return nil, fs.ErrNotExist
			}
		}
	}

	file, err := fsys.FS.Open(name)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return dotFileHidingFile{file.(fs.ReadDirFile)}, nil
	} else {
		return file, nil
	}
}

type dotFileHidingFile struct {
	fs.ReadDirFile
}

func (dir dotFileHidingFile) ReadDir(n int) (filtered []fs.DirEntry, err error) {
	files, err := dir.ReadDirFile.ReadDir(n)
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			filtered = append(filtered, f)
		}
	}
	return
}
