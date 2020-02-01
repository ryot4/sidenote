package main

import (
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	notePath := "/path/to/note"
	dir := newDirectory(notePath)

	tests := []struct {
		path   string
		expect string
	}{
		{"foo", "foo"},
		{"/foo", "foo"},
		{"foo.txt", "foo.txt"},
		{"foo.", "foo."},
		{"foo/bar", "foo/bar"},
		{"/foo//bar", "foo/bar"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			path, err := dir.FilePath(tt.path)
			if err != nil {
				t.Error(err)
			}
			relPath, err := filepath.Rel(notePath, path)
			if err != nil {
				t.Error(err)
			}
			if relPath != tt.expect {
				t.Errorf("expect %q, got %q", tt.expect, relPath)
			}
		})
	}
}

func TestFilePathDotFileError(t *testing.T) {
	dir := newDirectory("/path/to/note")

	tests := []string{
		".foo",
		"/.foo",
		"../foo",
		"foo/.bar",
		"foo/.bar/fizz",
		"foo/./bar",
		"foo/../bar",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			_, err := dir.FilePath(tt)
			if err != ErrDotFileInPath {
				t.Errorf("expect ErrContainsDotFile, got %v", err)
			}
		})
	}
}
