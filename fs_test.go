package main

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func setupFS() *fstest.MapFS {
	return &fstest.MapFS{
		".dotfile": {},
		".dotdotfile": {},
		".dotdir/file": {},
		"TODO.txt": {},
		"log/.dotfile": {},
		"log/.dotdir/file": {},
		"log/2026-08-30.log": {},
	}
}

func TestDotFileHidingFsOpen(t *testing.T) {
	fsys := dotFileHidingFs{setupFS()}

	tests := []struct{
		name string
		readable bool
	}{
		{".", true},
		{".dotfile", false},
		{"..dotdotfile", false},
		{".dotdir", false},
		{".dotdir/file", false},
		{"TODO.txt", true},
		{"log", true},
		{"log/.dotfile", false},
		{"log/.dotdir", false},
		{"log/.dotdir/file", false},
		{"log/2026-08-30.log", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fsys.Open(tt.name)
			if tt.readable {
				if err != nil {
					t.Errorf("%q should be readable (%v)", tt.name, err)
				}
			} else {
				if err == nil {
					t.Errorf("%q should be unreadable", tt.name)
				} else if !errors.Is(err, fs.ErrNotExist) {
					t.Errorf("%q is unreadable, but error is unexpected (%v)", tt.name, err)
				}
			}
		})
	}
}

func TestDotFileHidingFsReadDir(t *testing.T) {
	fsys := dotFileHidingFs{setupFS()}

	tests := []struct{
		name string
		readable bool
	}{
		{".dotfile", false},
		{"..dotdotfile", false},
		{".dotdir", false},
		{"TODO.txt", true},
		{"log", true},
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		t.Errorf("ReadDir() failed: %v", err)
	}

	names := make(map[string]struct{})
	for _, entry := range entries {
		names[entry.Name()] = struct{}{}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := names[tt.name]
			if tt.readable && !exists {
				t.Errorf("%q should be listed", tt.name)
			} else if !tt.readable && exists {
				t.Errorf("%q should not be listed", tt.name)
			}
		})
	}
}
