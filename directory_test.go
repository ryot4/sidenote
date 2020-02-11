package main

import (
	"testing"
)

func TestFilePath(t *testing.T) {
	noteDir := ".notes"
	dir := NewDirectory(noteDir)

	tests := []struct {
		path   string
		expect string
	}{
		{"foo", ".notes/foo"},
		{"/foo", ".notes/foo"},
		{"foo.txt", ".notes/foo.txt"},
		{"foo.", ".notes/foo."},
		{"foo/bar", ".notes/foo/bar"},
		{"/foo//bar", ".notes/foo/bar"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			path, err := dir.JoinPath(tt.path)
			if err != nil {
				t.Error(err)
			}
			if path != tt.expect {
				t.Errorf("expect %q, got %q", tt.expect, path)
			}
		})
	}
}

func TestFilePathDotFileError(t *testing.T) {
	dir := NewDirectory("/path/to/note")

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
			_, err := dir.JoinPath(tt)
			_, ok := err.(*DotFileError)
			if !ok {
				t.Errorf("expect DotFileError, got %v", err)
			}
		})
	}
}
